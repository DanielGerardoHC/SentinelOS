package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

func InterfacesHandler(w http.ResponseWriter, r *http.Request) {

	ifaces, err := system.GetInterfaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ifaces)
}
