package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (cfg *apiCfg) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return cfg.EmailSecret, nil
	})
	if err != nil || !token.Valid {
		// Redirect to frontend with error
		http.Redirect(w, r, "https://access.soldbyghost.com/email-verified?status=error", http.StatusSeeOther)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims["typ"] != "email_verification" {
		http.Redirect(w, r, "https://access.soldbyghost.com/email-verified?status=error", http.StatusSeeOther)
		return
	}
	email, _ := claims["email"].(string)
	// mark verified in DB
	err = cfg.DB.VerifyEmail(r.Context(), email)
	if err != nil {
		http.Redirect(w, r, "https://access.soldbyghost.com/email-verified?status=error", http.StatusSeeOther)
		return
	}
	// redirect to success page on frontend
	http.Redirect(w, r, "https://access.soldbyghost.com/email-verified?status=success", http.StatusSeeOther)
}

func (cfg *apiCfg) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	// Get email from request body
	type ResendEmailForm struct {
		Email string `json:"email"`
	}
	var form ResendEmailForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
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

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Verification email sent"})
}

func (cfg *apiCfg) CreateEmailAddress(w http.ResponseWriter, r *http.Request) {
	type emailRequest struct {
		Email     string `json:"email"`
		Type      string `json:"type"`
		IsPrimary bool   `json:"is_primary"`
	}

	// Get Contact ID from URL params
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Contact ID", err)
		return
	}

	// Parse request body
	var req emailRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Create email address in DB
	err = cfg.DB.EnterEmail(r.Context(), database.EnterEmailParams{
		ContactID:    uuid.NullUUID{UUID: contactUUID, Valid: true},
		EmailAddress: req.Email,
		Type:         sql.NullString{String: req.Type, Valid: req.Type != ""},
		IsPrimary:    sql.NullBool{Bool: req.IsPrimary, Valid: req.IsPrimary},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create email address", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) UpdateEmailAddress(w http.ResponseWriter, r *http.Request) {
	type emailRequest struct {
		Email     string `json:"email"`
		Type      string `json:"type"`
		IsPrimary bool   `json:"is_primary"`
	}

	// Get Email ID from URL params
	emailUUID, err := GetUUIDFromUrl("emailID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Email ID", err)
		return
	}

	// Parse request body
	var req emailRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Update email address in DB
	err = cfg.DB.UpdateEmail(r.Context(), database.UpdateEmailParams{
		ID:           emailUUID,
		EmailAddress: req.Email,
		Type:         sql.NullString{String: req.Type, Valid: req.Type != ""},
		IsPrimary:    sql.NullBool{Bool: req.IsPrimary, Valid: req.IsPrimary},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update email address", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiCfg) DeleteEmailAddress(w http.ResponseWriter, r *http.Request) {
	// Get Email ID from URL params
	emailUUID, err := GetUUIDFromUrl("emailID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Email ID", err)
		return
	}

	// Delete email address from DB
	err = cfg.DB.DeleteEmail(r.Context(), emailUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete email address", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
