package main

import (
	"net/http"

	"github.com/ViktorKharts/chirpy/internal/auth"
)

func (c *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	t, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = c.DB.RevokeToken(t)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "Token was revoked")
}
