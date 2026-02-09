package handlers 

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

func PoliciesHandler(w http.ResponseWriter, r *http.Request) {

	status, err := system.GetPolicies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}