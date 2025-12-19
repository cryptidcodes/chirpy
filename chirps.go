package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/cryptidcodes/chirpy/internal/auth"
	"github.com/cryptidcodes/chirpy/internal/database"
	"github.com/google/uuid"
)

// DO NOT DELETE: USED IN RESPONSE STRUCTURES
// database.Chirps DOES NOT HAVE JSON TAGS
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// specify request and response structures
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// validate JWT from headers
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", nil)
		return
	}
	UserID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// decode the incoming JSON body into a parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		Body:   cleaned,
		UserID: UserID,
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp in database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// RESPOND WITH CLEANED CHIRP
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// returns all chirps in the db

	// check query params
	q := r.URL.Query()
	author_ID := q.Get("author_id")
	sortOrder := q.Get("sort")
	if author_ID != "" {
		userID, err := uuid.Parse(author_ID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "couldn't parse userID", err)
			return
		}

		chirps, err := cfg.dbQueries.GetAllChirpsByUser(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "couldnt get user's chirps", err)
			return
		}

		resp := make([]Chirp, len(chirps))
		// parse chirps into response format
		for i := range chirps {
			resp[i] = Chirp{
				ID:        chirps[i].ID,
				CreatedAt: chirps[i].CreatedAt,
				UpdatedAt: chirps[i].UpdatedAt,
				Body:      chirps[i].Body,
				UserID:    chirps[i].UserID,
			}
		}

		// sort if needed
		if sortOrder == "desc" {
			// simple bubble sort for descending order by CreatedAt
			sort.Slice(resp, func(i, j int) bool { return resp[i].CreatedAt.After(resp[j].CreatedAt) })
		}

		respondWithJSON(w, http.StatusOK, resp)
		return
	}

	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	resp := make([]Chirp, len(chirps))
	// parse chirps into response format
	for i := range chirps {
		resp[i] = Chirp{
			ID:        chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body:      chirps[i].Body,
			UserID:    chirps[i].UserID,
		}
	}

	// sort if needed
	if sortOrder == "desc" {
		// simple bubble sort for descending order by CreatedAt
		sort.Slice(resp, func(i, j int) bool { return resp[i].CreatedAt.After(resp[j].CreatedAt) })
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
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// authenticate user via JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", nil)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// extract chirpID from URL
	chirpIDstring := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstring)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// get chirp to verify it exists
	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	// verify that the authenticated user is the owner of the chirp
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to delete this chirp", nil)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	// respond with no content
	respondWithJSON(w, http.StatusNoContent, nil)
}
