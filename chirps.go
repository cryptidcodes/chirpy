package main

import (
	"github.com/cryptidcodes/chirpy/internal/database"
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"github.com/google/uuid"
	"time"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// CONFIG
	// specify request and response structures
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
	}

	// HANDLE JSON REQUEST
	// decode the incoming JSON body into a parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// params is now a struct with data populated successfully

	// VALIDATE AND PROCESS
	// validate the body length
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// replace bad words
	strictWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedWords := []string{}
	uncleaned := strings.Split(params.Body, " ")
	for _, word := range uncleaned {
		for _, badWord := range strictWords {
			if strings.ToLower(word) == badWord {
				word = "****"
			}
		}
		cleanedWords = append(cleanedWords, word)
	}
	cleaned := strings.Join(cleanedWords, " ")

	// CREATE SQL ENTRY
	chirpParams := database.CreateChirpParams{
		Body: cleaned,
		UserID: uuid.NullUUID{
			UUID:  params.UserID,
			Valid: true,
		},
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp in database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// DEBUG
	println("Created chirp: ", chirp.Body)

	// RESPOND WITH CLEANED CHIRP
	respondWithJSON(w, http.StatusCreated, returnVals{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}