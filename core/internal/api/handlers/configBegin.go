package handlers

import (
	"errors"
	"net/http"

	"sentinelos/core/internal/api/middleware"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

// BeginConfigHandler godoc
// @Summary Begin Configuration Session
// @Description Start a new candidate configuration session.
// @Tags system, config
// @Produce json
// @Success 200 {object} map[string]string "message: config session started"
// @Failure 401 {object} utils.APIError "Missing authorization token"
// @Failure 409 {object} utils.APIError "Config session already active"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Security ApiKeyAuth
// @Router /api/config/begin [post]
func BeginConfigHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.UserFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "ERR_SEC_5001", "Missing authorization token", "")
		return
	}
	username := claims.Username

	err := config_engine.BeginConfig(username)
	if err != nil {
		var apiErr *utils.APIError
		if errors.As(err, &apiErr) {
			utils.SendError(w, http.StatusConflict, apiErr.Code, apiErr.Message, apiErr.Details)
		} else {
			utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "config session started"}`))
}
