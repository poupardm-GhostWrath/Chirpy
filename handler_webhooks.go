package main

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/poupardm-GhostWrath/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	const defaultEvent = "user.upgraded"

	type parameters struct {
		Event	string	`json:"event"`
		Data	struct {
			UserID	uuid.UUID	`json:"user_id"`
		} `json:"data"`
	}

	// Check Auth Header
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}
	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	// Decode Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Check Event
	if params.Event != defaultEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Update User
	_, err = cfg.db.UpgradeUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}