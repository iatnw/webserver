package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

var profaneRe = regexp.MustCompile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "could not read request", err)
		return
	}
	params := parameters{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, 500, "could not unmarshal parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
	} else {
		respondWithJSON(w, http.StatusOK, responseBody{
			CleanedBody: replaceProfaneWords(params.Body),
		})
	}
}
func replaceProfaneWords(originalText string) string {
	return profaneRe.ReplaceAllString(originalText, "****")
}
