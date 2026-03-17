package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"webserver/internal/database"

	"github.com/google/uuid"
)

var profaneRe = regexp.MustCompile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not read request", err)
		return
	}
	params := parameters{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not unmarshal parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleanedBody := replaceProfaneWords(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      cleanedBody,
			UserID:    params.UserID,
		},
	},
	)
}

func replaceProfaneWords(originalText string) string {
	return profaneRe.ReplaceAllString(originalText, "****")
}
