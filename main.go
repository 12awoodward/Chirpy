package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/12awoodward/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	secret string
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		},
	)
}

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
		secret: os.Getenv("SECRET"),
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", readinessEndpoint)

	mux.HandleFunc("POST /api/login", apiCfg.loginPostEndpoint)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshPostEndpoint)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokePostEndpoint)

	mux.HandleFunc("POST /api/users", apiCfg.usersPostEndpoint)
	mux.HandleFunc("PUT /api/users", apiCfg.usersPutEndpoint)

	mux.HandleFunc("GET /api/chirps", apiCfg.chirpsGetEndpoint)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.chirpsGetByIDEndpoint)
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpsPostEndpoint)

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

