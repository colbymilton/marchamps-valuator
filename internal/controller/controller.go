package controller

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	marvel "github.com/colbymilton/marchamps-valuator/internal/marvelcdb"
	"github.com/colbymilton/marchamps-valuator/internal/utils"
	mw "github.com/colbymilton/marchamps-valuator/pkg/mongoWrapper"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	cPacks      = "packs"
	cCards      = "cards"
	cDecks      = "decks"
	cHeroes     = "heroes"
	cCardValues = "card-values"
	cPackValues = "pack-values"
	cMeta       = "meta"

	cMetaId     = 1
	cUpdateFreq = time.Hour * 12
)

type Valuator struct {
	mCli *marvel.MarvelClient
	db   *mw.MongoDB

	cards map[string]*Card
	mutex sync.Mutex
}

func NewValuator() *Valuator {
	v := &Valuator{cards: make(map[string]*Card)}

	mcli, err := marvel.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	v.mCli = mcli

	mongoConnStr := os.Getenv("MONGO_CONN_STRING")
	v.db = mw.NewMongoDB(mongoConnStr, "marchamps-valuator")

	if os.Getenv("DELETE_ALL_ON_STARTUP") == "true" {
		v.db.EmptyCollection(cCards)
		v.db.EmptyCollection(cPacks)
		v.db.EmptyCollection(cHeroes)
		v.db.EmptyCollection(cCardValues)
		v.db.EmptyCollection(cPackValues)
	}

	return v
}

