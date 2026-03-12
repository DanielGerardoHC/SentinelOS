package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NatHandler(w http.ResponseWriter, r *http.Request) {

	natRules, err := system.GetNatRules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(natRules)
}

func EditNatHandler(w http.ResponseWriter, r *http.Request) {
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

	var nat *model.NATRule
	for _, n := range fw.NATRules {
		if n.ID == id {
			nat = n
			break
		}
	}

	if nat == nil {
		http.Error(w, "NAT rule not found", 404)
		return
	}

	var req struct {
		SrcZone      string `json:"src-zone"`
		DstZone      string `json:"dst-zone"`
		OutInterface string `json:"out-interface"`
		Action       string `json:"action"`
		Description  string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if req.SrcZone != "" {
		zonePtr, exists := fw.Zones[req.SrcZone]
		if !exists {
			http.Error(w, "src-zone does not exist", 400)
			return
		}
		nat.SrcZone = zonePtr
	}

	if req.DstZone != "" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			http.Error(w, "dst-zone does not exist", 400)
			return
		}
		nat.DstZone = zonePtr
	}

	if req.Action != "" {
		action := model.NATAction(req.Action)
		if action != model.Masquerade && action != model.SNAT && action != model.DNAT {
			http.Error(w, "Invalid Action. Use masquerade, snat o dnat", 400)
			return
		}
		nat.Action = action
	}

	// datos strings
	if req.OutInterface != "" {
		nat.OutInterface = req.OutInterface
	}

	if req.Description != "" {
		nat.Description = req.Description
	}

	w.Write([]byte("NAT rule updated in candidate"))
}

func CreateNatHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	var req struct {
		SrcZone      string `json:"src-zone"`
		DstZone      string `json:"dst-zone"`
		OutInterface string `json:"out-interface"`
		Action       string `json:"action"`
		Description  string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	id := config_engine.NextNATID()

	natRule := &model.NATRule{
		ID:           id,
		OutInterface: req.OutInterface,
		Description:  req.Description,
	}
	if req.SrcZone != "" {
		zonePtr, exists := fw.Zones[req.SrcZone]
		if !exists {
			http.Error(w, "src-zone does not exist", 400)
			return
		}
		natRule.SrcZone = zonePtr
	}

	if req.DstZone != "" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			http.Error(w, "dst-zone does not exist", 400)
			return
		}
		natRule.DstZone = zonePtr
	}

	if req.Action != "" {
		action := model.NATAction(req.Action)
		if action != model.Masquerade && action != model.SNAT && action != model.DNAT {
			http.Error(w, "Invalid Action. Use masquerade, snat o dnat", 400)
			return
		}
		natRule.Action = action
	}
	fw.NATRules = append(fw.NATRules, natRule)

	w.Write([]byte("NAT rule added in candidate"))

}

func DeleteNatHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing route id", 400)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid NatRule id", 400)
		return
	}

	index := -1
	for i, natRule := range fw.NATRules {
		if natRule.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "NAT rule not found", 404)
		return
	}

	fw.NATRules = append(fw.NATRules[:index], fw.NATRules[index+1:]...)

	w.Write([]byte("NAT rule deleted in candidate"))

}
