package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tanay-Verma/chirpy/internal/auth"
	"github.com/Tanay-Verma/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	ID          uuid.UUID `json:"id"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	newUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add new user to DB", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	user_id, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't generate hashed password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             user_id,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't update the user credentials", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          updatedUser.ID,
			CreatedAt:   updatedUser.CreatedAt,
			UpdatedAt:   updatedUser.UpdatedAt,
			Email:       updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
		},
	})
}
