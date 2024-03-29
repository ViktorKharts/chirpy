package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`	
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (c *apiConfig) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err := decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 8)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} 

	user, err := c.DB.CreateUser(params.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, User{
		ID: user.ID,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
