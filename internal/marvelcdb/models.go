package marvel

import (
	"strings"
	"time"
)

type Decklist struct {
	Id             int            `json:"id" bson:"_id"`
	DateCreatedStr string         `json:"date_creation"`
	DateUpdatedStr string         `json:"date_update"`
	Slots          map[string]int `json:"slots"`
	Meta           string         `json:"meta"`
	HeroCode       string         `json:"investigator_code"`
}

func (d *Decklist) DateCreated() time.Time {
	t, _ := time.Parse(time.RFC3339, d.DateCreatedStr)
	return t
}

func (d *Decklist) DateUpdated() time.Time {
	t, _ := time.Parse(time.RFC3339, d.DateUpdatedStr)
	return t
}

func (d *Decklist) Aspects() []string {
	a := make([]string, 0)
	if strings.Contains(d.Meta, "aggression") {
		a = append(a, "aggression")
	}
	if strings.Contains(d.Meta, "justice") {
		a = append(a, "justice")
	}
	if strings.Contains(d.Meta, "leadership") {
		a = append(a, "leadership")
	}
	if strings.Contains(d.Meta, "protection") {
		a = append(a, "protection")
	}
	return a
}

type Card struct {
	Code        string   `json:"code" bson:"_id"`
	Name        string   `json:"name"`
	SubName     string   `json:"subname"`
	PackCode    string   `json:"pack_code"`
	TypeCode    string   `json:"type_code"`
	FactionCode string   `json:"faction_code"`
	Traits      string   `json:"traits"`
	DuplicateOf string   `json:"duplicate_of_code"`
	DuplicateBy []string `json:"duplicated_by"`
	Text        string   `json:"text"`
	CardSetName string   `json:"card_set_name"`
	LinkedCard  *Card    `json:"linked_card"`
}

type Pack struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Id           int    `json:"id" bson:"_id"`
	AvailableStr string `json:"available"`
}

func (p *Pack) DateAvailable() time.Time {
	t, _ := time.Parse("2006-01-02", p.AvailableStr)
	return t
}
