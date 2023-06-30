package controller

import "github.com/colbymilton/marchamps-valuator/internal/utils"

type CardValue struct {
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Value              float64 `json:"value"`
	NewMod             float64 `json:"newMod"`
	PopularityMod      float64 `json:"popularityMod"`
	EligableDecksCount int     `json:"eligableDecksCount"`
	InDecksCount       int     `json:"inDecksCount"`
}

func (cv *CardValue) Calculate() {
	cv.PopularityMod = 1
	if cv.EligableDecksCount > 0 {
		cv.PopularityMod += utils.RoundFloat3(float64(cv.InDecksCount) / float64(cv.EligableDecksCount))
	}
	cv.Value = 1 * cv.NewMod * cv.PopularityMod
	cv.Value = utils.RoundFloat3(cv.Value)
}

type PackValue struct {
	Code       string       `json:"code" bson:"_id"`
	ValueSum   float64      `json:"valueSum"`
	CardValues []*CardValue `json:"cardValues"`
}

func (pv *PackValue) Calculate() {
	pv.ValueSum = 0
	for _, cv := range pv.CardValues {
		pv.ValueSum += cv.Value
	}
	pv.ValueSum = utils.RoundFloat3(pv.ValueSum)
}
