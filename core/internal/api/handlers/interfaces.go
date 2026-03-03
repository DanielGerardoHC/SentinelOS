package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
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

type InterfaceEditRequest struct {
	IP         string   `json:"ip"`
	Zone       string   `json:"zone"`
	State      string   `json:"state"`
	Management []string `json:"management"`
}

func EditInterfaceHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	// parametro chi
	ifaceName := chi.URLParam(r, "name")

	iface, ok := fw.Interfaces[ifaceName]
	if !ok {
		http.Error(w, "interface not found", 404)
		return
	}

	var req InterfaceEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
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

	w.Write([]byte("interface updated in candidate"))
}
