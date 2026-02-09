package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/system"
)

func NatHandler(w http.ResponseWriter, r *http.Request) {

	natRules, err := system.GetNatRules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(natRules)
}