// ValueAllCards handles the /card_values endpoint
func (v *Valuator) ValueAllCards(owned []string) ([]*CardValue, error) {
	if err := v.updateIfNeeded(); err != nil {
		return nil, err
	}

	// grab base card values from db
	cvs, err := mw.GetMany[CardValue](v.db, cCardValues, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// make map of owned cards
	ownedCards, err := v.getCardsFromPacks(owned)
	if err != nil {
		return nil, err
	}

	// get all heroes
	allHeroes, err := mw.GetMany[Hero](v.db, cHeroes, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// make map of owned heroes
	ownedHeroes := map[string]*Hero{}
	for _, hero := range allHeroes {
		if utils.StringsContains(owned, hero.PackCode) {
			ownedHeroes[hero.Code] = hero
		}
	}

	// modify base pack values based on owned cards
	for _, cv := range cvs {
		adjustCardValue(cv, ownedCards, ownedHeroes, allHeroes, "", map[string]float64{})
	}

	sort.Slice(cvs, func(i, j int) bool { return cvs[i].Value > cvs[j].Value })

	return cvs, nil
}

// ValueAllPacks handles the /pack_values endpoint
func (v *Valuator) ValueAllPacks(owned []string, aspectWeights map[string]float64) ([]*PackValue, error) {
	if err := v.updateIfNeeded(); err != nil {
		return nil, err
	}

	// grab base pack values from db
	pvs, err := mw.GetMany[PackValue](v.db, cPackValues, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// make map of owned cards
	ownedCards, err := v.getCardsFromPacks(owned)
	if err != nil {
		return nil, err
	}

	// get all heroes
	allHeroes, err := mw.GetMany[Hero](v.db, cHeroes, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// make map of owned heroes
	ownedHeroes := map[string]*Hero{}
	for _, hero := range allHeroes {
		if utils.StringsContains(owned, hero.PackCode) {
			ownedHeroes[hero.Code] = hero
		}
	}

	// modify base pack values based on owned cards
	for _, pv := range pvs {
		for _, cv := range pv.CardValues {
			adjustCardValue(cv, ownedCards, ownedHeroes, allHeroes, pv.Code, aspectWeights)
		}
		sort.Slice(pv.CardValues, func(i, j int) bool { return pv.CardValues[i].Value > pv.CardValues[j].Value })
		pv.Calculate()
	}

	sort.Slice(pvs, func(i, j int) bool { return pvs[i].ValueSum > pvs[j].ValueSum })

	return pvs, nil
}

// GetPacks handles the /packs endpoint
func (v *Valuator) GetPacks() ([]*marvel.Pack, error) {
	if err := v.updateIfNeeded(); err != nil {
		return nil, err
	}
	return mw.GetMany[marvel.Pack](v.db, cPacks, mw.BsonNoneD, bson.M{"availablestr": 1})
}

func (v *Valuator) updateIfNeeded() error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	// add meta data (if it already exists, this will be ignored)
	if err := mw.CreateMany[Meta](v.db, cMeta, []*Meta{{Id: cMetaId, LastUpdated: time.Time{}}}); err != nil {
		return err
	}

	// get meta data from db
	meta, err := mw.GetOne[Meta](v.db, cMeta, mw.BuildEqualsFilter("_id", cMetaId), mw.BsonNoneM)
	if err != nil {
		return err
	}

	setup := meta.LastUpdated.IsZero()
	needsUpdated := meta.LastUpdated.Add(cUpdateFreq).Before(time.Now())

	if setup {
		// since this is a first-time setup, we need to block until the database has been initialized
		return v.updateAll()

	} else if needsUpdated {
		// since we already have some data, we can do the update in the background
		go func() {
			if err := v.updateAll(); err != nil {
				log.Println("error when updating in the background:", err)
			}
		}()
		return nil
	}

	// doesn't need to be updated at all
	return nil
}

func (v *Valuator) updateAll() error {
	// check mongo first, just to save a marvel endpoint call
	if err := v.db.Ping(); err != nil {
		return err
	}

	// update packs
	if err := v.updatePacks(); err != nil {
		return err
	}

	// update cards
	if err := v.updateCards(); err != nil {
		return err
	}

	// update heroes
	if err := v.updateHeroes(); err != nil {
		return err
	}

	// update decks
	decksAdded, err := v.updateDecks()
	if err != nil {
		return err
	}

	// update card values
	if true { // todo remove true
		if err := v.updateCardValues(); err != nil {
			return err
		}
	}

	// update pack values
	if decksAdded || true { // todo remove true
		if err := v.updatePackValues(); err != nil {
			return err
		}
	}

	// update lastUpdated time
	return mw.ReplaceOneID[Meta](v.db, cMeta, &Meta{Id: cMetaId, LastUpdated: time.Now()})
}

func (v *Valuator) updatePacks() error {
	log.Println("Updating local list of packs.")
	packs, err := v.mCli.GetAllPacks()
	if err != nil {
		return err
	}

	// defer log.Println("Local pack count:", mw.GetCollectionSize(v.db, cPacks))

	return mw.CreateMany(v.db, cPacks, packs)
}

func (v *Valuator) updateCards() error {
	log.Println("Updating local list of cards.")

	// get packs
	packs, err := mw.GetAll[marvel.Pack](v.db, cPacks)
	if err != nil {
		return err
	}

	// build pack map
	packMap := make(map[string]*marvel.Pack)
	for _, pack := range packs {
		packMap[pack.Code] = pack
	}

	// get latest cards
	mCards, err := v.mCli.GetAllCards()
	if err != nil {
		return err
	}

	// convert marvel api cards to local cards
	dups := make([]*marvel.Card, 0)
	for _, mCard := range mCards {
		if mCard.DuplicateOf != "" {
			dups = append(dups, mCard)
			continue // handle dups later
		}

		card := &Card{
			Code:          mCard.Code,
			Name:          mCard.Name,
			Subname:       mCard.SubName,
			PackCodes:     []string{mCard.PackCode},
			TypeCode:      mCard.TypeCode,
			Aspect:        mCard.FactionCode,
			Traits:        strings.Split(mCard.Traits, ". "),
			LockingTraits: parseLockingTraits(mCard.Text),
			DateAvailable: packMap[mCard.PackCode].DateAvailable(),
			DuplicateBy:   []string{},
			Text:          mCard.Text,
			CardSetName:   mCard.CardSetName,
			ImageSrc:      mCard.ImageSrc,
		}

		if mCard.LinkedCard != nil {
			card.LinkedCardCode = mCard.LinkedCard.Code
		}

		v.cards[card.Code] = card
	}
	// handle duplicates
	for _, dup := range dups {
		// get the original card
		oCard := v.cards[dup.DuplicateOf]
		if oCard == nil {
			return fmt.Errorf("could not find duplicate card")
		}
		oCard.PackCodes = append(oCard.PackCodes, dup.PackCode)
		oCard.DuplicateBy = append(oCard.DuplicateBy, dup.Code)
		v.cards[dup.Code] = oCard // point to the same card
	}

	// defer log.Println("Local card count:", mw.GetCollectionSize(v.db, cCards))

	return mw.CreateMany(v.db, cCards, v.getUniqueCards())
}

func (v *Valuator) updateHeroes() error {
	log.Println("Updating local list of heroes.")

	// get hero cards
	identityCards, err := mw.GetMany[Card](v.db, cCards, mw.BuildEqualsFilter("aspect", "hero"), mw.BsonNoneM)
	if err != nil {
		return err
	}

	// convert to heroes
	rawHeroes := []*Hero{}
	for _, heroCard := range identityCards {
		if heroCard.TypeCode == "hero" {
			hero := &Hero{
				Code:     heroCard.Code,
				Name:     heroCard.Name,
				PackCode: heroCard.PackCodes[0],
				Traits:   heroCard.Traits,
			}
			if heroCard.LinkedCardCode != "" {
				linkedCard := v.cards[heroCard.LinkedCardCode]
				if linkedCard == nil {
					return fmt.Errorf("could not find linked card")
				}
				hero.Merge(&Hero{
					Code:   linkedCard.Code,
					Traits: linkedCard.Traits,
				})
			}
			rawHeroes = append(rawHeroes, hero)
		}
	}

	// merge same name & same pack heroes
	// this covers cases like "Ironheart" who has more than 1 hero card
	heroesBy := map[string]*Hero{}
	for _, hero := range rawHeroes {
		u := hero.Name + hero.PackCode
		if _, ok := heroesBy[u]; !ok {
			heroesBy[u] = hero
		} else {
			heroesBy[u].Merge(hero)
		}
	}

	heroes := []*Hero{}
	for _, hero := range heroesBy {
		// add possible granted traits from hero cards to hero
		heroCards, err := mw.GetMany[Card](v.db, cCards, mw.BuildEqualsFilter("cardsetname", hero.Name), mw.BsonNoneM)
		if err != nil {
			return err
		}

		for _, heroCard := range heroCards {
			grantedTrait := parseGrantedTrait(heroCard)
			if grantedTrait != "" {
				hero.Traits = append(hero.Traits, grantedTrait)
			}
		}
		hero.SanitizeTraits()
		heroes = append(heroes, hero)
	}

	// defer log.Println("Local hero count:", mw.GetCollectionSize(v.db, cHeroes))

	return mw.CreateMany(v.db, cHeroes, heroes)
}

func (v *Valuator) updateDecks() (isNewDecks bool, err error) {
	log.Println("Updating local list of decks.")

	isNewDecks = false

	// default latest time to our starting date
	latestTime, _ := time.Parse("2006-01-02", os.Getenv("DECKLISTS_FROM_TIME"))

	// get the latest deck we have stored and update latest time if needed
	deck, err := mw.GetOne[marvel.Decklist](v.db, cDecks, mw.BsonNoneD, bson.M{"datecreatedstr": -1})
	if err != nil {
		return false, err
	}
	if deck != nil {
		latestTime = deck.DateCreated().Add(time.Hour * 24)
	}

	// get decks since latest time
	for {
		if latestTime.After(time.Now()) {
			break
		}

		decks, err := v.mCli.GetDecklists(latestTime)
		if err != nil {
			// marvelcdb api seems to return a 500 if there are simply no decks, just try the next day
			if !strings.Contains(err.Error(), "500 Internal Server Error") {
				return false, err
			}
		}
		log.Println("Adding more decks:", len(decks))

		if len(decks) > 0 {
			isNewDecks = true
			if err := mw.CreateMany(v.db, cDecks, decks); err != nil {
				return false, err
			}
		}

		latestTime = latestTime.Add(time.Hour * 24)
	}

	// defer log.Println("Local deck count:", mw.GetCollectionSize(v.db, cDecks))

	return isNewDecks, nil
}

func (v *Valuator) updateCardValues() error {
	log.Println("Updating local list of base card values.")

	// get all decks
	allDecks, err := mw.GetAll[marvel.Decklist](v.db, cDecks)
	if err != nil {
		return err
	}

	// get all heroes
	allHeroes, err := mw.GetAll[Hero](v.db, cHeroes)
	if err != nil {
		return err
	}

	// prepare hero map
	heroesByCode := map[string]*Hero{}
	for _, hero := range allHeroes {
		heroesByCode[hero.Code[:len(hero.Code)-1]] = hero
	}

	// get all cards
	allCards := v.getUniqueCards()

	// loop through every card and check if each deck is eligible to run that card or not and if it does
	cardValues := []*CardValue{}
	for _, card := range allCards {
		cardValue := &CardValue{
			Code:      card.Code,
			Card:      card,
			NewMod:    1,
			WeightMod: 1,
		}

		for _, deck := range allDecks {
			hero := heroesByCode[deck.HeroCode[:len(deck.HeroCode)-1]]
			if hero == nil {
				return fmt.Errorf("could not find hero from decklist")
			}

			if isCardEligibleForDeck(card, deck, hero) {
				cardValue.EligibleDecksCount += 1

				// check if card (or duplicates) are in the deck
				toCheck := []string{card.Code}
				toCheck = append(toCheck, card.DuplicateBy...)
				inUse := false
				for _, code := range toCheck {
					if count, ok := deck.Slots[code]; ok && count > 0 {
						inUse = true
						break
					}
				}
				if inUse {
					cardValue.InDecksCount += 1
				}
			}
		}

		cardValue.Calculate()
		cardValues = append(cardValues, cardValue)
	}

	sort.Slice(cardValues, func(i, j int) bool { return cardValues[i].Value > cardValues[j].Value })

	// defer log.Println("Local card values count:", mw.GetCollectionSize(v.db, cCardValues))

	return mw.CreateMany(v.db, cCardValues, cardValues)
}

func (v *Valuator) updatePackValues() error {
	log.Println("Updating local list of base pack values.")

	// get packs
	allPacks, err := mw.GetAll[marvel.Pack](v.db, cPacks)
	if err != nil {
		return err
	}

	packValues := make([]*PackValue, 0)
	for _, pack := range allPacks {
		// get cards in pack
		packCards, err := v.getCardsFromPack(pack.Code)
		if err != nil {
			return err
		}

		// skip "empty" packs (likely scenario packs)
		if len(packCards) == 0 {
			continue
		}

		// get card values for those cards
		cardCodes := make([]string, len(packCards))
		for i, pCard := range packCards {
			cardCodes[i] = pCard.Code
		}
		filter := bson.D{{Key: "_id", Value: bson.D{{"$in", cardCodes}}}}
		cvs, err := mw.GetMany[CardValue](v.db, cCardValues, filter, bson.M{"Value": -1})
		if err != nil {
			return err
		}

		packValue := &PackValue{
			Code:       pack.Code,
			Pack:       pack,
			CardValues: cvs,
		}
		packValue.Calculate()
		packValues = append(packValues, packValue)
	}

	sort.Slice(packValues, func(i, j int) bool { return packValues[i].ValueSum > packValues[j].ValueSum })

	// defer log.Println("Local pack values count:", mw.GetCollectionSize(v.db, cPackValues))

	return mw.CreateMany(v.db, cPackValues, packValues)
}

func (v *Valuator) getCardsFromPack(packCode string) ([]*Card, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"packcodes", packCode}},
				bson.D{{"aspect", bson.D{{"$in", []string{"basic", "justice", "protection", "aggression", "leadership"}}}}},
			},
		},
	}

	return mw.GetMany[Card](v.db, cCards, filter, mw.BsonNoneM)
}

