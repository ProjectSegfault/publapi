package utils

import (
	"strings"
)

func Dedup(input string) string {
	unique := []string{}
	words := strings.Split(input, " ")
	for _, word := range words {
		if contains(unique, word) {
			continue
		}
		unique = append(unique, word)
	}
	return strings.Join(unique, " ")
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
