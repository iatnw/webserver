package main

import (
	"encoding/json"
	"io"
	"net/http"

	"webserver/internal/auth"
	"webserver/internal/database"
)

func (apiConfig *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not read request", err)
		return
	}
	var params parameters
	if err := json.Unmarshal(data, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	updated, err := apiConfig.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          updated.ID,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
		Email:       updated.Email,
		IsChirpyRed: updated.IsChirpyRed,
	})
}
