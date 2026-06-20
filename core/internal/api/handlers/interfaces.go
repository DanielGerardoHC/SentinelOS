package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"

	"github.com/go-chi/chi/v5"
)

// InterfacesHandler godoc
// @Summary List Interfaces
// @Description Get a list of all network interfaces.
// @Tags networking, interfaces
// @Accept json
// @Produce json
// @Param layer query string false "Filter by OSI layer (2 or 3)"
// @Success 200 {array} system.InterfaceInfo "List of interfaces"
// @Failure 500 {object} utils.APIError "ERR_SYS_4001 Internal server error"
// @Security ApiKeyAuth
// @Router /api/interfaces [get]
func InterfacesHandler(w http.ResponseWriter, r *http.Request) {
	layerFilter := r.URL.Query().Get("layer")

	var ifaces []system.InterfaceInfo
	var err error

	if layerFilter == "3" {
		ifaces, err = system.GetInterfacesL3()
	} else if layerFilter == "2" {
		ifaces, err = system.GetInterfacesL2()
	} else {
		ifaces, err = system.GetInterfaces()
	}

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

// EditInterfaceHandler godoc
// @Summary Edit Interface
// @Description Update an existing network interface.
// @Tags networking, interfaces
// @Accept json
// @Produce json
// @Param name path string true "Interface Name"
// @Param request body InterfaceEditRequest true "Interface details"
// @Success 200 {object} map[string]string "message: interface updated in candidate"
// @Failure 400 {object} utils.APIError "Invalid JSON or Zone does not exist"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/interfaces/{name} [put]
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

	if req.Zone != "" && req.Zone != iface.Zone {
		newZone, exists := fw.Zones[req.Zone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Zone does not exist", req.Zone)
			return
		}

		if iface.Zone != "" {
			if oldZone, oldExists := fw.Zones[iface.Zone]; oldExists {
				removeInterfaceFromZone(oldZone, ifaceName)
			}
		}

		addInterfaceToZone(newZone, ifaceName)
		iface.Zone = req.Zone
	}

	if req.IP != "" {
		iface.IP = req.IP
	}
	if req.State != "" {
		iface.State = req.State
	}
	if req.Management != nil {
		iface.Management = req.Management
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "interface updated in candidate"}`))
}
