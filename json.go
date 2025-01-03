package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		log.Println(err)
	}

	if status > 499 {
		log.Printf("Responding with 5XX error: %s", message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, status, errorResponse{
		Error: message,
	})
}

func respondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(status)
	w.Write(dat)
}
