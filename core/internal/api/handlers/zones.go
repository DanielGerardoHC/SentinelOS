package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

func ZonesHandler(w http.ResponseWriter, r *http.Request) {

	zones, err := system.GetZones()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}