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

func DhcpHandler(w http.ResponseWriter, r *http.Request) {
	dhcp, err := system.GetDhcpInfo()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dhcp)
}

func EditDhcpHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	ifaceName := chi.URLParam(r, "interface")
	if ifaceName == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing interface name in URL")
		return
	}

	var dhcpPool *model.DHCP
	for _, d := range fw.DHCPConfigs {
		if d.Interface == ifaceName {
			dhcpPool = d
			break
		}
	}

	if dhcpPool == nil {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "dhcp pool for interface "+ifaceName)
		return
	}

	var req struct {
		Start_ip   string   `json:"start_ip"`
		End_ip     string   `json:"end_ip"`
		Gateway    string   `json:"gateway"`
		Dns        []string `json:"dns"`
		Lease_time int      `json:"lease_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Start_ip != "" {
		dhcpPool.StartIP = req.Start_ip
	}
	if req.End_ip != "" {
		dhcpPool.EndIP = req.End_ip
	}
	if req.Gateway != "" {
		dhcpPool.Gateway = req.Gateway
	}
	if req.Dns != nil || len(req.Dns) != 0 {
		dhcpPool.DNS = req.Dns
	}
	if req.Lease_time != 0 {
		dhcpPool.LeaseTimeMin = req.Lease_time
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "dhcp pool updated in candidate"}`))
}

func CreateDhcpHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		Interface  string   `json:"interface"`
		Start_ip   string   `json:"start_ip"`
		End_ip     string   `json:"end_ip"`
		Gateway    string   `json:"gateway"`
		Dns        []string `json:"dns"`
		Lease_time int      `json:"lease_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	dhcpPool := &model.DHCP{
		Interface:    req.Interface,
		StartIP:      req.Start_ip,
		EndIP:        req.End_ip,
		Gateway:      req.Gateway,
		DNS:          req.Dns,
		LeaseTimeMin: req.Lease_time,
	}

	fw.DHCPConfigs = append(fw.DHCPConfigs, dhcpPool)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "dhcp pool created in candidate"}`))
}

func DeleteDhcpHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	ifaceName := chi.URLParam(r, "interface")
	if ifaceName == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing interface name in URL")
		return
	}

	var dhcpPool *model.DHCP
	index := -1
	for i, d := range fw.DHCPConfigs {
		if d.Interface == ifaceName {
			dhcpPool = d
			index = i
			break
		}
	}

	if dhcpPool == nil {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "dhcp pool for interface "+ifaceName)
		return
	}

	fw.DHCPConfigs = append(fw.DHCPConfigs[:index], fw.DHCPConfigs[index+1:]...)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "dhcp pool deleted in candidate"}`))
}
