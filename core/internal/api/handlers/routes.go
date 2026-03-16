package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes, err := system.GetRoutes()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

func EditRouteHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid or missing ID", "invalid route id format")
		return
	}

	var route *model.Route
	for _, r := range fw.Routes {
		if r.ID == id {
			route = r
			break
		}
	}

	if route == nil {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "route ID "+idStr)
		return
	}

	var req struct {
		Destination string `json:"destination"`
		Gateway     string `json:"gateway"`
		Interface   string `json:"interface"`
		Metric      int    `json:"metric"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Destination != "" {
		route.Destination = req.Destination
	}
	if req.Gateway != "" {
		route.Gateway = req.Gateway
	}
	if req.Interface != "" {
		route.Interface = req.Interface
	}
	if req.Metric != 0 {
		route.Metric = req.Metric
	}
	if req.Description != "" {
		route.Description = req.Description
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "route updated in candidate"}`))
}

func CreateRouteHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		Destination string `json:"destination"`
		Gateway     string `json:"gateway"`
		Interface   string `json:"interface"`
		Metric      int    `json:"metric"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	id := config_engine.NextRouteID()

	route := &model.Route{
		ID:          id,
		Destination: req.Destination,
		Gateway:     req.Gateway,
		Interface:   req.Interface,
		Metric:      req.Metric,
		Description: req.Description,
	}

	fw.Routes = append(fw.Routes, route)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "route created in candidate"}`))
}

func DeleteRouteHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing route id in URL")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid or missing ID", "invalid route id format")
		return
	}

	index := -1
	for i, route := range fw.Routes {
		if route.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "route ID "+idParam)
		return
	}

	fw.Routes = append(fw.Routes[:index], fw.Routes[index+1:]...)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "route deleted from candidate"}`))
}