func (v *Valuator) getCardsFromPacks(packCodes []string) (map[string]*Card, error) {
	allCards := map[string]*Card{}
	for _, packCode := range packCodes {
		cards, err := v.getCardsFromPack(packCode)
		if err != nil {
			return nil, err
		}
		for _, card := range cards {
			allCards[card.Code] = card
		}
	}
	return allCards, nil
}

func (v *Valuator) getUniqueCards() []*Card {
	cards := make([]*Card, 0)
	for oCode, card := range v.cards {
		if card.Code == oCode { // avoid duplicates
			cards = append(cards, card)
		}
	}
	return cards
}

func isCardEligibleForDeck(card *Card, deck *marvel.Decklist, hero *Hero) bool {
	// not eligible if the deck was made before the card released
	if deck.DateUpdated().Before(card.DateAvailable) {
		return false
	}

	// card aspect needs to be basic or match the deck
	if card.Aspect != "basic" && !utils.StringsContains(deck.Aspects(), card.Aspect) {
		return false
	}

	// does the card name match the heroes name
	// TODO

	// does the card have a locking trait
	// what hero is the deck running
	// does that hero have the locking trait
	if len(card.LockingTraits) > 0 {
		for _, trait := range card.LockingTraits {
			if utils.StringsContains(hero.Traits, trait) {
				return true
			}
		}
		return false
	}

	return true
}

