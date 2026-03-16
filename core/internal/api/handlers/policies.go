package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

func PoliciesHandler(w http.ResponseWriter, r *http.Request) {
	status, err := system.GetPolicies()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		SrcZone  string   `json:"src-zone"`
		DstZone  string   `json:"dst-zone"`
		SrcAddr  string   `json:"src-addr"`
		DstAddr  string   `json:"dst-addr"`
		Services []string `json:"services"`
		Action   string   `json:"action"`
		Log      bool     `json:"log"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	id := config_engine.NextPolicyID()

	policy := &model.Policy{
		ID:  id,
		Log: req.Log,
	}

	if req.SrcZone != "" {
		zonePtr, exists := fw.Zones[req.SrcZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-zone does not exist: "+req.SrcZone)
			return
		}
		policy.SrcZone = zonePtr
	}

	if req.DstZone != "" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-zone does not exist: "+req.DstZone)
			return
		}
		policy.DstZone = zonePtr
	}

	if req.SrcAddr != "" {
		addrPtr, exists := fw.Addresses[req.SrcAddr]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-addr does not exist: "+req.SrcAddr)
			return
		}
		policy.SrcAddr = addrPtr
	}

	if req.DstAddr != "" {
		addrPtr, exists := fw.Addresses[req.DstAddr]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-addr does not exist: "+req.DstAddr)
			return
		}
		policy.DstAddr = addrPtr
	}

	if req.Services != nil {
		var svcPtrs []*model.Service
		for _, svcName := range req.Services {
			svcPtr, exists := fw.Services[svcName]
			if !exists {
				utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "service does not exist: "+svcName)
				return
			}
			svcPtrs = append(svcPtrs, svcPtr)
		}
		policy.Services = svcPtrs
	}

	if req.Action != "" {
		action := model.Action(req.Action)
		if action != model.Allow && action != model.Deny {
			utils.SendError(w, http.StatusBadRequest, "ERR_SEC_1002", "Invalid policy action", "Use allow or deny")
			return
		}
		policy.Action = action
	}

	fw.Policies = append(fw.Policies, policy)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "policy created in candidate"}`))
}

func EditPolicyHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid or missing ID", "invalid policy id format")
		return
	}

	var policy *model.Policy
	for _, p := range fw.Policies {
		if p.ID == id {
			policy = p
			break
		}
	}

	if policy == nil {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", fmt.Sprintf("policy ID %d", id))
		return
	}

	var req struct {
		SrcZone  string   `json:"src-zone"`
		DstZone  string   `json:"dst-zone"`
		SrcAddr  string   `json:"src-addr"`
		DstAddr  string   `json:"dst-addr"`
		Services []string `json:"services"`
		Action   string   `json:"action"`
		Log      *bool    `json:"log"`
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
		policy.SrcZone = zonePtr
	}

	if req.DstZone != "" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-zone does not exist: "+req.DstZone)
			return
		}
		policy.DstZone = zonePtr
	}

	if req.SrcAddr != "" {
		addrPtr, exists := fw.Addresses[req.SrcAddr]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-addr does not exist: "+req.SrcAddr)
			return
		}
		policy.SrcAddr = addrPtr
	}

	if req.DstAddr != "" {
		addrPtr, exists := fw.Addresses[req.DstAddr]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-addr does not exist: "+req.DstAddr)
			return
		}
		policy.DstAddr = addrPtr
	}

	if req.Services != nil {
		var svcPtrs []*model.Service
		for _, svcName := range req.Services {
			svcPtr, exists := fw.Services[svcName]
			if !exists {
				utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "service does not exist: "+svcName)
				return
			}
			svcPtrs = append(svcPtrs, svcPtr)
		}
		policy.Services = svcPtrs
	}

	if req.Action != "" {
		action := model.Action(req.Action)
		if action != model.Allow && action != model.Deny {
			utils.SendError(w, http.StatusBadRequest, "ERR_SEC_1002", "Invalid policy action", "Use allow or deny")
			return
		}
		policy.Action = action
	}

	if req.Log != nil {
		policy.Log = *req.Log
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "policy updated in candidate"}`))
}

func DeletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing policy id in URL")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid or missing ID", "invalid policy id format")
		return
	}

	index := -1
	for i, policy := range fw.Policies {
		if policy.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", fmt.Sprintf("policy ID %d", id))
		return
	}

	fw.Policies = append(fw.Policies[:index], fw.Policies[index+1:]...)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "policy deleted from candidate"}`))
}
