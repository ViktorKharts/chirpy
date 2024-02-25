package main

import (
	"net/http"
	"encoding/json"
	"os"
)

const (
	USER_UPGRADED = "user.upgraded"
	UPGRADED_STATUS = true
)

func (c *apiConfig) updatePolkaHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Event string `json:"event"`
		Data struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	d := json.NewDecoder(r.Body)
	params := requestPayload{} 
	err := d.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	if params.Event != USER_UPGRADED {
		respondWithJson(w, http.StatusOK, "OK")
		return
	}

	user, err := c.DB.GetUserById(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user with such id")
		return
	}

	_, err = c.DB.UpdateUser(user.ID, user.Email, user.Password, UPGRADED_STATUS)
	if err == os.ErrNotExist {
		respondWithError(w, http.StatusNotFound, "No user with such id")
		return
	} else if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)	
}
