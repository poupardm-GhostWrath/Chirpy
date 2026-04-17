package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/poupardm-GhostWrath/Chirpy/internal/auth"
	"github.com/poupardm-GhostWrath/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 				string 	`json:"email"`
		Password 			string 	`json:"password"`
	}

	type response struct {
		User
		Token			string	`json:"token"`
		RefreshToken	string	`json:"refresh_token"`
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

	// Create JWT Token
	jwtToken, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	// Create Refresh Token
	refreshToken := auth.MakeRefreshToken()
	dbRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save Refresh Token", err)
		return
	}



	// Respond with User
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email: dbUser.Email,
		},
		Token: jwtToken,
		RefreshToken: dbRefreshToken.Token,
	})
}

