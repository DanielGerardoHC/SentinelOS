package handlers

import (
	"errors"
	"net/http"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

func CommitHandler(w http.ResponseWriter, r *http.Request) {
	err := config_engine.Commit()
	if err != nil {
		var apiErr *utils.APIError
		if errors.As(err, &apiErr) {
			utils.SendError(w, http.StatusBadRequest, apiErr.Code, apiErr.Message, apiErr.Details)
		} else {
			utils.SendError(w, http.StatusBadRequest, "ERR_SYS_4000", "Unknown commit error", err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "commit successful"}`))
}
