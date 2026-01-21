package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

func (cfg *apiCfg) GetAllSmartLists(w http.ResponseWriter, r *http.Request) {
	// Get userID from Context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	smartList, err := cfg.DB.GetAllSmartListsWithCounts(r.Context(), uuid.NullUUID{UUID: userUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get smart lists", err)
		return
	}

	respondWithJSON(w, http.StatusOK, smartList)
}

func (cfg *apiCfg) CreateSmartList(w http.ResponseWriter, r *http.Request) {
	// Get userID from Context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	type SmartListRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var smartListReq SmartListRequest
	if err := json.NewDecoder(r.Body).Decode(&smartListReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	smartList, err := cfg.DB.CreateSmartList(r.Context(), database.CreateSmartListParams{
		UserID:      uuid.NullUUID{UUID: userUUID, Valid: true},
		Name:        smartListReq.Name,
		Description: sql.NullString{String: smartListReq.Description, Valid: smartListReq.Description != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create smart list", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, smartList)
}

func (cfg *apiCfg) SetSmartListFilterCriteria(w http.ResponseWriter, r *http.Request) {
	cfg.logger.Info("SetSmartListFilterCriteria called")
	// Get smartListID from url
	smartListUUID, err := GetUUIDFromUrl("smartListID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smart list ID", err)
		return
	}

	type FilterCriteriaRequest struct {
		FilterCriteria json.RawMessage `json:"filter_criteria"`
	}

	var filterCriteriaReq FilterCriteriaRequest
	if err := json.NewDecoder(r.Body).Decode(&filterCriteriaReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	cfg.logger.Info("Received Filter Criteria: ", "Criteria:", string(filterCriteriaReq.FilterCriteria))
	cfg.logger.Info("Filter Criteria Length: ", "RawMessage", filterCriteriaReq.FilterCriteria)

	smartList, err := cfg.DB.SetSmartListFilterCriteria(r.Context(), database.SetSmartListFilterCriteriaParams{
		ID:             smartListUUID,
		FilterCriteria: pqtype.NullRawMessage{RawMessage: filterCriteriaReq.FilterCriteria, Valid: len(filterCriteriaReq.FilterCriteria) > 0},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to set filter criteria", err)
		return
	}

	cfg.logger.Info("Updated Smart List Criteria: ", "Criteria:", string(smartList.FilterCriteria.RawMessage))
	respondWithJSON(w, http.StatusOK, smartList)
}

func (cfg *apiCfg) UpdateSmartList(w http.ResponseWriter, r *http.Request) {
	// Get smartListID from url
	smartListUUID, err := GetUUIDFromUrl("smartListID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smart list ID", err)
		return
	}

	type UpdateSmartListRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var updateSmartListReq UpdateSmartListRequest
	if err := json.NewDecoder(r.Body).Decode(&updateSmartListReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	smartList, err := cfg.DB.UpdateSmartList(r.Context(), database.UpdateSmartListParams{
		ID:          smartListUUID,
		Name:        updateSmartListReq.Name,
		Description: sql.NullString{String: updateSmartListReq.Description, Valid: updateSmartListReq.Description != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update smart list", err)
		return
	}

	respondWithJSON(w, http.StatusOK, smartList)
}

func (cfg *apiCfg) ReorderSmartLists(w http.ResponseWriter, r *http.Request) {
	// Get userID from Context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	// Parse reques body
	var smartlists []database.SmartList
	err = json.NewDecoder(r.Body).Decode(&smartlists)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Start transaction
	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start database transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	// Prepare reorder params
	IDs := make([]uuid.UUID, len(smartlists))
	Indexes := make([]int32, len(smartlists))
	for i, smartlist := range smartlists {
		IDs[i] = smartlist.ID
		Indexes[i] = int32(smartlist.ListIndex)
	}

	// Bump Indexes by 1000 to avoid unique constraint conflicts
	err = qtx.BumpSmartListIndexes(r.Context(), uuid.NullUUID{UUID: userUUID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to bump smart list indexes", err)
		return
	}

	// Update order of each smart list
	err = qtx.ReorderSmartLists(r.Context(), database.ReorderSmartListsParams{
		Column1: IDs,
		Column2: Indexes,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reorder smart lists", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Smart lists reordered successfully"})
}

func (cfg *apiCfg) DeleteSmartList(w http.ResponseWriter, r *http.Request) {
	// Get smartListID from url
	smartListUUID, err := GetUUIDFromUrl("smartListID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smart list ID", err)
		return
	}

	err = cfg.DB.DeleteSmartList(r.Context(), smartListUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete smart list", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Smart list deleted successfully"})
}
