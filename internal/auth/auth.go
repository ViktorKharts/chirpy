package auth

import (
	"fmt"
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
	JWT_MAX_HOURS_INT = 24 * 60 * 60 
	JWT_MAX_HOURS_STRING = "86400s"
	JWT_ISSUER = "chirpy"
	TOKEN_PREFIX = "Bearer "
	JWT_SECRET = "JWT_SECRET"
)


func GenerateJWToken(secret string, jwtDuration *int, userId int) (string, error) {
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
	st, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return st, nil
}

func GetBearerToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("No Authorization header included")
	}
	t, _ := strings.CutPrefix(h, TOKEN_PREFIX) 
	return t, nil
}

func ValidateJWToken(t string) (string, error) {
	token, err := jwt.ParseWithClaims(t, jwt.MapClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(JWT_SECRET)), nil
	})	
	if err != nil {
		return "", errors.New("Invalid JSON Web Token")
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("Authorization error. Please, login.") 
	}

	return userId, nil
}

func HashPassword(pw string) (string, error)  {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		return "", err
	} 
	return string(hashedPassword), nil
} 

func CompareHashToPassword(h, pw string ) error {
	if err := bcrypt.CompareHashAndPassword([]byte(h), []byte(pw)); err != nil {
		return err
	}
	return nil
}

