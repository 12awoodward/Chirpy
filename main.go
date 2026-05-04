package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", readinessEndpoint)
	mux.HandleFunc("POST /api/validate_chirp", validateChripEndpoint)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsEndpoint)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsEndpoint)

	err := server.ListenAndServe()
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
		err = respondWithError(w, http.StatusBadRequest, "Request should be json including 'body' string")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if len(chirp.Body) > 140 {
		err = respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		if err != nil {
			fmt.Println(err)
		}
		
	} else {
		response := validResponse{CleanedBody: filterWords(chirp.Body)}
		respondWithJSON(w, http.StatusOK, response)
	}
}
