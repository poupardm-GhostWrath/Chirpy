package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/poupardm-GhostWrath/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 				string 	`json:"email"`
		Password 			string 	`json:"password"`
		ExpiresInSeconds	int		`json:"expires_in_seconds"`
	}

	// Decode Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Get User from DB
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Check Password
	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Set Time Duration
	var expireTime time.Duration = time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds <= 3600 {
		expireTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	// Create Token
	token, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, expireTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	// Respond with User
	respondWithJSON(w, http.StatusOK, User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		Token: token,
	})
}

