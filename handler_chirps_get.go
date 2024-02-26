package main

import (
	"net/http"
	"strconv"
	"sort"

	"github.com/ViktorKharts/chirpy/internal/database"
)

func (c *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("author_id")

	dbChirps, err := c.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirps := []database.Chirp{}
	for _, c := range dbChirps {
		chirps = append(chirps, c)
	}

	if userID != "" {
		userIDI, err := strconv.Atoi(userID) 	
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		chirps = filterChirpsByUserID(chirps, userIDI)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJson(w, http.StatusOK, chirps)
}

func filterChirpsByUserID(c []database.Chirp, userID int) (r []database.Chirp) {
	for _, chirp := range c {
		if chirp.AuthorID == userID {
			r = append(r, chirp)
		}
	}
	return r
}

