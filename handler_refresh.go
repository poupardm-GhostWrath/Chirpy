package main

import (
	"time"
	"net/http"

	"github.com/poupardm-GhostWrath/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshTokens(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Check Refresh Token Header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	// Get Refresh Token in DB
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	// Create JWT Token
	jwtToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	// Respond with JWT Token
	respondWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})
}