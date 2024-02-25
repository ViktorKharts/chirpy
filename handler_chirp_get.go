package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (c *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := chi.URLParam(r, "chirpID")
	chirpIDInt, _ := strconv.Atoi(chirpID)

	chirp, err := c.DB.GetChirp(chirpIDInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJson(w, http.StatusOK, chirp)
} 

