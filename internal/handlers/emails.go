package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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
