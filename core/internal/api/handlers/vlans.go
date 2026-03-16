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

func VlansHandler(w http.ResponseWriter, r *http.Request) {

	vlans, err := system.GetVlans()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vlans)
}

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

	var req struct {
		Parent     string   `json:"parent"`
		ID         int      `json:"id"`
		IP         string   `json:"ip"`
		Zone       string   `json:"zone"`
		State      string   `json:"state"`
		Management []string `json:"management"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
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
	if req.Zone != "" {
		vlan.Zone = req.Zone
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

func CreateVlanHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		Name       string   `json:"name"`
		Parent     string   `json:"parent"`
		ID         int      `json:"id"`
		IP         string   `json:"ip"`
		Zone       string   `json:"zone"`
		State      string   `json:"state"`
		Management []string `json:"management"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
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

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "vlan created in candidate"}`))
}

func DeleteVlanHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	vlanName := chi.URLParam(r, "name")

	_, ok := fw.Vlans[vlanName]
	if !ok {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "vlan "+vlanName)
		return
	}

	delete(fw.Vlans, vlanName)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "vlan deleted from candidate"}`))
}
