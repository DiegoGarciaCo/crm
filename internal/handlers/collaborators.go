package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) AddCollaborator(w http.ResponseWriter, r *http.Request) {
	type req struct {
		ContactID string `json:"contact_id"`
		UserID    string `json:"user_id"`
		Role      string `json:"role"`
	}

	var request req
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Parse ID values to UUID format
	contactUUID, err := uuid.Parse(request.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Start DB transaction
	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	// Add collaborator to database
	err = qtx.AddCollaborator(r.Context(), database.AddCollaboratorParams{
		ContactID: contactUUID,
		UserID:    userUUID,
		Role:      request.Role,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add collaborator", err)
		return
	}

	// Notify the collaborator
	_, err = qtx.CreateNotification(r.Context(), database.CreateNotificationParams{
		UserID:    userUUID,
		Message:   "You have been added as a collaborator.",
		Type:      "collaborator_added",
		ContactID: uuid.NullUUID{UUID: contactUUID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create notification", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	// Get collaborator ID from URL parameters
	collaboratorUUID, err := GetUUIDFromUrl("collaboratorID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid collaborator ID", err)
		return
	}

	// Get contact ID from URL parameters
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	// Remove collaborator from database
	err = cfg.DB.RemoveCollaborator(r.Context(), database.RemoveCollaboratorParams{
		UserID:    collaboratorUUID,
		ContactID: contactUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to remove collaborator", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
