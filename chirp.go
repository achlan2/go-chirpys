package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	type returnVal struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	re := regexp.MustCompile(`\b(?i)(kerfuffle|sharbert|fornax)\b`)

	respondWithJSON(w, 200, returnVal{Valid: true, CleanedBody: re.ReplaceAllString(params.Body, "****")})
}
