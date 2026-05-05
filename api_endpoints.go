package main

import (
	"encoding/json"
	"net/http"

	"github.com/12awoodward/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) usersEndpoint(w http.ResponseWriter, r * http.Request) {
	type emailJSON struct {
		Email string `json:"email"`
	}

	newUser := emailJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "New user requires email")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), newUser.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to add user")
		return
	}

	addedUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusCreated, addedUser)

}

func (cfg *apiConfig) chirpsEndpoint(w http.ResponseWriter, r *http.Request) {
	type newChirpJSON struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	newChirp := newChirpJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newChirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Expected chirp body and user id")
		return
	}

	if len(newChirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
		
	} 

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: filterWords(newChirp.Body),
		UserID: newChirp.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to add Chirp")
		return
	}

	addedChirp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, addedChirp)
}
