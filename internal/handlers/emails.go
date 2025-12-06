package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiCfg) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return cfg.EmailSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithError(w, http.StatusBadRequest, "Invalid or expired token", err)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["typ"] != "email_verification" {
		respondWithError(w, http.StatusBadRequest, "Invalid token type", nil)
		return
	}

	email, _ := claims["email"].(string)

	// mark verified in DB
	err = cfg.DB.VerifyEmail(r.Context(), email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to verify email", err)
		return
	}

	// redirect to success page
	http.Redirect(w, r, "/email-verified", http.StatusSeeOther)
}
