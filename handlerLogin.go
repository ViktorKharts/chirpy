package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type loginPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := loginPayload{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	user, err := c.DB.GetUser(params.Email)	
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: wrong login or password")
		return
	}
	
	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		Email: user.Email,
	})
}
