package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"webserver/internal/auth"

	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != apiConfig.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not read request", err)
		return
	}
	var params parameters
	if err := json.Unmarshal(body, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	err = apiConfig.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
