package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"

	"github.com/go-chi/chi/v5"
)

// VlansHandler godoc
// @Summary List Vlans
// @Description Get a list of all vlans.
// @Tags networking, vlans
// @Produce json
// @Success 200 {object} map[string]model.Vlan "List of vlans"
// @Failure 500 {object} utils.APIError "ERR_SYS_4001 Internal server error"
// @Security ApiKeyAuth
// @Router /api/vlans [get]
func VlansHandler(w http.ResponseWriter, r *http.Request) {

	vlans, err := system.GetVlans()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vlans)
}

// VlanCreateRequest represents the payload for creating a vlan
type VlanCreateRequest struct {
	Name       string   `json:"name"`
	Parent     string   `json:"parent"`
	ID         int      `json:"id"`
	IP         string   `json:"ip"`
	Zone       string   `json:"zone"`
	State      string   `json:"state"`
	Management []string `json:"management"`
}

// CreateVlanHandler godoc
// @Summary Create Vlan
// @Description Create a new vlan object.
// @Tags networking, vlans
// @Accept json
// @Produce json
// @Param request body VlanCreateRequest true "Vlan details"
// @Success 200 {object} map[string]string "message: vlan created in candidate"
// @Failure 400 {object} utils.APIError "Invalid JSON or Zone does not exist"
// @Security ApiKeyAuth
// @Router /api/vlans [post]
func CreateVlanHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req VlanCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Zone != "" {
		if _, exists := fw.Zones[req.Zone]; !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Zone does not exist", req.Zone)
			return
		}
	}

	vlan := &model.Vlan{
		Name:       req.Name,
		Parent:     req.Parent,
		ID:         req.ID,
		IP:         req.IP,
		Zone:       req.Zone,
		State:      req.State,
		Management: req.Management,
	}

	if fw.Vlans == nil {
		fw.Vlans = make(map[string]*model.Vlan)
	}

	fw.Vlans[vlan.Name] = vlan

	if vlan.Zone != "" {
		addInterfaceToZone(fw.Zones[vlan.Zone], vlan.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "vlan created in candidate"}`))
}

// VlanEditRequest represents the payload for editing a vlan
type VlanEditRequest struct {
	Parent     string   `json:"parent"`
	ID         int      `json:"id"`
	IP         string   `json:"ip"`
	Zone       string   `json:"zone"`
	State      string   `json:"state"`
	Management []string `json:"management"`
}

// EditVlanHandler godoc
// @Summary Edit Vlan
// @Description Update an existing vlan object.
// @Tags networking, vlans
// @Accept json
// @Produce json
// @Param name path string true "Vlan Name"
// @Param request body VlanEditRequest true "Vlan details"
// @Success 200 {object} map[string]string "message: vlan updated in candidate"
// @Failure 400 {object} utils.APIError "Invalid JSON or Zone does not exist"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/vlans/{name} [put]
func EditVlanHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	vlanName := chi.URLParam(r, "name")

	vlan, ok := fw.Vlans[vlanName]
	if !ok {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "vlan "+vlanName)
		return
	}

	var req VlanEditRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Zone != "" && req.Zone != vlan.Zone {
		newZone, exists := fw.Zones[req.Zone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Zone does not exist", req.Zone)
			return
		}

		if vlan.Zone != "" {
			if oldZone, oldExists := fw.Zones[vlan.Zone]; oldExists {
				removeInterfaceFromZone(oldZone, vlanName)
			}
		}

		addInterfaceToZone(newZone, vlanName)
		vlan.Zone = req.Zone
	}

	if req.Parent != "" {
		vlan.Parent = req.Parent
	}
	if req.ID != 0 {
		vlan.ID = req.ID
	}
	if req.IP != "" {
		vlan.IP = req.IP
	}
	if req.State != "" {
		vlan.State = req.State
	}
	if req.Management != nil {
		vlan.Management = req.Management
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "vlan updated in candidate"}`))
}

// DeleteVlanHandler godoc
// @Summary Delete Vlan
// @Description Delete an existing vlan object.
// @Tags networking, vlans
// @Produce json
// @Param name path string true "Vlan Name"
// @Success 200 {object} map[string]string "message: vlan deleted from candidate"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/vlans/{name} [delete]
func DeleteVlanHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	vlanName := chi.URLParam(r, "name")

	vlan, ok := fw.Vlans[vlanName]
	if !ok {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "vlan "+vlanName)
		return
	}

	if vlan.Zone != "" {
		if z, zExists := fw.Zones[vlan.Zone]; zExists {
			removeInterfaceFromZone(z, vlanName)
		}
	}

	delete(fw.Vlans, vlanName)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "vlan deleted from candidate"}`))
}

func addInterfaceToZone(zone *model.Zone, ifaceName string) {
	if zone == nil {
		return
	}
	for _, name := range zone.Interfaces {
		if name == ifaceName {
			return
		}
	}
	zone.Interfaces = append(zone.Interfaces, ifaceName)
}

func removeInterfaceFromZone(zone *model.Zone, ifaceName string) {
	if zone == nil {
		return
	}
	for i, name := range zone.Interfaces {
		if name == ifaceName {
			zone.Interfaces = append(zone.Interfaces[:i], zone.Interfaces[i+1:]...)
			break
		}
	}
}
