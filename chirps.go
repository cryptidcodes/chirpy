package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"github.com/google/uuid"
	"time"
	"github.com/cryptidcodes/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// specify request and response structures
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

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
		UserID: params.UserID,
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp in database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// RESPOND WITH CLEANED CHIRP
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// returns all chirps in the db

	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	resp := make([]Chirp, len(chirps))
	// parse chirps into response format
	for i := range chirps {
		resp[i] = Chirp{
			ID: chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body: chirps[i].Body,
			UserID: chirps[i].UserID,
		}
	}


	// respond with JSON
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	// extract chirpID from URL
	chirpID := r.PathValue("chirpID")
	ID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	
	c, err := cfg.dbQueries.GetChirpByID(r.Context(), ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	// respond with JSON
	respondWithJSON(w, http.StatusOK, Chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body: c.Body,
		UserID: c.UserID,
	})
}