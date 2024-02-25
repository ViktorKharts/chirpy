package main

import (
	"encoding/json"
	"net/http"

	"github.com/ViktorKharts/chirpy/internal/auth"
)

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type loginPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresIn *int `json:"expires_in_seconds,omitempty"`
	}

	type loginResponse struct {
		User
		Token string `json:"token"`
	}

	params := loginPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	user, err := c.DB.GetUser(params.Email)	
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if err = auth.CompareHashToPassword(user.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: wrong login or password")
		return
	}

	signedToken, err := auth.GenerateJWToken(c.jwtSecret, params.ExpiresIn, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJson(w, http.StatusOK, loginResponse{
		User: User{
			ID: user.ID,
			Email: user.Email,
		},
		Token: signedToken,
	})
}

