package main

import (
	"github.com/google/uuid"
	"time"
	"net/http"
	"encoding/json"
)

type User struct {
		ID uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email string `json:"email"`
	}	

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// define request and response structures for this endpoint
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		User
	}
	
	// decode JSON request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// create new user in database
	newUser, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	
	// create and send JSON response
	respondWithJSON(w, 201, response{
		User: User{
			ID: newUser.ID,
			Created_at: newUser.CreatedAt,
			Updated_at: newUser.UpdatedAt,
			Email: newUser.Email,
		},
	})
}
