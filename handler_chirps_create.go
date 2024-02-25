package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"slices"
)

func (c *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body string `json:"body"`
    	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err := decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := c.DB.CreateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, chirp)
}

func cleanBody(ds string) (cs string) {	
	dw := []string{"kerfuffle", "sharbert", "fornax"}
	s := strings.Split(ds, " ")
	for i, v := range s {
		if slices.Contains[[]string, string](dw, strings.ToLower(v)) {
			s[i] = "****"
		}
	}
	return strings.Join(s, " ")
}

