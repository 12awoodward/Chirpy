package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	type jsonErr struct {
		Error string `json:"error"`
	}

	respErr := jsonErr{Error: msg}
	respBody, err := json.Marshal(respErr)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)

	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	respBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)
	
	return nil
}