package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) GetNotifications(w http.ResponseWriter, r *http.Request) {
	// Get User ID from request context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Get Limit and Offset from query parameters
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Validate and convert limit and offset
	if limit == "" {
		limit = "10" // default limit
	}
	if offset == "" {
		offset = "0" // default offset
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid offset parameter", err)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid limit parameter", err)
		return
	}

	// Fetch notifications for the user
	notifications, err := cfg.DB.GetNotificationsByUserID(r.Context(), database.GetNotificationsByUserIDParams{
		UserID: userUUID,
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch notifications", err)
		return
	}

	// Respond with the notifications
	respondWithJSON(w, http.StatusOK, notifications)
}

func (cfg *apiCfg) CreateNotification(w http.ResponseWriter, r *http.Request) {
	type req struct {
		UserID        string `json:"user_id"`
		Message       string `json:"message"`
		Type          string `json:"type"`
		ContactID     string `json:"contact_id"`
		AppointmentID string `json:"appointment_id"`
		TaskID        string `json:"task_id"`
	}

	var request req
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Validate required fields
	if request.UserID == "" || request.Message == "" || request.Type == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	// Parse User ID
	userUUID, err := uuid.Parse(request.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var contactUUID, appointmentUUID, taskUUID uuid.UUID
	// Parse IDs if necessary
	if request.ContactID != "" {
		contactUUID, err = uuid.Parse(request.ContactID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
			return
		}
	}
	if request.AppointmentID != "" {
		appointmentUUID, err = uuid.Parse(request.AppointmentID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid appointment ID", err)
			return
		}
	}
	if request.TaskID != "" {
		taskUUID, err = uuid.Parse(request.TaskID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid task ID", err)
			return
		}
	}

	// Create a new notification
	notification, err := cfg.DB.CreateNotification(r.Context(), database.CreateNotificationParams{
		UserID:        userUUID,
		Message:       request.Message,
		Type:          request.Type,
		ContactID:     uuid.NullUUID{UUID: contactUUID, Valid: request.ContactID != ""},
		AppointmentID: uuid.NullUUID{UUID: appointmentUUID, Valid: request.AppointmentID != ""},
		TaskID:        uuid.NullUUID{UUID: taskUUID, Valid: request.TaskID != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create notification", err)
		return
	}

	// Respond with the created notification
	respondWithJSON(w, http.StatusCreated, notification)
}

func (cfg *apiCfg) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	// Get Notification ID from URL parameters
	notificationUUID, err := GetUUIDFromUrl("notificationID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	// Mark the notification as read
	err = cfg.DB.MarkNotificationAsRead(r.Context(), notificationUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to mark notification as read", err)
		return
	}

	// Respond with no content
	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	// Get Notification ID from URL parameters
	notificationUUID, err := GetUUIDFromUrl("notificationID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	// Delete the notification
	err = cfg.DB.DeleteNotification(r.Context(), notificationUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete notification", err)
		return
	}

	// Respond with no content
	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	// Get User ID from request context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Mark all notifications as read for the user
	err = cfg.DB.MarkAllNotificationsAsRead(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to mark all notifications as read", err)
		return
	}

	// Respond with no content
	respondWithJSON(w, http.StatusNoContent, nil)
}
