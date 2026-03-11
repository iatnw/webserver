package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		Valid bool `json:"valid"`
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "could not read request", err)
		return
	}
	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, 500, "could not unmarshal parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "", err)
	} else {
		respondWithJSON(w, http.StatusOK, responseBody{
			Valid: true,
		})
	}
}
