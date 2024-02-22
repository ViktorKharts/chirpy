package main

import (
	"net/http"
	"sort"

	"github.com/ViktorKharts/chirpy/internal/database"
)

func (c *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := c.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirps := []database.Chirp{}
	for _, c := range dbChirps {
		chirps = append(chirps, c)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJson(w, http.StatusOK, chirps)
}

