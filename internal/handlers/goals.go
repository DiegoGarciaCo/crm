package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) SetGoal(w http.ResponseWriter, r *http.Request) {
	// Get User ID from url
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		slog.Error("Invalid User ID", "error", err, "UserID", userUUID)
		respondWithError(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}
	type request struct {
		Year                              int    `json:"year"`
		Month                             int    `json:"month"`
		Income_goal                       string `json:"income_goal"`
		Transaction_goal                  string `json:"transaction_goal"`
		Estimated_average_sale_price      string `json:"estimated_average_sale_price"`
		Estimated_average_commission_rate string `json:"estimated_average_commission_rate"`
	}

	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Parse transaction goal to int
	transactionGoal, err := strconv.Atoi(req.Transaction_goal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Transaction Goal", err)
		return
	}

	// Parse and validate request body

	goal, err := cfg.DB.SetGoal(r.Context(), database.SetGoalParams{
		UserID:                         uuid.NullUUID{UUID: userUUID, Valid: true},
		Year:                           int32(req.Year),
		Month:                          int32(req.Month),
		IncomeGoal:                     sql.NullString{String: req.Income_goal, Valid: req.Income_goal != ""},
		TransactionGoal:                sql.NullInt32{Int32: int32(transactionGoal), Valid: transactionGoal != 0},
		EstimatedAverageSalePrice:      sql.NullString{String: req.Estimated_average_sale_price, Valid: req.Estimated_average_sale_price != ""},
		EstimatedAverageCommissionRate: sql.NullString{String: req.Estimated_average_commission_rate, Valid: req.Estimated_average_commission_rate != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to set goal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, goal)
}

func (cfg *apiCfg) GetGoalByUserAndYear(w http.ResponseWriter, r *http.Request) {
	// Get User ID from Context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}

	// Get Year from url params
	year := r.URL.Query().Get("year")
	if year == "" {
		respondWithError(w, http.StatusBadRequest, "Year is required", nil)
		return
	}

	// Parse year to int
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Year", err)
		return
	}

	goal, err := cfg.DB.GetGoalByUserAndYear(r.Context(), database.GetGoalByUserAndYearParams{
		UserID: uuid.NullUUID{UUID: userUUID, Valid: true},
		Year:   int32(yearInt),
	})
	if err == sql.ErrNoRows {
		cfg.logger.Info("No goal found for user", "ID", userUUID, "Year:", yearInt)
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"goal": nil,
		})
		return
	}
	if err != nil {
		fmt.Println("Error fetching goal:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get goal", err)
		return
	}

	cfg.logger.Info("Goal fetched successfully", "UserID", userUUID, "Year:", yearInt, "goal", goal)
	respondWithJSON(w, http.StatusOK, goal)
}

func (cfg *apiCfg) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	// Get Goal ID from url
	goalUUID, err := GetUUIDFromUrl("GoalID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Goal ID", err)
		return
	}

	type request struct {
		Income_goal                       string `json:"income_goal"`
		Transaction_goal                  int    `json:"transaction_goal"`
		Estimated_average_sale_price      string `json:"estimated_average_sale_price"`
		Estimated_average_commission_rate string `json:"estimated_average_commission_rate"`
	}

	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	goal, err := cfg.DB.UpdateGoal(r.Context(), database.UpdateGoalParams{
		ID:                             goalUUID,
		IncomeGoal:                     sql.NullString{String: req.Income_goal, Valid: true},
		TransactionGoal:                sql.NullInt32{Int32: int32(req.Transaction_goal), Valid: true},
		EstimatedAverageSalePrice:      sql.NullString{String: req.Estimated_average_sale_price, Valid: true},
		EstimatedAverageCommissionRate: sql.NullString{String: req.Estimated_average_commission_rate, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update goal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, goal)
}
