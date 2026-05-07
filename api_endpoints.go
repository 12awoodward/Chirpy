package main

import (
	"encoding/json"
	"net/http"
	"time"

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

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to generate auth token")
		return
	}

	token_entry, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: auth.MakeRefreshToken(),
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(24 * 60)),
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to generate refresh token")
		return
	}

	type userResponse struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, userResponse{
		toJSONUser(user),
		token,
		token_entry.Token,
	})
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	type newChirpJSON struct {
		Body string `json:"body"`
	}

	newChirp := newChirpJSON{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newChirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Expected chirp body")
		return
	}

	if len(newChirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
		
	} 

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: filterWords(newChirp.Body),
		UserID: userID,
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

func (cfg *apiConfig) refreshPostEndpoint(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Expected token")
		return
	}
	
	token_entry, err := cfg.db.GetValidRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	token, err := auth.MakeJWT(token_entry.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to generate token")
		return
	}

	type newToken struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, newToken{token})
}

func (cfg *apiConfig) revokePostEndpoint(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Expected token")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
