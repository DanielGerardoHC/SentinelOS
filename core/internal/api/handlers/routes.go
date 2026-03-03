package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
)

func RoutesHandler(w http.ResponseWriter, r *http.Request) {

	routes, err := system.GetRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

func EditRouteHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid route id", 400)
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
		http.Error(w, "route not found", 404)
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
		http.Error(w, "invalid json", 400)
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

	w.Write([]byte("route updated in candidate"))
}

func CreateRouteHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
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
		http.Error(w, "invalid json", 400)
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

	w.Write([]byte("route created in candidate"))
}

func DeleteRouteHandler(w http.ResponseWriter, r *http.Request) {

	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	// obtener id desde la URL
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing route id", 400)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid route id", 400)
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
		http.Error(w, "route not found", 404)
		return
	}

	// eliminar del slice
	fw.Routes = append(fw.Routes[:index], fw.Routes[index+1:]...)

	w.Write([]byte("route deleted from candidate"))
}
