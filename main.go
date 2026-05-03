package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", readinessEndpoint)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func readinessEndpoint(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
}
