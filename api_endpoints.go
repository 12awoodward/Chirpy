package main

import (
	"encoding/json"
	"net/http"

	"github.com/12awoodward/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) usersPostEndpoint(w http.ResponseWriter, r * http.Request) {
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

func (cfg *apiConfig) chirpsPostEndpoint(w http.ResponseWriter, r *http.Request) {
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

	addedChirp := toJSONChirp(chirp)
	respondWithJSON(w, http.StatusCreated, addedChirp)
}

func (cfg *apiConfig) chirpsGetEndpoint(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get chirps")
		return
	}

	chirps_response := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirps_response[i] = toJSONChirp(chirp)
	}

	respondWithJSON(w, http.StatusOK, chirps_response)
}

func (cfg *apiConfig) chirpsGetByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	uuid, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}
	
	chirp, err := cfg.db.GetChirp(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}

	foundChirp := toJSONChirp(chirp)
	respondWithJSON(w, http.StatusOK, foundChirp)
}
