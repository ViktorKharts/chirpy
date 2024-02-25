package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"slices"
	"strings"

	"github.com/ViktorKharts/chirpy/internal/auth"
)

func (c *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body string `json:"body"`
    	}

	t, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized")
		return 
	}

	validatedToken, err := auth.ValidateJWToken(t)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized")
		return 
	}

	authorId, err := validatedToken.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized")
		return 
	}

	authorIdInt, err := strconv.Atoi(authorId) 
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return 
	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err = decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := c.DB.CreateChirp(params.Body, authorIdInt)
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

