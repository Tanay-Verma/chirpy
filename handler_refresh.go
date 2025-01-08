package main

import (
	"net/http"

	"github.com/Tanay-Verma/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find refresh token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find refresh token", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(req.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
