package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := system.GetServices()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Ports    []int  `json:"ports"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "service name is required")
		return
	}

	if fw.Services == nil {
		fw.Services = make(map[string]*model.Service)
	}

	if _, exists := fw.Services[req.Name]; exists {
		utils.SendError(w, http.StatusConflict, "ERR_NET_2006", "Resource already exists", "service "+req.Name)
		return
	}

	service := &model.Service{
		Name:  req.Name,
		Ports: req.Ports,
	}

	if req.Protocol != "" {
		proto := model.Protocol(req.Protocol)
		if proto != model.TCP && proto != model.UDP {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_1006", "Invalid protocol", "Use tcp or udp")
			return
		}
		service.Protocol = proto
	}

	fw.Services[req.Name] = service

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "service created in candidate"}`))
}

func EditServiceHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing service name in URL")
		return
	}

	service, exists := fw.Services[nameParam]
	if !exists {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "service "+nameParam)
		return
	}

	var req struct {
		Protocol string `json:"protocol"`
		Ports    []int  `json:"ports"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Protocol != "" {
		proto := model.Protocol(req.Protocol)
		if proto != model.TCP && proto != model.UDP {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_1006", "Invalid protocol", "Use tcp or udp")
			return
		}
		service.Protocol = proto
	}

	if req.Ports != nil {
		service.Ports = req.Ports
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "service updated in candidate"}`))
}

func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing service name in URL")
		return
	}

	if _, exists := fw.Services[nameParam]; !exists {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "service "+nameParam)
		return
	}

	delete(fw.Services, nameParam)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "service deleted from candidate"}`))
}
