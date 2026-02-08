package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	status := map[string]interface{}{
		"firewall":   system.FirewallRunning(),
		"interfaces": system.InterfacesCount(),
		"routes":     system.RoutesCount(),
		"dhcp":       system.DHCPRunning(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
