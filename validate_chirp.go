package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_Body: cleaned,
	})
}