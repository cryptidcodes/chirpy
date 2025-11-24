package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cryptidcodes/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	key, err := auth.GetAPIKey(r.Header)
	if err != nil || key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid API key", err)
		return
	}

	// parse request body
	var params parameters
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Printf("decode error: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// only handle user.upgraded events
	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "event not handled", nil)
	} else {
		// upgrade user to chirpy red
		_, err := cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
		if err != nil {
			http.Error(w, "Failed to upgrade user", http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}
