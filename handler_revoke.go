package main

import (
	"net/http"

	"github.com/poupardm-GhostWrath/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeTokens(w http.ResponseWriter, r *http.Request) {
	// Check Refresh Token Header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Revoke Token in DB
	_, err = cfg.db.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke Token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}