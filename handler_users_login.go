package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"webserver/internal/auth"
)

func (apiConfig *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var params parameters
	if err := json.Unmarshal(data, &params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := apiConfig.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}
	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	const maxExpiry = time.Hour
	expiry := maxExpiry
	if params.ExpiresInSeconds != nil {
		requested := time.Duration(*params.ExpiresInSeconds) * time.Second
		if requested < expiry {
			expiry = requested
		}
	}

	token, err := auth.MakeJWT(user.ID, os.Getenv("JWT_SECRET"), expiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
