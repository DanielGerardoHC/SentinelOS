package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/api/middleware"
	"sentinelos/core/pkg/utils"
)

// MeHandler godoc
// @Summary Current User
// @Description Get current authenticated user details.
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{} "User details"
// @Failure 401 {object} utils.APIError "Missing authorization token"
// @Security ApiKeyAuth
// @Router /api/me [get]
func MeHandler(w http.ResponseWriter, r *http.Request) {

	claims, ok := middleware.UserFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "ERR_SEC_5001", "Missing authorization token", "")
		return
	}

	resp := map[string]interface{}{
		"username": claims.Username,
		"role":     claims.Role,
		"expires":  claims.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
