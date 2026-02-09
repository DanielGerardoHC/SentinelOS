package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)
func RoutesHandler(w http.ResponseWriter, r *http.Request) {

	routes, err := system.GetRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}