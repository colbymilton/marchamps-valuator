package controller

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	marvel "github.com/colbymilton/marchamps-valuator/internal/marvelcdb"
	"github.com/colbymilton/marchamps-valuator/internal/utils"
	mw "github.com/colbymilton/marchamps-valuator/pkg/mongoWrapper"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	cPacks        = "packs"
	cCards        = "cards"
	cDecks        = "decks"
	cCardValues   = "card-values"
	cPackValues   = "pack-values"
	cStartingDate = "2020-01-01"
)

type Valuator struct {
	mCli  *marvel.MarvelClient
	db    *mw.MongoDB
	cards map[string]*Card
}

func NewValuator() *Valuator {
	v := &Valuator{cards: make(map[string]*Card)}

	mcli, err := marvel.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	v.mCli = mcli
	v.db = mw.NewMongoDB("mongodb://localhost:27017", "marchamps-valuator")
	return v
}

func (v *Valuator) UpdateDatabase() error {
	// get latest packs
	packs, err := v.mCli.GetAllPacks()
	if err != nil {
		return err
	}
	if err := mw.CreateMany(v.db, cPacks, packs); err != nil {
		return err
	}
	// build pack map for future use when convert cards
	packMap := make(map[string]*marvel.Pack)
	for _, pack := range packs {
		packMap[pack.Code] = pack
	}

	// get latest cards
	mCards, err := v.mCli.GetAllCards()
	if err != nil {
		return err
	}

	// convert cards to local
	dups := make([]*marvel.Card, 0)
	for _, mCard := range mCards {
		if mCard.FactionCode == "hero" {
			continue // skip hero cards
		}

		if mCard.DuplicateOf != "" {
			dups = append(dups, mCard)
			continue // handle dups later
		}

		card := &Card{
			Code:          mCard.Code,
			Name:          mCard.Name,
			PackCodes:     []string{mCard.PackCode},
			TypeCode:      mCard.TypeCode,
			Aspect:        mCard.FactionCode,
			Traits:        strings.Split(mCard.Traits, ". "),
			LockingTrait:  "", // TODO
			DateAvailable: packMap[mCard.PackCode].DateAvailable(),
			DuplicateBy:   []string{},
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
	if err := mw.CreateMany(v.db, cCards, v.getCards()); err != nil {
		return err
	}

	// get latest decks
	// when was the last deck stored?
	addedDecks := false
	lt, _ := time.Parse("2006-01-02", cStartingDate)
	deck, err := mw.GetOne[marvel.Decklist](v.db, cDecks, mw.BsonNoneD, bson.M{"datecreatedstr": -1})
	if err != nil {
		return err
	}
	if deck != nil {
		lt = deck.DateCreated().Add(time.Hour * 24)
	}
	// get decks since then
	for {
		if lt.After(time.Now()) {
			break
		}

		decks, err := v.mCli.GetDecklists(lt)
		if err != nil {
			// marvelcdb api seems to return a 500 if there are simply no decks, just try the next day
			if !strings.Contains(err.Error(), "500 Internal Server Error") {
				return err
			}
		}
		log.Println("Adding more decks:", len(decks))

		if len(decks) > 0 {
			addedDecks = true
			if err := mw.CreateMany(v.db, cDecks, decks); err != nil {
				return err
			}
		}

		lt = lt.Add(time.Hour * 24)
	}

	// update values of cards and packs
	if addedDecks || true {
		// get all decks (for value calculation below)
		allDecks, err := mw.GetMany[marvel.Decklist](v.db, cDecks, mw.BsonNoneD, mw.BsonNoneM)
		if err != nil {
			return err
		}

		// get the base value of every card
		baseCVs := v.calculateCardValues(allDecks, v.getCards())
		if err := mw.CreateMany(v.db, cCardValues, baseCVs); err != nil {
			return err
		}

		// get the pase value of every pack
		basePVs, err := v.calculateAllPackValues()
		if err != nil {
			return err
		}
		if err := mw.CreateMany(v.db, cPackValues, basePVs); err != nil {
			return err
		}
	}

	return nil
}

func (v *Valuator) ValueAllCards(owned []string) ([]*CardValue, error) {
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

	// modify base pack values based on owned cards
	for _, cv := range cvs {
		adjustCardValue(cv, ownedCards)
	}

	sort.Slice(cvs, func(i, j int) bool { return cvs[i].Value > cvs[j].Value })

	return cvs, nil
}

func (v *Valuator) ValueAllPacks(owned []string) ([]*PackValue, error) {
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

	// modify base pack values based on owned cards
	for _, pv := range pvs {
		for _, cv := range pv.CardValues {
			adjustCardValue(cv, ownedCards)
		}
		sort.Slice(pv.CardValues, func(i, j int) bool { return pv.CardValues[i].Value > pv.CardValues[j].Value })
		pv.Calculate()
	}

	sort.Slice(pvs, func(i, j int) bool { return pvs[i].ValueSum > pvs[j].ValueSum })

	return pvs, nil
}

func (v *Valuator) GetPacks() ([]*marvel.Pack, error) {
	return mw.GetMany[marvel.Pack](v.db, cPacks, mw.BsonNoneD, bson.M{"availablestr": 1})
}

func (v *Valuator) getCardsFromPack(packCode string) ([]*Card, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"packcodes", packCode}},
				bson.D{{"factioncode", bson.D{{"$ne", "hero"}}}},
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

func (v *Valuator) calculateAllPackValues() ([]*PackValue, error) {
	// get packs
	packs, err := mw.GetMany[marvel.Pack](v.db, cPacks, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	pvs := make([]*PackValue, 0)
	for _, pack := range packs {
		// get cards in pack
		pCards, err := v.getCardsFromPack(pack.Code)
		if err != nil {
			return nil, err
		}

		// get card values for those cards
		cardCodes := make([]string, len(pCards))
		for i, pCard := range pCards {
			cardCodes[i] = pCard.Code
		}
		filter := bson.D{{"_id", bson.D{{"$in", cardCodes}}}}
		cvs, err := mw.GetMany[CardValue](v.db, cCardValues, filter, bson.M{"Value": -1})
		if err != nil {
			return nil, err
		}
		sort.Slice(cvs, func(i, j int) bool { return cvs[i].Value > cvs[j].Value })

		pv := &PackValue{
			Code:       pack.Code,
			CardValues: cvs,
		}
		pv.Calculate()
		pvs = append(pvs, pv)
	}

	sort.Slice(pvs, func(i, j int) bool { return pvs[i].ValueSum > pvs[j].ValueSum })

	return pvs, nil
}

func (v *Valuator) calculateCardValues(decks []*marvel.Decklist, cards []*Card) []*CardValue {
	// prepare card values
	cvs := make([]*CardValue, len(cards))
	for i, card := range cards {
		cvs[i] = &CardValue{
			Code:   card.Code,
			Name:   card.Name,
			NewMod: 1,
		}
	}

	// loop through every deck and check if each card (or duplicate versions) are used
	for _, deck := range decks {
		deckTime := deck.DateUpdated()
		for i, cv := range cvs {
			card := cards[i]

			if deckTime.Before(card.DateAvailable) {
				continue // not eligable if the deck was made before the card released
			}

			if card.Aspect != "basic" && !utils.SliceContains(deck.Aspects(), card.Aspect) {
				continue // card aspect needs to be basic or match the deck
			}

			cv.EligableDecksCount += 1

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
				cv.InDecksCount += 1
			}
		}
	}
	for _, cv := range cvs {
		cv.Calculate()
	}
	sort.Slice(cvs, func(i, j int) bool { return cvs[i].Value > cvs[j].Value })
	return cvs
}

func (v *Valuator) getCards() []*Card {
	cards := make([]*Card, 0)
	for oCode, card := range v.cards {
		if card.Code == oCode { // avoid duplicates
			cards = append(cards, card)
		}
	}
	return cards
}

func adjustCardValue(cv *CardValue, ownedCards map[string]*Card) {
	if _, ok := ownedCards[cv.Code]; ok {
		cv.NewMod = 0
		cv.Calculate()
	}
}
