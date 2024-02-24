package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JWT_MAX_HOURS_STRING = "86400s"
	JWT_MAX_HOURS_INT = 86400
	JWT_ISSUER = "chirpy"
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: wrong login or password")
		return
	}

	signedToken, err := generateJWToken(c, params.ExpiresIn, user.ID)
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

func generateJWToken(c *apiConfig, jwtDuration *int, userId int) (string, error) {
	expiresIn, err := time.ParseDuration(JWT_MAX_HOURS_STRING)
	if err != nil {
		return "", err
	}

	if jwtDuration != nil {
		if *jwtDuration < JWT_MAX_HOURS_INT {
			s := fmt.Sprintf("%ds", *jwtDuration)
			if expiresIn, err = time.ParseDuration(s); err != nil {
				return "", err
			}
		}
	}

	claims := jwt.RegisteredClaims{
		Issuer: JWT_ISSUER,
		IssuedAt: jwt.NewNumericDate(time.Now().Local()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: strconv.Itoa(userId),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	st, err := t.SignedString([]byte(c.jwtSecret))
	if err != nil {
		return "", err
	}

	return st, nil
}

