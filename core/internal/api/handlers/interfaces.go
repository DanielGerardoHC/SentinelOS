package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"

	"github.com/go-chi/chi/v5"
)

func InterfacesHandler(w http.ResponseWriter, r *http.Request) {
	ifaces, err := system.GetInterfaces()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ifaces)
}

type InterfaceEditRequest struct {
	IP         string   `json:"ip"`
	Zone       string   `json:"zone"`
	State      string   `json:"state"`
	Management []string `json:"management"`
}

func EditInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	ifaceName := chi.URLParam(r, "name")
	iface, ok := fw.Interfaces[ifaceName]
	if !ok {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "interface "+ifaceName)
		return
	}

	var req InterfaceEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.IP != "" {
		iface.IP = req.IP
	}
	if req.Zone != "" {
		iface.Zone = req.Zone
	}
	if req.State != "" {
		iface.State = req.State
	}
	if req.Management != nil {
		iface.Management = req.Management
	}

	// buena practica devolver JSON en los casos de exito
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "interface updated in candidate"}`))
}
