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
	cPackValues   = "pack-values"
	cStartingDate = "2020-01-01"
)

type Valuator struct {
	mCli *marvel.MarvelClient
	db   *mw.MongoDB
}

func NewValuator() *Valuator {
	v := &Valuator{}

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

	// get latest cards
	cards, err := v.mCli.GetAllCards()
	if err != nil {
		return err
	}
	if err := mw.CreateMany(v.db, cCards, cards); err != nil {
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

	// update pack values
	if addedDecks || true {
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

func (v *Valuator) ValueCard(cardCode string) (*CardValue, error) {
	// get all the decks
	decks, err := mw.GetMany[marvel.Decklist](v.db, cDecks, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// get card
	card, err := mw.GetOne[marvel.Card](v.db, cCards, bson.D{{"_id", cardCode}}, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, fmt.Errorf("card not found")
	}

	return calculateCardValue(decks, card), nil
}

func (v *Valuator) ValuePack(packCode string) (*PackValue, error) {
	// get all the decks
	decks, err := mw.GetMany[marvel.Decklist](v.db, cDecks, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// get cards
	cards, err := v.getCardsFromPack(packCode)

	// valuate
	cvs := calculateCardValues(decks, cards)
	pv := &PackValue{Code: packCode, CardValues: cvs}
	pv.Calculate()

	return pv, nil
}

func (v *Valuator) ValueAllPacks(owned []string) ([]*PackValue, error) {
	// grab base pack values from db
	pvs, err := mw.GetMany[PackValue](v.db, cPackValues, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	// make map of owned cards
	ownedCards := make(map[string]*marvel.Card)
	for _, code := range owned {
		cards, err := v.getCardsFromPack(code)
		if err != nil {
			return nil, err
		}

		for _, card := range cards {
			ownedCards[card.Code] = card

			// handle duplicates
			dups, err := v.GetDuplicateCodes(card)
			if err != nil {
				return nil, err
			}
			for _, dup := range dups {
				ownedCards[dup] = card
			}
		}
	}

	// modify pack values based owned cards
	for _, pv := range pvs {
		for _, cv := range pv.CardValues {
			if _, ok := ownedCards[cv.Code]; ok {
				cv.NewMod = 0
				cv.Calculate()
			}
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

func (v *Valuator) GetNewestDeck() (*marvel.Decklist, error) {
	return mw.GetOne[marvel.Decklist](v.db, cDecks, mw.BsonNoneD, bson.M{"datecreatedstr": -1})
}

func (v *Valuator) calculateAllPackValues() ([]*PackValue, error) {
	// get packs
	packs, err := mw.GetMany[marvel.Pack](v.db, cPacks, mw.BsonNoneD, mw.BsonNoneM)
	if err != nil {
		return nil, err
	}

	pvs := make([]*PackValue, 0)
	for _, pack := range packs {
		pv, _ := v.ValuePack(pack.Code)
		if pv.ValueSum != 0 {
			pvs = append(pvs, pv)
		}
	}

	sort.Slice(pvs, func(i, j int) bool { return pvs[i].ValueSum > pvs[j].ValueSum })

	return pvs, nil
}

func (v *Valuator) getCardsFromPack(packCode string) ([]*marvel.Card, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"packcode", bson.D{{"$eq", packCode}}}},
				bson.D{{"factioncode", bson.D{{"$ne", "hero"}}}},
			},
		},
	}

	return mw.GetMany[marvel.Card](v.db, cCards, filter, mw.BsonNoneM)
}

func (v *Valuator) GetDuplicateCodes(card *marvel.Card) ([]string, error) {
	if card.DuplicateBy != nil {
		return card.DuplicateBy, nil
	}

	if card.DuplicateOf != "" {
		// get the original card
		oCard, err := mw.GetOne[marvel.Card](v.db, cCards, bson.D{{"_id", card.DuplicateOf}}, mw.BsonNoneM)
		if err != nil {
			return nil, err
		}
		return append(oCard.DuplicateBy, card.DuplicateOf), nil
	}

	return []string{}, nil
}

func calculateCardValue(decks []*marvel.Decklist, card *marvel.Card) *CardValue {
	cv := &CardValue{
		Code:               card.Code,
		Name:               card.Name,
		NewMod:             1,
		EligableDecksCount: 0,
		InDecksCount:       0,
	}
	for _, deck := range decks {
		eligable := false
		if card.FactionCode == "basic" {
			eligable = true
		} else {
			if utils.SliceContains(deck.Aspects(), card.FactionCode) {
				eligable = true
			}
		}

		if eligable {
			cv.EligableDecksCount += 1
			if count, ok := deck.Slots[card.Code]; ok && count > 0 {
				cv.InDecksCount += 1
			}
		}
	}
	cv.Calculate()
	return cv
}

func calculateCardValues(decks []*marvel.Decklist, cards []*marvel.Card) []*CardValue {
	cvs := make([]*CardValue, len(cards))
	for i, card := range cards {
		cvs[i] = &CardValue{
			Code:   card.Code,
			Name:   card.Name,
			NewMod: 1,
		}
	}
	for _, deck := range decks {
		for i, cv := range cvs {
			card := cards[i]
			eligable := false
			if card.FactionCode == "basic" {
				eligable = true
			} else {
				if utils.SliceContains(deck.Aspects(), card.FactionCode) {
					eligable = true
				}
			}

			if eligable {
				cv.EligableDecksCount += 1
				if count, ok := deck.Slots[card.Code]; ok && count > 0 {
					cv.InDecksCount += 1
				}
			}
		}
	}
	for _, cv := range cvs {
		cv.Calculate()
	}
	sort.Slice(cvs, func(i, j int) bool { return cvs[i].Value > cvs[j].Value })
	return cvs
}
