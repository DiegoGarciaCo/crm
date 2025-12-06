package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateStage(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ClientType  string `json:"client_type"`
		OrderIndex  int    `json:"order_index"`
	}

	// Get ownerID from Context
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	stage, err := cfg.DB.CreateStage(r.Context(), database.CreateStageParams{
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		ClientType:  database.ClientType(req.ClientType),
		OrderIndex:  int32(req.OrderIndex),
		OwnerID:     uuid.NullUUID{UUID: ownerUUID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create stage", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, stage)
}

func (cfg *apiCfg) GetStages(w http.ResponseWriter, r *http.Request) {
	// Get ownerID from Context
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	stages, err := cfg.DB.GetAllStages(r.Context(), uuid.NullUUID{UUID: ownerUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve stages", err)
		return
	}

	respondWithJSON(w, http.StatusOK, stages)
}

func (cfg *apiCfg) GetStagesByClientType(w http.ResponseWriter, r *http.Request) {
	// Get ownerID from Context
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	// Get clientType from Url
	clientTypeStr := r.URL.Query().Get("client")
	if clientTypeStr == "" {
		respondWithError(w, http.StatusBadRequest, "Client type is required", nil)
		return
	}

	stages, err := cfg.DB.GetStagesByClientType(r.Context(), database.GetStagesByClientTypeParams{
		OwnerID:    uuid.NullUUID{UUID: ownerUUID, Valid: true},
		ClientType: database.ClientType(clientTypeStr),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve stages", err)
		return
	}

	respondWithJSON(w, http.StatusOK, stages)
}
