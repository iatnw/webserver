package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirp, err := apiConfig.db.GetChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (apiConfig *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := apiConfig.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not list chirps", err)
		return
	}
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
