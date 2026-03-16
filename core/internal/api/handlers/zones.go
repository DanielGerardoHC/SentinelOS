package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"

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
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		Name       string   `json:"name"`
		Type       string   `json:"type"`
		Interfaces []string `json:"interfaces"`
		Networks   []string `json:"networks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "zone name is required")
		return
	}

	if fw.Zones == nil {
		fw.Zones = make(map[string]*model.Zone)
	}

	if _, exists := fw.Zones[req.Name]; exists {
		utils.SendError(w, http.StatusConflict, "ERR_NET_2006", "Resource already exists", "zone "+req.Name)
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
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_1005", "Invalid zone type", "Use l2 or l3")
			return
		}
		zone.Type = zType
	}

	fw.Zones[req.Name] = zone

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "zone added in candidate"}`))
}

func EditZoneHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing zone name in URL")
		return
	}

	zone, exists := fw.Zones[nameParam]
	if !exists {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "zone "+nameParam)
		return
	}

	var req struct {
		Type       string   `json:"type"`
		Interfaces []string `json:"interfaces"`
		Networks   []string `json:"networks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Type != "" {
		zType := model.ZoneType(req.Type)
		if zType != model.ZoneL2 && zType != model.ZoneL3 {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_1005", "Invalid zone type", "Use l2 or l3")
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

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "zone updated in candidate"}`))
}

func DeleteZoneHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing zone name in URL")
		return
	}

	if _, exists := fw.Zones[nameParam]; !exists {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "zone "+nameParam)
		return
	}

	delete(fw.Zones, nameParam)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "zone deleted in candidate"}`))
}
