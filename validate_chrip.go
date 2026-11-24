package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChripsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	words := strings.Split(params.Body, " ")
	var cleanBody []string

	for _, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			cleanBody = append(cleanBody, "****")
			continue
		}
		cleanBody = append(cleanBody, word)
	}
	resp := strings.Join(cleanBody, " ")

	respondWithJSON(w, http.StatusOK, validResponse{CleanBody: resp})

}
