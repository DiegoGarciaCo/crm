package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateTask(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID string `json:"contact_id"`
		Title     string `json:"title"`
		Type      string `json:"type"`
		Date      string `json:"date"`
		Status    string `json:"status"`
		Priority  string `json:"priority"`
		Note      string `json:"note"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	AssignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No User ID in Context", err)
		return
	}

	// Convert string IDs to UUIDs and handle errors
	ContactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact_id", err)
		return
	}

	// Parse date string to time.Time
	var parsedDate time.Time
	if req.Date != "" {
		parsedDate, err = time.Parse("2006-01-02T15:04", req.Date)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid date format", err)
			return
		}
	}

	// Call the CreateTask method from the queries
	task, err := cfg.DB.CreateTask(r.Context(), database.CreateTaskParams{
		ContactID:    uuid.NullUUID{UUID: ContactUUID, Valid: true},
		AssignedToID: uuid.NullUUID{UUID: AssignedToUUID, Valid: true},
		Title:        req.Title,
		Type:         database.NullTaskType{TaskType: database.TaskType(req.Type), Valid: req.Type != ""},
		Date:         sql.NullTime{Time: parsedDate, Valid: req.Date != ""},
		Status:       database.NullTaskStatus{TaskStatus: database.TaskStatus(req.Status), Valid: req.Status != ""},
		Priority:     database.NullTaskPriority{TaskPriority: database.TaskPriority(req.Priority), Valid: req.Priority != ""},
		Note:         sql.NullString{String: req.Note, Valid: req.Note != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create task", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (cfg *apiCfg) GetTasksByContactID(w http.ResponseWriter, r *http.Request) {
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	tasks, err := cfg.DB.GetTasksByContactID(r.Context(), uuid.NullUUID{UUID: contactUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tasks", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (cfg *apiCfg) GetTaskByAssignedToID(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}

	tasks, err := cfg.DB.GetTaskByAssignedToID(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tasks", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (cfg *apiCfg) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskUUID, err := GetUUIDFromUrl("taskID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	task, err := cfg.DB.GetTaskByID(r.Context(), taskUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get task", err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (cfg *apiCfg) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskUUID, err := GetUUIDFromUrl("taskID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	err = cfg.DB.DeleteTask(r.Context(), taskUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete task", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) GetOverdueTasks(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}

	tasks, err := cfg.DB.GetOverdueTasks(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get overdue tasks", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (cfg *apiCfg) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
	}

	taskUUID, err := GetUUIDFromUrl("taskID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedTask, err := cfg.DB.UpdateTaskStatus(r.Context(), database.UpdateTaskStatusParams{
		ID:     taskUUID,
		Status: database.NullTaskStatus{TaskStatus: database.TaskStatus(req.Status), Valid: req.Status != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update task status", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedTask)
}

func (cfg *apiCfg) UpdateTask(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID    string `json:"contact_id"`
		AssignedToID string `json:"assigned_to_id"`
		Title        string `json:"title"`
		Type         string `json:"type"`
		Date         string `json:"date"`
		Status       string `json:"status"`
		Priority     string `json:"priority"`
		Note         string `json:"note"`
	}

	taskUUID, err := GetUUIDFromUrl("taskID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID", err)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get AssignedToUUID from context
	AssignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No user ID in context", err)
		return
	}

	ContactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact_id", err)
		return
	}

	var parsedDate time.Time
	if req.Date != "" {
		parsedDate, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid date format", err)
			return
		}
	}

	updatedTask, err := cfg.DB.UpdateTask(r.Context(), database.UpdateTaskParams{
		ID:           taskUUID,
		ContactID:    uuid.NullUUID{UUID: ContactUUID, Valid: true},
		AssignedToID: uuid.NullUUID{UUID: AssignedToUUID, Valid: true},
		Title:        req.Title,
		Type:         database.NullTaskType{TaskType: database.TaskType(req.Type), Valid: req.Type != ""},
		Date:         sql.NullTime{Time: parsedDate, Valid: req.Date != ""},
		Status:       database.NullTaskStatus{TaskStatus: database.TaskStatus(req.Status), Valid: req.Status != ""},
		Priority:     database.NullTaskPriority{TaskPriority: database.TaskPriority(req.Priority), Valid: req.Priority != ""},
		Note:         sql.NullString{String: req.Note, Valid: req.Note != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update task", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedTask)
}

func (cfg *apiCfg) GetTasksDueToday(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}

	tasks, err := cfg.DB.GetTaskDueToday(r.Context(), uuid.NullUUID{UUID: assignedToUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tasks due today", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}
