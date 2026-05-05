package main

import (
	"encoding/json"
	"net/http"

	"github.com/12awoodward/chirpy/internal/auth"
	"github.com/12awoodward/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginPostEndpoint(w http.ResponseWriter, r *http.Request) {
	type loginUserJSON struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	toLogin := loginUserJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&toLogin)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	user, err := cfg.db.GetUser(r.Context(), toLogin.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	canLogin, err := auth.CheckPasswordHash(toLogin.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if !canLogin {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, toJSONUser(user))
}

func (cfg *apiConfig) usersPostEndpoint(w http.ResponseWriter, r *http.Request) {
	type newUserJSON struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	newUser := newUserJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "New user requires email")
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to add user")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: newUser.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to add user")
		return
	}

	respondWithJSON(w, http.StatusCreated, toJSONUser(user))
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

	respondWithJSON(w, http.StatusCreated, toJSONChirp(chirp))
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

	respondWithJSON(w, http.StatusOK, toJSONChirp(chirp))
}
