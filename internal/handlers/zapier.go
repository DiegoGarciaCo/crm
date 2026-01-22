package handlers

import (
	"net/http"
)

func (cfg *apiCfg) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user, err := cfg.DB.GetUserNameById(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user info", err)
		return
	}

	type userInfoResponse struct {
		UserName string `json:"username"`
	}

	respondWithJSON(w, http.StatusOK, userInfoResponse{
		UserName: user,
	})
}
