package main

import (
	"encoding/json"
	"strconv"
	"net/http"
	"github.com/ViktorKharts/chirpy/internal/auth"
)

type requestPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (c *apiConfig) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	t, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWToken(t)
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

	user, err := c.DB.UpdateUser(userIdInt, params.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		Email: user.Email,
	})
}
