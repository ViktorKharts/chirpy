package main

import (
	"net/http"
	"strconv"
	"os"

	"github.com/ViktorKharts/chirpy/internal/auth"
	"github.com/go-chi/chi"
)

func (c *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	t, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedToken, err := auth.ValidateJWToken(t)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := validatedToken.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get data from token")
		return
	}

	userIDInt, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	chirpID := chi.URLParam(r, "chirpID")
	chirpIDInt, err := strconv.Atoi(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = c.DB.DeleteChirp(chirpIDInt, userIDInt)
	if err == os.ErrPermission {
		respondWithError(w, http.StatusForbidden, "You can only delete your Chirps")
		return
	} else if err == os.ErrNotExist {
		respondWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	} else if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "Chirp was deleted")
}

