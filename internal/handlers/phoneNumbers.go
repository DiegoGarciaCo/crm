package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PhoneNumber string `json:"phone_number"`
		Type        string `json:"type"`
		IsPrimary   bool   `json:"is_primary"`
	}

	// Get Contact ID from URL parameters
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	// Decode the request body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" {
		respondWithError(w, http.StatusBadRequest, "Phone number is required", err)
		return
	}

	err = cfg.DB.EnterPhoneNumber(r.Context(), database.EnterPhoneNumberParams{
		ContactID:   uuid.NullUUID{UUID: contactUUID, Valid: true},
		PhoneNumber: req.PhoneNumber,
		Type:        sql.NullString{String: req.Type, Valid: req.Type != ""},
		IsPrimary:   sql.NullBool{Bool: req.IsPrimary, Valid: req.IsPrimary},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create phone number", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) DeletePhoneNumber(w http.ResponseWriter, r *http.Request) {
	// Get Phone Number ID from URL parameters
	phoneNumberUUID, err := GetUUIDFromUrl("phoneNumberID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid phone number ID", err)
		return
	}

	err = cfg.DB.DeletePhoneNumber(r.Context(), phoneNumberUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete phone number", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) UpdatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PhoneNumber string `json:"phone_number"`
		Type        string `json:"type"`
		IsPrimary   bool   `json:"is_primary"`
	}

	// Get Phone Number ID from URL parameters
	phoneNumberUUID, err := GetUUIDFromUrl("phoneNumberID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid phone number ID", err)
		return
	}

	// Decode the request body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" {
		respondWithError(w, http.StatusBadRequest, "Phone number is required", err)
		return
	}

	err = cfg.DB.UpdatePhoneNumber(r.Context(), database.UpdatePhoneNumberParams{
		ID:          phoneNumberUUID,
		PhoneNumber: req.PhoneNumber,
		Type:        sql.NullString{String: req.Type, Valid: req.Type != ""},
		IsPrimary:   sql.NullBool{Bool: req.IsPrimary, Valid: req.IsPrimary},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update phone number", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
