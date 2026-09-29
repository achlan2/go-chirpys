package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/achlan2/go-chirpy/internal/auth"
	"github.com/achlan2/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type response struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("id")

	id, parseErr := uuid.Parse(chirpId)
	if parseErr != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse uuid", parseErr)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}
	result := []response{}
	for _, chirp := range chirps {
		result = append(result, response{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, result)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	token, errBearer := auth.GetBearerToken(r.Header)
	if errBearer != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer Unauthorized!", nil)
		return
	}

	userId, errJWT := auth.ValidateJWT(token, cfg.secret)
	if userId == uuid.Nil && errJWT != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized!", nil)
		return
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

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   re.ReplaceAllString(params.Body, "****"),
		UserID: userId,
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
