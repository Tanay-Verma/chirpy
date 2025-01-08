package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tanay-Verma/chirpy/internal/auth"
	"github.com/Tanay-Verma/chirpy/internal/database"
)

const expiresIn = time.Hour

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get the user", err)
		return
	}

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password does not match please try again...", err)
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	const refresh_token_expiration_duration = 60 * 24 * time.Hour
	newRefreshToken, err := cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refresh_token_expiration_duration),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        jwtToken,
		RefreshToken: newRefreshToken.Token,
	})
}
