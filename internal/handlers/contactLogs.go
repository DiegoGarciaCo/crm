package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) LogContact(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID     string `json:"contact_id"`
		ContactMethod string `json:"contact_method"`
		Note          string `json:"note"`
	}
	createdBy, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	log, err := cfg.DB.LogContact(r.Context(), database.LogContactParams{
		ContactID:     uuid.NullUUID{UUID: contactUUID, Valid: true},
		ContactMethod: req.ContactMethod,
		CreatedBy:     uuid.NullUUID{UUID: createdBy, Valid: createdBy != uuid.Nil},
		Note:          sql.NullString{String: req.Note, Valid: req.Note != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create contact log", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, log)
}

func (cfg *apiCfg) GetContactLogsByContactID(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID            uuid.UUID
		ContactID     uuid.NullUUID
		ContactMethod string
		CreatedBy     uuid.NullUUID
		Note          sql.NullString
		CreatedAt     sql.NullTime
		UpdatedAt     sql.NullTime
	}

	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	logs, err := cfg.DB.GetContactLogsByContactID(r.Context(), uuid.NullUUID{UUID: contactUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get contact logs", err)
		return
	}

	var resp []response
	for _, log := range logs {
		resp = append(resp, response{
			ID:            log.ID,
			ContactID:     log.ContactID,
			ContactMethod: log.ContactMethod,
			CreatedBy:     log.CreatedBy,
			Note:          log.Note,
			CreatedAt:     log.CreatedAt,
			UpdatedAt:     log.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}
