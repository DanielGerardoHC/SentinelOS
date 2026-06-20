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

// NatHandler godoc
// @Summary List NAT Rules
// @Description Get a list of all NAT rules.
// @Tags security, nat
// @Produce json
// @Param type query string false "Filter by NAT type (e.g. snat, dnat-ip)"
// @Success 200 {object} []model.NATRule "List of NAT rules"
// @Failure 500 {object} utils.APIError "ERR_SYS_4001 Internal server error"
// @Security ApiKeyAuth
// @Router /api/nat [get]
func NatHandler(w http.ResponseWriter, r *http.Request) {
	actionFilter := r.URL.Query().Get("type")
	natRules, err := system.GetNatRules(actionFilter)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(natRules)
}

// NatCreateRequest represents the payload for creating a NAT rule
type NatCreateRequest struct {
	Type           string `json:"type"`
	SrcZone        string `json:"src-zone"`
	DstZone        string `json:"dst-zone"`
	SrcAddress     string `json:"src-addr"`
	DstAddress     string `json:"dst-addr"`
	Service        string `json:"service"`
	OutInterface   string `json:"out-interface"`
	TranslatedIP   string `json:"translated-ip"`
	TranslatedPort string `json:"translated-port"`
	Description    string `json:"description"`
}

// CreateNatHandler godoc
// @Summary Create NAT Rule
// @Description Create a new NAT rule.
// @Tags security, nat
// @Accept json
// @Produce json
// @Param request body NatCreateRequest true "NAT details"
// @Success 200 {object} map[string]string "message: NAT rule added in candidate"
// @Failure 400 {object} utils.APIError "Invalid JSON or NAT type"
// @Security ApiKeyAuth
// @Router /api/nat [post]
func CreateNatHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req NatCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	
	natType := model.NATType(req.Type)
	if natType != model.TypeSNAT && natType != model.TypeDNATIP && natType != model.TypeDNATPort {
		utils.SendError(w, http.StatusBadRequest, "ERR_SEC_1001", "Invalid NAT type", "Use snat, dnat-ip or dnat-port")
		return
	}

	id := config_engine.NextNATID()

	natRule := &model.NATRule{
		ID:             id,
		Type:           natType,
		SrcAddress:     req.SrcAddress,
		DstAddress:     req.DstAddress,
		Service:        req.Service,
		OutInterface:   req.OutInterface,
		TranslatedIP:   req.TranslatedIP,
		TranslatedPort: req.TranslatedPort,
		Description:    req.Description,
	}

	if req.SrcZone != "" && req.SrcZone != "any" {
		zonePtr, exists := fw.Zones[req.SrcZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-zone does not exist: "+req.SrcZone)
			return
		}
		natRule.SrcZone = zonePtr
	}

	if req.DstZone != "" && req.DstZone != "any" {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-zone does not exist: "+req.DstZone)
			return
		}
		natRule.DstZone = zonePtr
	}

	fw.NATRules = append(fw.NATRules, natRule)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "NAT rule added in candidate"}`))
}

// NatEditRequest represents the payload for editing a NAT rule
type NatEditRequest struct {
	Type           string `json:"type"`
	SrcZone        string `json:"src-zone"`
	DstZone        string `json:"dst-zone"`
	SrcAddress     string `json:"src-addr"`
	DstAddress     string `json:"dst-addr"`
	Service        string `json:"service"`
	OutInterface   string `json:"out-interface"`
	TranslatedIP   string `json:"translated-ip"`
	TranslatedPort string `json:"translated-port"`
	Description    string `json:"description"`
}

// EditNatHandler godoc
// @Summary Edit NAT Rule
// @Description Update an existing NAT rule.
// @Tags security, nat
// @Accept json
// @Produce json
// @Param id path int true "NAT Rule ID"
// @Param request body NatEditRequest true "NAT details"
// @Success 200 {object} map[string]string "message: NAT rule updated in candidate"
// @Failure 400 {object} utils.APIError "Invalid JSON or ID"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/nat/{id} [put]
func EditNatHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid or missing ID", "invalid nat id format")
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

	var req NatEditRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Type != "" {
		natType := model.NATType(req.Type)
		if natType != model.TypeSNAT && natType != model.TypeDNATIP && natType != model.TypeDNATPort {
			utils.SendError(w, http.StatusBadRequest, "ERR_SEC_1001", "Invalid NAT type", "Use snat, dnat-ip or dnat-port")
			return
		}
		nat.Type = natType
	}

if req.SrcZone == "" || req.SrcZone == "any" {
		nat.SrcZone = nil
	} else {
		zonePtr, exists := fw.Zones[req.SrcZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "src-zone does not exist: "+req.SrcZone)
			return
		}
		nat.SrcZone = zonePtr
	}

	if req.DstZone == "" || req.DstZone == "any" {
		nat.DstZone = nil
	} else {
		zonePtr, exists := fw.Zones[req.DstZone]
		if !exists {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_2003", "Resource references unknown entity", "dst-zone does not exist: "+req.DstZone)
			return
		}
		nat.DstZone = zonePtr
	}

	if req.SrcAddress != "" { nat.SrcAddress = req.SrcAddress }
	if req.DstAddress != "" { nat.DstAddress = req.DstAddress }
	if req.Service != "" { nat.Service = req.Service }
	if req.OutInterface != "" { nat.OutInterface = req.OutInterface }
	if req.TranslatedIP != "" { nat.TranslatedIP = req.TranslatedIP }
	if req.TranslatedPort != "" { nat.TranslatedPort = req.TranslatedPort }
	if req.Description != "" { nat.Description = req.Description }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "NAT rule updated in candidate"}`))
}

