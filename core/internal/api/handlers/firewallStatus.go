package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
)

// StatusHandler godoc
// @Summary Firewall Status
// @Description Get the system and firewall status.
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{} "Status details"
// @Security ApiKeyAuth
// @Router /api/status [get]
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
