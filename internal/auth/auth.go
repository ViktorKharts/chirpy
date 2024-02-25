package auth

import (
	"net/http"
	"strconv"
	"strings"
	"errors"
	"time"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	ACCESS_JWT_ISSUER = "chirpy-access"
	ACCESS_JWT_DURATION = 60 * 60
	REFRESH_JWT_ISSUER = "chirpy-refresh"
	REFRESH_JWT_DURATION = 60 * 24 * 60 * 60 
	TOKEN_PREFIX = "Bearer "
	POLKA_API_PREFIX = "ApiKey "
	JWT_SECRET = "JWT_SECRET"
)


func GenerateJWToken(secret, issuer string, userId int) (string, error) {
	if issuer == ACCESS_JWT_ISSUER {
		return signToken(userId, ACCESS_JWT_ISSUER, secret, ACCESS_JWT_DURATION)
	}
	return signToken(userId, REFRESH_JWT_ISSUER, secret, REFRESH_JWT_DURATION)
}

func signToken(userId int, issuer, secret string, expiresIn int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: issuer,
		IssuedAt: jwt.NewNumericDate(time.Now().Local()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiresIn))),
		Subject: strconv.Itoa(userId),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func GetBearerToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("No Authorization header included")
	}
	t, _ := strings.CutPrefix(h, TOKEN_PREFIX) 
	return t, nil
}

func GetPolkaApiKey(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("No Authorization header included")
	}
	k, _ := strings.CutPrefix(h, POLKA_API_PREFIX) 
	return k, nil
}

func ValidateJWToken(t string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(t, jwt.MapClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(JWT_SECRET)), nil
	})	
	if err != nil {
		return jwt.MapClaims{}, errors.New("Invalid JSON Web Token")
	}

	return token.Claims, nil
}

func HashPassword(pw string) (string, error)  {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		return "", err
	} 
	return string(hashedPassword), nil
} 

func CompareHashToPassword(h, pw string ) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(pw))
}

