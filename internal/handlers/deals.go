package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateDeal(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID            string  `json:"contact_id"`
		Title                string  `json:"title"`
		Price                float64 `json:"price"`
		ClosingDate          string  `json:"closing_date"`
		EarnestMoneyDueDate  string  `json:"earnest_money_due_date"`
		MutualAcceptanceDate string  `json:"mutual_acceptance_date"`
		InspectionDate       string  `json:"inspection_date"`
		AppraisalDate        string  `json:"appraisal_date"`
		FinalWalkthroughDate string  `json:"final_walkthrough_date"`
		PossessionDate       string  `json:"possession_date"`
		ClosedDate           string  `json:"closed_date"`
		Commission           float64 `json:"commission"`
		CommissionSplit      float64 `json:"commission_split"`
		PropertyAddress      string  `json:"property_address"`
		PropertyCity         string  `json:"property_city"`
		PropertyState        string  `json:"property_state"`
		PropertyZip          string  `json:"property_zip"`
		Description          string  `json:"description"`
		StageID              string  `json:"stage_id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get userUUID from context
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned_to_id format", err)
		return
	}

	// Convert ContactID from string to uuid.UUID
	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact_id format", err)
		return
	}

	// Convert StageID from string to uuid.UUID
	stageUUID, err := uuid.Parse(req.StageID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid stage_id format", err)
		return
	}
	var closingDate time.Time
	var earnestMoneyDueDate time.Time
	var mutualAcceptanceDate time.Time
	var inspectionDate time.Time
	var appraisalDate time.Time
	var finalWalkthroughDate time.Time
	var possessionDate time.Time
	var closedDate time.Time

	if req.ClosedDate != "" {
		date, err := time.Parse("2006-01-02T15:04", req.ClosedDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid closed_date format", err)
			return
		}
		closedDate = date
	}

	// Convert dates from string to sql.NullTime
	if req.ClosingDate != "" {
		date, err := time.Parse("2006-01-02T15:04", req.ClosingDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid closing_date format", err)
			return
		}
		closingDate = date
	}

	if req.EarnestMoneyDueDate != "" {

		date, err := time.Parse("2006-01-02T15:04", req.EarnestMoneyDueDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid earnest_money_due_date format", err)
			return
		}
		earnestMoneyDueDate = date
	}

	if req.MutualAcceptanceDate != "" {

		date, err := time.Parse("2006-01-02T15:04", req.MutualAcceptanceDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid mutual_acceptance_date format", err)
			return
		}
		mutualAcceptanceDate = date
	}

	if req.InspectionDate != "" {

		date, err := time.Parse("2006-01-02T15:04", req.InspectionDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid inspection_date format", err)
			return
		}
		inspectionDate = date
	}

	if req.AppraisalDate != "" {
		date, err := time.Parse("2006-01-02T15:04", req.AppraisalDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid appraisal_date format", err)
			return
		}
		appraisalDate = date
	}

	if req.FinalWalkthroughDate != "" {
		date, err := time.Parse("2006-01-02T15:04", req.FinalWalkthroughDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid final_walkthrough_date format", err)
			return
		}
		finalWalkthroughDate = date
	}

	if req.PossessionDate != "" {
		date, err := time.Parse("2006-01-02T15:04", req.PossessionDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid possession_date format", err)
			return
		}
		possessionDate = date
	}

	// Turn commission and commission split into strings
	commissionStr := strconv.FormatFloat(req.Commission, 'f', 2, 64)
	commissionSplitStr := strconv.FormatFloat(req.CommissionSplit, 'f', 2, 64)

	deal, err := cfg.DB.CreateDeal(r.Context(), database.CreateDealParams{
		ContactID:            uuid.NullUUID{UUID: contactUUID, Valid: true},
		AssignedToID:         uuid.NullUUID{UUID: assignedToUUID, Valid: true},
		Title:                req.Title,
		Price:                int32(req.Price),
		ClosingDate:          sql.NullTime{Time: closingDate, Valid: closingDate != time.Time{}},
		EarnestMoneyDueDate:  sql.NullTime{Time: earnestMoneyDueDate, Valid: earnestMoneyDueDate != time.Time{}},
		MutualAcceptanceDate: sql.NullTime{Time: mutualAcceptanceDate, Valid: mutualAcceptanceDate != time.Time{}},
		InspectionDate:       sql.NullTime{Time: inspectionDate, Valid: inspectionDate != time.Time{}},
		AppraisalDate:        sql.NullTime{Time: appraisalDate, Valid: appraisalDate != time.Time{}},
		FinalWalkthroughDate: sql.NullTime{Time: finalWalkthroughDate, Valid: finalWalkthroughDate != time.Time{}},
		PossessionDate:       sql.NullTime{Time: possessionDate, Valid: possessionDate != time.Time{}},
		Commission:           sql.NullString{String: commissionStr, Valid: req.Commission != 0},
		CommissionSplit:      sql.NullString{String: commissionSplitStr, Valid: req.CommissionSplit != 0},
		PropertyAddress:      sql.NullString{String: req.PropertyAddress, Valid: req.PropertyAddress != ""},
		PropertyCity:         sql.NullString{String: req.PropertyCity, Valid: req.PropertyCity != ""},
		PropertyState:        sql.NullString{String: req.PropertyState, Valid: req.PropertyState != ""},
		PropertyZipCode:      sql.NullString{String: req.PropertyZip, Valid: req.PropertyZip != ""},
		Description:          sql.NullString{String: req.Description, Valid: req.Description != ""},
		StageID:              uuid.NullUUID{UUID: stageUUID, Valid: true},
		ClosedDate:           sql.NullTime{Time: closedDate, Valid: closedDate != time.Time{}},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create deal", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, deal)
}

func (cfg *apiCfg) GetDealByID(w http.ResponseWriter, r *http.Request) {
	dealUUID, err := GetUUIDFromUrl("dealID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID", err)
		return
	}

	deal, err := cfg.DB.GetDealById(r.Context(), dealUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get deal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, deal)
}

func (cfg *apiCfg) UpdateDeal(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContactID            string  `json:"contact_id"`
		Title                string  `json:"title"`
		Price                float64 `json:"price"`
		ClosingDate          string  `json:"closing_date"`
		EarnestMoneyDueDate  string  `json:"earnest_money_due_date"`
		MutualAcceptanceDate string  `json:"mutual_acceptance_date"`
		InspectionDate       string  `json:"inspection_date"`
		AppraisalDate        string  `json:"appraisal_date"`
		FinalWalkthroughDate string  `json:"final_walkthrough_date"`
		PossessionDate       string  `json:"possession_date"`
		ClosedDate           string  `json:"closed_date"`
		Commission           float64 `json:"commission"`
		CommissionSplit      float64 `json:"commission_split"`
		PropertyAddress      string  `json:"property_address"`
		PropertyCity         string  `json:"property_city"`
		PropertyState        string  `json:"property_state"`
		PropertyZip          string  `json:"property_zip"`
		Description          string  `json:"description"`
		StageID              string  `json:"stage_id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get assignedToUUID from context
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		cfg.logger.Info("Error parsing assigned to ID:", "Error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid assigned_to_id format", err)
		return
	}

	dealUUID, err := GetUUIDFromUrl("dealID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID", err)
		return
	}

	// Convert ContactID from string to uuid.UUID
	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		cfg.logger.Info("Error parsing contact ID:", "Error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid contact_id format", err)
		return
	}

	// Convert StageID from string to uuid.UUID
	stageUUID, err := uuid.Parse(req.StageID)
	if err != nil {
		cfg.logger.Info("Error parsing stage ID:", "Error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid stage_id format", err)
		return
	}

	var closingDate time.Time
	var earnestMoneyDueDate time.Time
	var mutualAcceptanceDate time.Time
	var inspectionDate time.Time
	var appraisalDate time.Time
	var finalWalkthroughDate time.Time
	var possessionDate time.Time
	var closedDate time.Time

	if req.ClosedDate != "" {
		date, err := time.Parse(time.RFC3339, req.ClosedDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid closed_date format", err)
			return
		}
		closedDate = date
	}

	// Convert dates from string to sql.NullTime
	if req.ClosingDate != "" {
		date, err := time.Parse(time.RFC3339, req.ClosingDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid closing_date format", err)
			return
		}
		closingDate = date
	}

	if req.EarnestMoneyDueDate != "" {

		date, err := time.Parse(time.RFC3339, req.EarnestMoneyDueDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid earnest_money_due_date format", err)
			return
		}
		earnestMoneyDueDate = date
	}

	if req.MutualAcceptanceDate != "" {

		date, err := time.Parse(time.RFC3339, req.MutualAcceptanceDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid mutual_acceptance_date format", err)
			return
		}
		mutualAcceptanceDate = date
	}

	if req.InspectionDate != "" {

		date, err := time.Parse(time.RFC3339, req.InspectionDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid inspection_date format", err)
			return
		}
		inspectionDate = date
	}

	if req.AppraisalDate != "" {
		date, err := time.Parse(time.RFC3339, req.AppraisalDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid appraisal_date format", err)
			return
		}
		appraisalDate = date
	}

	if req.FinalWalkthroughDate != "" {
		date, err := time.Parse(time.RFC3339, req.FinalWalkthroughDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid final_walkthrough_date format", err)
			return
		}
		finalWalkthroughDate = date
	}

	if req.PossessionDate != "" {
		date, err := time.Parse(time.RFC3339, req.PossessionDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid possession_date format", err)
			return
		}
		possessionDate = date
	}

	// Convert commission and commission split into strings
	commissionStr := strconv.FormatFloat(req.Commission, 'f', 2, 64)
	commissionSplitStr := strconv.FormatFloat(req.CommissionSplit, 'f', 2, 64)

	deal, err := cfg.DB.UpdateDeal(r.Context(), database.UpdateDealParams{
		ID:                   dealUUID,
		ContactID:            uuid.NullUUID{UUID: contactUUID, Valid: true},
		AssignedToID:         uuid.NullUUID{UUID: assignedToUUID, Valid: true},
		Title:                req.Title,
		Price:                int32(req.Price),
		ClosingDate:          sql.NullTime{Time: closingDate, Valid: closingDate != time.Time{}},
		EarnestMoneyDueDate:  sql.NullTime{Time: earnestMoneyDueDate, Valid: earnestMoneyDueDate != time.Time{}},
		MutualAcceptanceDate: sql.NullTime{Time: mutualAcceptanceDate, Valid: mutualAcceptanceDate != time.Time{}},
		InspectionDate:       sql.NullTime{Time: inspectionDate, Valid: inspectionDate != time.Time{}},
		AppraisalDate:        sql.NullTime{Time: appraisalDate, Valid: appraisalDate != time.Time{}},
		FinalWalkthroughDate: sql.NullTime{Time: finalWalkthroughDate, Valid: finalWalkthroughDate != time.Time{}},
		PossessionDate:       sql.NullTime{Time: possessionDate, Valid: possessionDate != time.Time{}},
		Commission:           sql.NullString{String: commissionStr, Valid: req.Commission != 0},
		CommissionSplit:      sql.NullString{String: commissionSplitStr, Valid: req.CommissionSplit != 0},
		PropertyAddress:      sql.NullString{String: req.PropertyAddress, Valid: req.PropertyAddress != ""},
		PropertyCity:         sql.NullString{String: req.PropertyCity, Valid: req.PropertyCity != ""},
		PropertyState:        sql.NullString{String: req.PropertyState, Valid: req.PropertyState != ""},
		PropertyZipCode:      sql.NullString{String: req.PropertyZip, Valid: req.PropertyZip != ""},
		Description:          sql.NullString{String: req.Description, Valid: req.Description != ""},
		StageID:              uuid.NullUUID{UUID: stageUUID, Valid: true},
		ClosedDate:           sql.NullTime{Time: closedDate, Valid: closedDate != time.Time{}},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update deal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, deal)
}

