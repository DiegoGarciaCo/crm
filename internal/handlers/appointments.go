package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID   string `json:"contact_id"`
		Title       string `json:"title"`
		ScheduledAt string `json:"scheduled_at"`
		Notes       string `json:"notes"`
		Outcome     string `json:"outcome"`
		Location    string `json:"location"`
		Type        string `json:"type"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get assignedToUUID from context
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No User ID in Context", err)
		return
	}

	// Convert ContactID to UUID
	contactID, err := uuid.Parse(req.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	// Parse ScheduledAt to time.Time
	scheduledAt, err := time.Parse("2006-01-02T15:04", req.ScheduledAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid scheduled at format", err)
		return
	}

	appointment, err := cfg.DB.CreateAppointment(r.Context(), database.CreateAppointmentParams{
		AssignedToID: uuid.NullUUID{UUID: assignedToUUID, Valid: true},
		ContactID:    uuid.NullUUID{UUID: contactID, Valid: true},
		Title:        req.Title,
		ScheduledAt:  scheduledAt,
		Location:     sql.NullString{String: req.Location, Valid: req.Location != ""},
		Type:         database.NullAppointmentType{AppointmentType: database.AppointmentType(req.Type), Valid: req.Type != ""},
		Outcome:      database.NullAppointmentOutcome{AppointmentOutcome: database.AppointmentOutcome(req.Outcome), Valid: req.Outcome != ""},
		Note:         sql.NullString{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create appointment", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, appointment)
}

func (cfg *apiCfg) GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
	appointmentUUID, err := GetUUIDFromUrl("appointmentID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid appointment ID", err)
		return
	}

	appointment, err := cfg.DB.GetAppointmentById(r.Context(), appointmentUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get appointment", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointment)
}

func (cfg *apiCfg) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID   string `json:"contact_id"`
		Title       string `json:"title"`
		ScheduledAt string `json:"scheduled_at"`
		Notes       string `json:"notes"`
		Outcome     string `json:"outcome"`
		Location    string `json:"location"`
		Type        string `json:"type"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	appointmentUUID, err := GetUUIDFromUrl("appointmentID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid appointment ID", err)
		return
	}

	// Parse ScheduledAt to time.Time
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid scheduled at format", err)
		return
	}

	appointment, err := cfg.DB.UpdateAppointment(r.Context(), database.UpdateAppointmentParams{
		ID:          appointmentUUID,
		ContactID:   uuid.NullUUID{UUID: uuid.MustParse(req.ContactID), Valid: req.ContactID != ""},
		Title:       req.Title,
		ScheduledAt: scheduledAt,
		Location:    sql.NullString{String: req.Location, Valid: req.Location != ""},
		Type:        database.NullAppointmentType{AppointmentType: database.AppointmentType(req.Type), Valid: req.Type != ""},
		Outcome:     database.NullAppointmentOutcome{AppointmentOutcome: database.AppointmentOutcome(req.Outcome), Valid: req.Outcome != ""},
		Note:        sql.NullString{String: req.Notes, Valid: req.Notes != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update appointment", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointment)
}

func (cfg *apiCfg) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentUUID, err := GetUUIDFromUrl("appointmentID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid appointment ID", err)
		return
	}

	err = cfg.DB.DeleteAppointment(r.Context(), appointmentUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete appointment", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) ListAppointmentsToday(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}
	appointments, err := cfg.DB.ListTodaysAppointments(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list today's appointments", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointments)
}

func (cfg *apiCfg) ListAppointments(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}
	appointments, err := cfg.DB.ListAppointments(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list appointments", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointments)
}

func (cfg *apiCfg) ListUpcomingAppointments(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}
	appointments, err := cfg.DB.ListUpcomingAppointments(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list upcoming appointments", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointments)
}

func (cfg *apiCfg) ListAppointmentsByContactID(w http.ResponseWriter, r *http.Request) {
	contactUUID, err := GetUUIDFromUrl("ContactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	appointments, err := cfg.DB.ListAppointmentsByContactId(r.Context(), uuid.NullUUID{UUID: contactUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list appointments by contact ID", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointments)
}