func adjustCardValue(cv *CardValue, ownedCards map[string]*Card, ownedHeroes map[string]*Hero, allHeroes []*Hero, packCode string, aspectWeights map[string]float64) {
	// owned cards
	if _, ok := ownedCards[cv.Code]; ok {
		cv.NewMod = 0
	}

	// aspect weight
	if weight, ok := aspectWeights[cv.Card.Aspect]; ok {
		cv.WeightMod = weight
	}

	// trait-locked cards
	if cv.Card.LockingTraits != nil && len(cv.Card.LockingTraits) > 0 {
		// how many total heroes are there
		packHeroes := map[string]*Hero{}
		traitedHeroes := []*Hero{}
		for _, hero := range allHeroes {
			for _, trait := range cv.Card.LockingTraits {
				if hero.PackCode == packCode {
					packHeroes[hero.Code] = hero
				}
				if utils.StringsContains(hero.Traits, trait) {
					traitedHeroes = append(traitedHeroes, hero)
					break
				}
			}
		}
		cv.EligibleHeroCount = len(ownedHeroes)

		// how many owned
		count := 0
		for _, hero := range traitedHeroes {
			if _, ok := ownedHeroes[hero.Code]; ok {
				count++
			} else if _, ok := packHeroes[hero.Code]; ok {
				count++
			}
		}
		cv.OwnedHeroCount = count
	}

	cv.Calculate()
}

