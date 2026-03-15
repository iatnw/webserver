package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Email string `json:"email"`
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
	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "could not create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated,
		User{
			Id:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
	)
}
