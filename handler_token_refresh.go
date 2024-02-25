package main

import (
	"net/http"
	"strconv"

	"github.com/ViktorKharts/chirpy/internal/auth"
)

type refreshTokenHandlerResponsePayload struct {
	Token string `json:"token"`
}

func (c *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	t, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	validatedToken, err := auth.ValidateJWToken(t)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token error: unauthorized")
		return
	}

	userId, err := validatedToken.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token error: unauthorized")
		return
	}

	userIdInt, err := strconv.Atoi(userId) 
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}


	issuer, err := validatedToken.GetIssuer()	
	if issuer != auth.REFRESH_JWT_ISSUER {
		respondWithError(w, http.StatusUnauthorized, "Refresh token error: unauthorized")
		return
	}

	ok, err := c.DB.IsTokenRevoked(t)
	if ok || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token error: unauthorized")
		return
	}
	
	newAccessToken, err := auth.GenerateJWToken(c.jwtSecret, auth.ACCESS_JWT_ISSUER, userIdInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, refreshTokenHandlerResponsePayload{
		Token: newAccessToken,
	})
}