func (cfg *apiCfg) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	dealUUID, err := GetUUIDFromUrl("dealID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID", err)
		return
	}

	err = cfg.DB.DeleteDeal(r.Context(), dealUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete deal", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) ListDeals(w http.ResponseWriter, r *http.Request) {
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}
	// Get Query Parameters for filtering, pagination, etc.
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10" // Default limit
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0" // Default offset
	}

	// Convert limit and offset to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid limit parameter", err)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid offset parameter", err)
		return
	}

	deals, err := cfg.DB.ListDeals(r.Context(), database.ListDealsParams{
		AssignedToID: uuid.NullUUID{UUID: assignedToUUID, Valid: true},
		Limit:        int32(limitInt),
		Offset:       int32(offsetInt),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list deals", err)
		return
	}

	respondWithJSON(w, http.StatusOK, deals)
}

func (cfg *apiCfg) ListDealsByStageID(w http.ResponseWriter, r *http.Request) {
	// Get stageID from URL
	stageUUID, err := GetUUIDFromUrl("stageID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid stage ID", err)
		return
	}
	// Get assignedToID from Context
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}

	// Get Query Parameters for filtering, pagination, etc.
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10" // Default limit
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0" // Default offset
	}

	// Convert limit and offset to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid limit parameter", err)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid offset parameter", err)
		return
	}

	deals, err := cfg.DB.ListDealsByStage(r.Context(), database.ListDealsByStageParams{
		StageID:      uuid.NullUUID{UUID: stageUUID, Valid: true},
		AssignedToID: uuid.NullUUID{UUID: assignedToUUID, Valid: true},
		Limit:        int32(limitInt),
		Offset:       int32(offsetInt),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list deals by stage ID", err)
		return
	}

	respondWithJSON(w, http.StatusOK, deals)
}

func (cfg *apiCfg) ListDealsByContactID(w http.ResponseWriter, r *http.Request) {
	// Get contactID from URL
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	// Get assignedToID from Context
	assignedToUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assigned to ID", err)
		return
	}

	// Get Contact's Deals
	deals, err := cfg.DB.ListDealsByContactID(r.Context(), database.ListDealsByContactIDParams{
		ContactID:    uuid.NullUUID{UUID: contactUUID, Valid: true},
		AssignedToID: uuid.NullUUID{UUID: assignedToUUID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list deals by contact ID", err)
		return
	}

	respondWithJSON(w, http.StatusOK, deals)
}
