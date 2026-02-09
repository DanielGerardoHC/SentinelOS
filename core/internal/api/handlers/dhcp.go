package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

func DhcpHandler(w http.ResponseWriter, r *http.Request) {

	dhcp, err := system.GetDhcpInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dhcp)
}
