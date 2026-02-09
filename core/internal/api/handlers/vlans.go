package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

func VlansHandler(w http.ResponseWriter, r *http.Request) {

	vlans, err := system.GetVlans()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vlans)
}