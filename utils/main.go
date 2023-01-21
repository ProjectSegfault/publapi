package utils

import (
	_ "embed"
	"strings"
)

//go:embed templates/signup-script-template
var Bashscript string

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
