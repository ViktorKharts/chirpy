package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`	
}

func (c *apiConfig) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Email string `json:"email"`
	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err := decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	user, err := c.DB.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, User{
		ID: user.ID,
		Email: user.Email,
	})
}
