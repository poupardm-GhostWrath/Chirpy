package main

import (
	"time"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/poupardm-GhostWrath/Chirpy/internal/auth"
	"github.com/poupardm-GhostWrath/Chirpy/internal/database"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
	Password	string		`json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// Decode Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Hash Password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Create User in DB
	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
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
