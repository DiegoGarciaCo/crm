package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) GetCollaborators(w http.ResponseWriter, r *http.Request) {
	type req struct {
		OrgIDs []string `json:"org_ids"`
	}

	var request req
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	var collaborators []database.Collaborator
	// Fetch collaborators based on organization IDs
	for _, orgID := range request.OrgIDs {
		// parse orgID if necessary
		orgUUID, err := uuid.Parse(orgID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid organization ID", err)
			return
		}

		collaborators, err := cfg.DB.GetOrganizationMembers(r.Context(), orgUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch collaborators", err)
			return
		}
		// Append or process collaborators as needed
		collaborators = append(collaborators, collaborators...)
	}

	cfg.logger.Info("collaborators fetched", "collaborators", collaborators)

	if len(collaborators) == 0 {
		emptyList := []database.Collaborator{}
		respondWithJSON(w, http.StatusOK, emptyList)
		return
	}

	// Respond with the list of collaborators
	respondWithJSON(w, http.StatusOK, collaborators)
}
