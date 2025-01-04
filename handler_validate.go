package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const chirpMaxLength = 140
	if len(params.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	restrictedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(params.Body, restrictedWords)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}
