package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/12awoodward/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", readinessEndpoint)
	mux.HandleFunc("POST /api/validate_chirp", validateChripEndpoint)
	mux.HandleFunc("POST /api/users", apiCfg.usersEndpoint)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsEndpoint)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetEndpoint)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChripEndpoint(w http.ResponseWriter, r *http.Request) {
	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type chirpBody struct {
		Body string `json:"body"`
	}

	chirp := chirpBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Request should be json including 'body' string")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		
	} else {
		response := validResponse{CleanedBody: filterWords(chirp.Body)}
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (cfg *apiConfig) usersEndpoint(w http.ResponseWriter, r * http.Request) {
	type emailJSON struct {
		Email string `json:"email"`
	}

	newUser := emailJSON{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "New user requires email.")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), newUser.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to add user.")
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
