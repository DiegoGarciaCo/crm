package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CollectLandingPageForm(w http.ResponseWriter, r *http.Request) {
	// Get User ID from context
	userID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: unable to get user ID", err)
		return
	}
	// Get Source from query parameters
	source := r.URL.Query().Get("source")
	if source == "" {
		http.Error(w, "Source parameter is required", http.StatusBadRequest)
		return
	}

	type LandingPageForm struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}

	var form LandingPageForm
	err = json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Start DB transaction
	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start database transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	// Insert form data into the database
	contact, err := qtx.LandingPageEmails(r.Context(), database.LandingPageEmailsParams{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Source:    sql.NullString{String: source, Valid: true},
		OwnerID:   uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save form data", err)
		return
	}

	// Insert email into the database
	err = qtx.EnterEmail(r.Context(), database.EnterEmailParams{
		ContactID:    uuid.NullUUID{UUID: contact.ID, Valid: true},
		EmailAddress: form.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save email", err)
		return
	}

	// Insert phone number into the database
	err = qtx.EnterPhoneNumber(r.Context(), database.EnterPhoneNumberParams{
		ContactID:   uuid.NullUUID{UUID: contact.ID, Valid: true},
		PhoneNumber: form.PhoneNumber,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save phone number", err)
		return
	}

	// create a JWT
	token, err := cfg.GenerateEmailToken(form.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate email token", err)
		return
	}

	// send email
	err = cfg.SendVerificationEmail(form.Email, token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send verification email", err)
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Form data collected successfully"})
}
