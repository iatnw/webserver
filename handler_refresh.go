package main

import (
	"net/http"
	"os"
	"time"

	"webserver/internal/auth"
)

func (apiConfig *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	dbToken, err := apiConfig.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	if time.Now().After(dbToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "token expired", nil)
		return
	}
	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, os.Getenv("JWT_SECRET"), time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}