// DeleteNatHandler godoc
// @Summary Delete NAT Rule
// @Description Delete an existing NAT rule.
// @Tags security, nat
// @Produce json
// @Param id path int true "NAT Rule ID"
// @Success 200 {object} map[string]string "message: NAT rule deleted in candidate"
// @Failure 400 {object} utils.APIError "Invalid ID"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/nat/{id} [delete]
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


// NatMoveRequest represents the payload for moving a NAT rule
type NatMoveRequest struct {
	Position    string `json:"position"`
	ReferenceID int    `json:"reference_id"`
}

// MoveNatHandler godoc
// @Summary Move NAT Rule
// @Description Move a NAT rule relative to another rule (top, bottom, before, after).
// @Tags security, nat
// @Accept json
// @Produce json
// @Param id path int true "NAT Rule ID to move"
// @Param request body NatMoveRequest true "Move details"
// @Success 200 {object} map[string]string "message: NAT rule moved successfully"
// @Failure 400 {object} utils.APIError "Invalid JSON or Position"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 NAT Rule not found"
// @Security ApiKeyAuth
// @Router /api/nat/{id}/move [post]
func MoveNatHandler(w http.ResponseWriter, r *http.Request) {
    fw := config_engine.GetCandidate()
    if fw == nil {
        utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
        return
    }

    idStr := chi.URLParam(r, "id")
    targetID, err := strconv.Atoi(idStr)
    if err != nil {
        utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Invalid ID", "invalid nat rule id format")
        return
    }

    var req NatMoveRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
        return
    }

    srcIndex := -1
    for i, nat := range fw.NATRules {
        if nat.ID == targetID {
            srcIndex = i
            break
        }
    }

    if srcIndex == -1 {
        utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "NAT Rule not found", "")
        return
    }


    ruleToMove := fw.NATRules[srcIndex]
    fw.NATRules = append(fw.NATRules[:srcIndex], fw.NATRules[srcIndex+1:]...)

    var insertIndex int

    switch req.Position {
    case "top":
        insertIndex = 0
    case "bottom":
        insertIndex = len(fw.NATRules)
    case "before", "after":
        refIndex := -1
        for i, nat := range fw.NATRules {
            if nat.ID == req.ReferenceID {
                refIndex = i
                break
            }
        }

        if refIndex == -1 {
            utils.SendError(w, http.StatusBadRequest, "ERR_NET_1005", "Reference NAT rule not found", "")
            return
        }

        if req.Position == "before" {
            insertIndex = refIndex
        } else {
            insertIndex = refIndex + 1
        }
    default:
        utils.SendError(w, http.StatusBadRequest, "ERR_NET_1006", "Invalid position parameter", "Use top, bottom, before, or after")
        return
    }

   
    newRules := make([]*model.NATRule, 0, len(fw.NATRules)+1)
    newRules = append(newRules, fw.NATRules[:insertIndex]...)
    newRules = append(newRules, ruleToMove)
    newRules = append(newRules, fw.NATRules[insertIndex:]...)

    fw.NATRules = newRules

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "NAT rule moved successfully"}`))
}