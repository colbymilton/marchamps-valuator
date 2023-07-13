package marvel

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

const baseAddr = "http://marvelcdb.com/api/public/"

type MarvelClient struct {
}

func NewClient() (*MarvelClient, error) {
	mcli := &MarvelClient{}
	return mcli, nil
}

func (mcli *MarvelClient) get(endpoint string, body any) error {
	log.Println("sending:", baseAddr+endpoint)
	resp, err := sling.New().Get(baseAddr + endpoint).ReceiveSuccess(body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response (status %v): %v", resp.StatusCode, resp.Status)
	}
	return nil
}

// GetDecklists returns all Decklists posted on Marvelcdb on the specified day
func (mcli *MarvelClient) GetDecklists(date time.Time) ([]*Decklist, error) {
	formatted := date.Format("2006-01-02")
	var decklists []*Decklist
	err := mcli.get("decklists/by_date/"+formatted, &decklists)
	return decklists, err
}

// GetAllCards returns all the cards on Marvelcdb
func (mcli *MarvelClient) GetAllCards() ([]*Card, error) {
	var cards []*Card
	err := mcli.get("cards", &cards)
	return cards, err
}

// GetCards returns all cards that are part of the specified pack
func (mcli *MarvelClient) GetCards(packCode string) ([]*Card, error) {
	var cards []*Card
	err := mcli.get("cards/"+packCode, &cards)
	return cards, err
}

// GetAllPacks returns all the packs on Marvelcdb
func (mcli *MarvelClient) GetAllPacks() ([]*Pack, error) {
	var packs []*Pack
	err := mcli.get("packs", &packs)
	return packs, err
}
