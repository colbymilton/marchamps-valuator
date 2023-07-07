package utils

import (
	"encoding/json"
	"log"
	"strings"
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

func StringsContains(ss []string, os string) bool {
	for _, s := range ss {
		if strings.EqualFold(s, os) {
			return true
		}
	}

	return false
}
