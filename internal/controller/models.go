package controller

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	marvel "github.com/colbymilton/marchamps-valuator/internal/marvelcdb"
	"github.com/colbymilton/marchamps-valuator/internal/utils"
)

type Card struct {
	Code           string    `json:"code" bson:"_id"`
	Name           string    `json:"name"`
	PackCodes      []string  `json:"packCodes"`
	TypeCode       string    `json:"typeCode"`
	Aspect         string    `json:"aspect"`
	Traits         []string  `json:"traits"`
	LockingTraits  []string  `json:"lockingTraits"`
	DateAvailable  time.Time `json:"dateAvailable"`
	DuplicateBy    []string  `json:"duplicatedBy"`
	Text           string    `json:"text"`
	CardSetName    string    `json:"cardSetName"`
	LinkedCardCode string    `json:"linkedCard"`
}

type CardValue struct {
	Code               string  `json:"code" bson:"_id"`
	Card               *Card   `json:"card"`
	Value              int     `json:"value"`
	NewMod             float64 `json:"newMod"`
	PopularityMod      float64 `json:"popularityMod"`
	EligibleDecksCount int     `json:"eligibleDecksCount"`
	InDecksCount       int     `json:"inDecksCount"`
	TraitMod           float64 `json:"traitMod"`
	EligibleHeroCount  int     `json:"eligibleHeroCount"`
	OwnedHeroCount     int     `json:"ownedHeroCount"`
}

func (cv *CardValue) Calculate() {
	cv.PopularityMod = 1
	if cv.EligibleDecksCount > 0 {
		cv.PopularityMod += float64(cv.InDecksCount) / float64(cv.EligibleDecksCount)
	}
	cv.TraitMod = 1
	if cv.EligibleHeroCount > 0 {
		cv.TraitMod = float64(cv.OwnedHeroCount) / float64(cv.EligibleHeroCount)
	}
	cv.Value = int(math.Round(100 * cv.NewMod * cv.PopularityMod * cv.TraitMod))
}

type PackValue struct {
	Code       string       `json:"code" bson:"_id"`
	Pack       *marvel.Pack `json:"pack"`
	ValueSum   int          `json:"valueSum"`
	CardValues []*CardValue `json:"cardValues"`
}

func (pv *PackValue) Calculate() {
	pv.ValueSum = 0
	for _, cv := range pv.CardValues {
		pv.ValueSum += cv.Value
	}
}

type Hero struct {
	Code     string   `json:"code" bson:"_id"`
	PackCode string   `json:"packCode"`
	Name     string   `json:"name"`
	Traits   []string `json:"traits"`
}

func (h *Hero) Merge(h2 *Hero) {
	v1, _ := strconv.Atoi(h.Code[:len(h.Code)-1])
	v2, _ := strconv.Atoi(h2.Code[:len(h2.Code)-1])
	if v1 > 0 && v2 > 0 {
		h.Code = fmt.Sprintf("%05da", int(math.Min(float64(v1), float64(v2))))
	}

	h.Traits = append(h.Traits, h2.Traits...)
	h.SanitizeTraits()
}

func (h *Hero) SanitizeTraits() {
	newTraits := []string{}
	for _, trait := range h.Traits {
		t := strings.Trim(trait, ".")
		if !utils.SliceContains(newTraits, t) {
			newTraits = append(newTraits, t)
		}
	}

	h.Traits = newTraits
}
