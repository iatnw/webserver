package main

import (
	"net/http"

	"webserver/internal/auth"
)

func (apiConfig *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	err = apiConfig.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
