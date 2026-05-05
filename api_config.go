package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/12awoodward/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		},
	)
}

func (cfg *apiConfig) metricsEndpoint(w http.ResponseWriter, r *http.Request) {
	hitCount := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hitCount))
}

func (cfg *apiConfig) resetMetricsEndpoint(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset"))
}