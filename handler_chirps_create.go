package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"webserver/internal/auth"
	"webserver/internal/database"
)

var profaneRe = regexp.MustCompile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)

func (apiConfig *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, apiConfig.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
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

	dbChirp, err := apiConfig.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      cleanedBody,
			UserID:    userID,
		},
	})
}

func replaceProfaneWords(originalText string) string {
	return profaneRe.ReplaceAllString(originalText, "****")
}
