package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/achlan2/go-chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body", err)
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

	userID, parseErr := uuid.Parse(params.UserID)
	if parseErr != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse uuid", parseErr)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   re.ReplaceAllString(params.Body, "****"),
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
