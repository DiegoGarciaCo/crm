package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/auth"
	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateUser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Decode the request body into a User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	// Save user to Database
	createdUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	// Respond with the created user
	respondWithJSON(w, http.StatusCreated, Response{
		Name:  createdUser.Username,
		Email: createdUser.Email,
	})
}

func (cfg *apiCfg) CreateUserWithTeam(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		TeamID   string `json:"team_id"`
	}

	type Response struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		TeamID string `json:"team_id,omitempty"`
	}

	// Decode the request body into a User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Prepare TeamID
	TeamUUID, err := uuid.Parse(user.TeamID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team ID", err)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	// Save user to Database
	createdUser, err := cfg.DB.CreateUserWithTeam(r.Context(), database.CreateUserWithTeamParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		TeamID:       uuid.NullUUID{UUID: TeamUUID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	// Respond with the created user
	respondWithJSON(w, http.StatusCreated, Response{
		Name:   createdUser.Username,
		Email:  createdUser.Email,
		TeamID: createdUser.TeamID.UUID.String(),
	})
}
