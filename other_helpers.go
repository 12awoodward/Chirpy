package main

import "strings"

func filterWords(msg string) string {
	toRemove := map[string]struct{}{
		"kerfuffle": struct{}{},
		"sharbert": struct{}{},
		"fornax": struct{}{},
	}
	toReplace := "****"
	words := strings.Split(msg, " ")
	filteredText := make([]string, len(words))

	for i, word := range words {
		if _, ok := toRemove[strings.ToLower(word)]; ok {
			filteredText[i] = toReplace
		} else {
			filteredText[i] = word
		}
	}

	return strings.Join(filteredText, " ")
}