func parseGrantedTrait(card *Card) string {
	// "gain the BLANK trait"
	r1 := regexp.MustCompile(`gain the (.+) trait`)

	// "gains the BLANK trait"
	r2 := regexp.MustCompile(`gains the (.+) trait`)

	matches := r1.FindStringSubmatch(card.Text)
	if len(matches) > 1 {
		return strings.Trim(matches[1], "[].")
	}
	matches = r2.FindStringSubmatch(card.Text)
	if len(matches) > 1 {
		return strings.Trim(matches[1], "[].")
	}

	return ""
}

func parseLockingTraits(text string) []string {
	// "Play only if your identity has the BLANK or BLANK trait"
	r := regexp.MustCompile(`Play only if your identity has the (.+) or (.+) trait`)
	matches := r.FindStringSubmatch(text)
	if len(matches) > 2 {
		return []string{strings.Trim(matches[1], "[]."), strings.Trim(matches[2], "[].")}
	}

	// "Play only if your identity has the BLANK trait"
	r = regexp.MustCompile(`Play only if your identity has the (.+) trait`)
	matches = r.FindStringSubmatch(text)
	if len(matches) > 1 {
		return []string{strings.Trim(matches[1], "[].")}
	}

	// "Play only if you have the BLANK trait"
	r = regexp.MustCompile(`Play only if you have the (.+) trait`)
	matches = r.FindStringSubmatch(text)
	if len(matches) > 1 {
		return []string{strings.Trim(matches[1], "[].")}
	}

	return []string{}
}
