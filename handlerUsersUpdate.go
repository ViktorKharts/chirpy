package main

import (
	"encoding/json"
	"strconv"
	"net/http"
	"strings"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	TOKEN_PREFIX = "Bearer "
)

type requestPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (c *apiConfig) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	t := r.Header.Get("Authorization")
	t, ok := strings.CutPrefix(t, TOKEN_PREFIX) 
	if !ok {
		respondWithError(w, http.StatusNotFound, "No bearer token found")
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(t, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(JWT_SECRET)), nil
	})	
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JSON Web Token")
		return
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Authorization error. Please, login.")
		return
	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err = decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 8)
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
