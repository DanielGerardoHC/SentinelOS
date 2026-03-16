package handlers

import (
	"errors"
	"net/http"

	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

func BeginConfigHandler(w http.ResponseWriter, r *http.Request) {
	username := "admin" // TODO: obtener del contexto de autenticación
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
