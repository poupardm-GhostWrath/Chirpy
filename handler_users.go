package main

import (
	"time"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	// Decode Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Create User in DB
	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	// Create User From DBUser
	user := User{
		ID: 		dbUser.ID,
		CreatedAt: 	dbUser.CreatedAt,
		UpdatedAt:	dbUser.UpdatedAt,
		Email:		dbUser.Email,
	}

	// Respond with User
	respondWithJSON(w, http.StatusCreated, user)
}