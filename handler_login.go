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
		AccessToken string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := loginPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	user, err := c.DB.GetUserByEmail(params.Email)	
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if err = auth.CompareHashToPassword(user.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: wrong login or password")
		return
	}

	signedAccessToken, err := auth.GenerateJWToken(c.jwtSecret, auth.ACCESS_JWT_ISSUER, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	signedRefreshToken, err := auth.GenerateJWToken(c.jwtSecret, auth.REFRESH_JWT_ISSUER, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJson(w, http.StatusOK, loginResponse{
		User: User{
			ID: user.ID,
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		AccessToken: signedAccessToken,
		RefreshToken: signedRefreshToken,
	})
}

