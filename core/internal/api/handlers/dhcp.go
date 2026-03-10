package handlers

import (
	"encoding/json"
	"net/http"
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system/config_engine"

	"sentinelos/core/internal/system"

	"github.com/go-chi/chi/v5"
)

func DhcpHandler(w http.ResponseWriter, r *http.Request) {

	dhcp, err := system.GetDhcpInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dhcp)
}

func EditDhcpHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	ifaceName := chi.URLParam(r, "interface")
	if ifaceName == "" {
		http.Error(w, "invalid interface", 400)
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
		http.Error(w, "dhcp config not found", 404)
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
		http.Error(w, "invalid json", 400)
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
	if req.Dns == nil || len(req.Dns) == 0 {
		dhcpPool.DNS = req.Dns
	}
	if req.Lease_time == 0 {
		dhcpPool.LeaseTimeMin = req.Lease_time
	}

	w.Write([]byte("dhcp pool updated in candidate"))
}

func CreateDhcpHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
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
		http.Error(w, "invalid json", 400)
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

	w.Write([]byte("dhcp pool updated in candidate"))
}

func DeleteDhcpHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		http.Error(w, "no active config session", 400)
		return
	}

	ifaceName := chi.URLParam(r, "interface")
	if ifaceName == "" {
		http.Error(w, "invalid interface", 400)
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
		http.Error(w, "dhcp config not found", 404)
		return
	}

	fw.DHCPConfigs = append(fw.DHCPConfigs[:index], fw.DHCPConfigs[index+1:]...)

	w.Write([]byte("dhcp pool updated in candidate"))

}
