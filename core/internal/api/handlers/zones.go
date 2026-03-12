package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system/config_engine"

	"sentinelos/core/internal/system"

	"github.com/go-chi/chi/v5"
)

func ZonesHandler(w http.ResponseWriter, r *http.Request) {

	zones, err := system.GetZones()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

func CreateZoneHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	var req struct {
		Name       string   `json:"name"`
		Type       string   `json:"type"`
		Interfaces []string `json:"interfaces"`
		Networks   []string `json:"networks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if req.Name == "" {
		http.Error(w, "zone name is required", 400)
		return
	}

	if _, exists := fw.Zones[req.Name]; exists {
		http.Error(w, "zone already exists", 409)
		return
	}

	zone := &model.Zone{
		Name:       req.Name,
		Interfaces: req.Interfaces,
		Networks:   req.Networks,
	}

	if req.Type != "" {
		zType := model.ZoneType(req.Type)
		if zType != model.ZoneL2 && zType != model.ZoneL3 {
			http.Error(w, "invalid zone type. Use l2 or l3", 400)
			return
		}
		zone.Type = zType
	}

	fw.Zones[req.Name] = zone

	w.Write([]byte("zone added in candidate"))
}

func EditZoneHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		http.Error(w, "missing zone name", 400)
		return
	}

	zone, exists := fw.Zones[nameParam]
	if !exists {
		http.Error(w, "zone not found", 404)
		return
	}

	var req struct {
		Type       string   `json:"type"`
		Interfaces []string `json:"interfaces"`
		Networks   []string `json:"networks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if req.Type != "" {
		zType := model.ZoneType(req.Type)
		if zType != model.ZoneL2 && zType != model.ZoneL3 {
			http.Error(w, "invalid zone type. Use l2 or l3", 400)
			return
		}
		zone.Type = zType
	}

	if req.Interfaces != nil {
		zone.Interfaces = req.Interfaces
	}

	if req.Networks != nil {
		zone.Networks = req.Networks
	}

	w.Write([]byte("zone updated in candidate"))
}

func DeleteZoneHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		http.Error(w, "missing zone name", 400)
		return
	}

	if _, exists := fw.Zones[nameParam]; !exists {
		http.Error(w, "zone not found", 404)
		return
	}

	delete(fw.Zones, nameParam)

	w.Write([]byte("zone deleted in candidate"))
}
