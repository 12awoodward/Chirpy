package main

import (
	"context"
	"strings"

	"github.com/12awoodward/chirpy/internal/database"
	"github.com/google/uuid"
)

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

func getChirpSlice(cfg *apiConfig, ctx context.Context, authorQuery string) ([]database.Chirp, error) {
	if len(authorQuery) == 0 {
		chirps, err := cfg.db.GetChirps(ctx)
		if err != nil {
			return nil, err
		}
		return chirps, nil

	} else {
		authorID, err := uuid.Parse(authorQuery)
		if err != nil {
			return nil, err
		}

		chirps, err := cfg.db.GetChirpsByAuthor(ctx, authorID)
		if err != nil {
			return nil, err
		}
		return chirps, nil
	}
}
