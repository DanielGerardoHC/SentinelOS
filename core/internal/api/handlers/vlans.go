package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system/config_engine"

	"sentinelos/core/internal/system"

	"github.com/go-chi/chi/v5"
)

func VlansHandler(w http.ResponseWriter, r *http.Request) {

	vlans, err := system.GetVlans()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vlans)
}

func EditVlanHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	vlanName := chi.URLParam(r, "name")

	vlan, ok := fw.Vlans[vlanName]
	if !ok {
		http.Error(w, "interface not found", 404)
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
		http.Error(w, "invalid json", 400)
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

	w.Write([]byte("vlan updated in candidate"))
}

func CreateVlanHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
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
		http.Error(w, "invalid json", 400)
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

	fw.Vlans[vlan.Name] = vlan

	w.Write([]byte("vlan created in candidate"))
}

func DeleteVlanHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	vlanName := chi.URLParam(r, "name")

	_, ok := fw.Vlans[vlanName]
	if !ok {
		http.Error(w, "interface not found", 404)
		return
	}

	delete(fw.Vlans, vlanName)

	w.Write([]byte("vlan deleted from candidate"))
}
