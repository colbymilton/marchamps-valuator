package utils

import (
	"encoding/json"
	"log"
	"math"
)

func PrintJSON(i any) {
	if js, err := json.MarshalIndent(i, "", "  "); err != nil {
		log.Println("json error:", err)
	} else {
		log.Println(string(js))
	}
}

func SliceContains[T string](slice []T, find T) bool {
	for _, s := range slice {
		if s == find {
			return true
		}
	}

	return false
}

func RoundFloat3(in float64) float64 {
	return math.Round(in*1000) / 1000
}
