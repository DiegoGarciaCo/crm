package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiCfg) Get5NewestContacts(w http.ResponseWriter, r *http.Request) {
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}
	contacts, err := cfg.DB.Get5NewestContacts(r.Context(), uuid.NullUUID{UUID: ownerUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve contacts", err)
		return
	}

	respondWithJSON(w, http.StatusOK, contacts)
}

func (cfg *apiCfg) GetNewContactsCount(w http.ResponseWriter, r *http.Request) {
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	count, err := cfg.DB.NewContactsThisMonth(r.Context(), uuid.NullUUID{UUID: ownerUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve new contacts count", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"new_contacts_count": int(count)})
}

func (cfg *apiCfg) GetAppointmentsCount(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignedTo ID", err)
		return
	}
	count, err := cfg.DB.AppointmentsThisWeek(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve appointments count", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"appointments_count": int(count)})
}

func (cfg *apiCfg) GetTasksDueTodayCount(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignedTo ID", err)
		return
	}

	count, err := cfg.DB.TasksDueTodayCount(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tasks due today count", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"tasks_due_today_count": int(count)})
}

func (cfg *apiCfg) Get5UpcomingAppointments(w http.ResponseWriter, r *http.Request) {
	// Get assignedToID from URL
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignedTo ID", err)
		return
	}

	appointments, err := cfg.DB.GetUpcomingAppointments(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve upcoming appointments", err)
		return
	}

	respondWithJSON(w, http.StatusOK, appointments)
}

func (cfg *apiCfg) GetContactsCount(w http.ResponseWriter, r *http.Request) {
	// Get ownerID from URL
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	count, err := cfg.DB.ContactsCount(r.Context(), uuid.NullUUID{UUID: ownerUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve contacts count", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"contacts_count": int(count)})
}
