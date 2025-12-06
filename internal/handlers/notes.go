package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateNote(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID string `json:"contact_id"`
		Note      string `json:"note"`
		CreatedBy string `json:"created_by"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Parse ContactID and CreatedBy to UUIDs
	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}
	createdByUUID, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid created by ID", err)
		return
	}

	note, err := cfg.DB.CreateNote(r.Context(), database.CreateNoteParams{
		ContactID: uuid.NullUUID{UUID: contactUUID, Valid: req.ContactID != ""},
		Note:      req.Note,
		CreatedBy: uuid.NullUUID{UUID: createdByUUID, Valid: req.CreatedBy != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create note", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, note)
}

func (cfg *apiCfg) GetNotesByContactID(w http.ResponseWriter, r *http.Request) {
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	notes, err := cfg.DB.GetNotesByContactID(r.Context(), uuid.NullUUID{UUID: contactUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve notes", err)
		return
	}

	respondWithJSON(w, http.StatusOK, notes)
}
