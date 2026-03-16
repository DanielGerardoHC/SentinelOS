package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
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

func CreateNatHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
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
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
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
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-zone does not exist: "+req.SrcZone)
			return
		}
		natRule.SrcZone = zonePtr
	}

	if req.DstZone != "" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-zone does not exist: "+req.DstZone)
			return
		}
		natRule.DstZone = zonePtr
	}

	if req.Action != "" {
		action := model.NATAction(req.Action)
		if action != model.Masquerade && action != model.SNAT && action != model.DNAT {
			utils.SendError(w, http.StatusBadRequest, "ERR_SEC_1001", "Invalid NAT action", "Use masquerade, snat or dnat")
			return
		}
		natRule.Action = action
	}

	fw.NATRules = append(fw.NATRules, natRule)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "NAT rule added in candidate"}`))
}

func EditNatHandler(w http.ResponseWriter, r *http.Request) {
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

	var nat *model.NATRule
	for _, n := range fw.NATRules {
		if n.ID == id {
			nat = n
			break
		}
	}

	if nat == nil {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", fmt.Sprintf("NAT rule ID %d", id))
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
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.SrcZone != "" {
		zonePtr, exists := fw.Zones[req.SrcZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-zone does not exist: "+req.SrcZone)
			return
		}
		nat.SrcZone = zonePtr
	}

	if req.DstZone != "" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-zone does not exist: "+req.DstZone)
			return
		}
		nat.DstZone = zonePtr
	}

	if req.Action != "" {
		action := model.NATAction(req.Action)
		if action != model.Masquerade && action != model.SNAT && action != model.DNAT {
			utils.SendError(w, http.StatusBadRequest, "ERR_SEC_1001", "Invalid NAT action", "Use masquerade, snat or dnat")
			return
		}
		nat.Action = action
	}

	if req.OutInterface != "" {
		nat.OutInterface = req.OutInterface
	}

	if req.Description != "" {
		nat.Description = req.Description
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "NAT rule updated in candidate"}`))
}

func DeleteNatHandler(w http.ResponseWriter, r *http.Request) {
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
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid or missing ID", "invalid NatRule id format")
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
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", fmt.Sprintf("NAT rule ID %d", id))
		return
	}

	fw.NATRules = append(fw.NATRules[:index], fw.NATRules[index+1:]...)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "NAT rule deleted in candidate"}`))
}
