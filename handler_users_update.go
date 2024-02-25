package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ViktorKharts/chirpy/internal/auth"
)

func (c *apiConfig) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

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

	issuer, err := validatedToken.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get data from token")
		return
	}

	userId, err := validatedToken.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get data from token")
		return
	}

	if issuer == auth.REFRESH_JWT_ISSUER {
		respondWithError(w, http.StatusUnauthorized, "Wrong JSON Web Token provided")
		return
	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err = decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} 

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := c.DB.UpdateUser(userIdInt, params.Email, string(hashedPassword), params.IsChirpyRed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
