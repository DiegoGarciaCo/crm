package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateTag(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	type request struct {
		TagName        string `json:"tag_name"`
		TagDescription string `json:"tag_description"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	tag, err := cfg.DB.CreateTag(r.Context(), database.CreateTagParams{
		UserID:      uuid.NullUUID{UUID: userUUID, Valid: true},
		Name:        req.TagName,
		Description: sql.NullString{String: req.TagDescription, Valid: req.TagDescription != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create tag", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, tag)
}

func (cfg *apiCfg) GetAllTags(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	tags, err := cfg.DB.GetAllTags(r.Context(), uuid.NullUUID{UUID: userUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tags", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tags)
}

func (cfg *apiCfg) DeleteTag(w http.ResponseWriter, r *http.Request) {
	// Get TagID from URL
	tagUUID, err := GetUUIDFromUrl("tagID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	err = cfg.DB.DeleteTag(r.Context(), tagUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete tag", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Tag deleted successfully"})
}

func (cfg *apiCfg) AssignTagToContact(w http.ResponseWriter, r *http.Request) {
	// Get TagID from URL
	tagUUID, err := GetUUIDFromUrl("tagID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	// Get ContactID from URL
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	tag, err := cfg.DB.AssignTagToContact(r.Context(), database.AssignTagToContactParams{
		TagID:     tagUUID,
		ContactID: contactUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to assign tag to contact", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tag)
}

func (cfg *apiCfg) RemoveTagFromContact(w http.ResponseWriter, r *http.Request) {
	// Get TagID from URL
	tagUUID, err := GetUUIDFromUrl("tagID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	// Get ContactID from URL
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	err = cfg.DB.RemoveTagFromContact(r.Context(), database.RemoveTagFromContactParams{
		TagID:     tagUUID,
		ContactID: contactUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to remove tag from contact", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
