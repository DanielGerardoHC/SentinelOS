package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system"
	"sentinelos/core/internal/system/config_engine"
	"sentinelos/core/pkg/utils"
)

func AddressesHandler(w http.ResponseWriter, r *http.Request) {
	addresses, err := system.GetAddresses()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

func CreateAddressHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req struct {
		Name string   `json:"name"`
		IPs  []string `json:"ips"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "address name is required")
		return
	}

	if fw.Addresses == nil {
		fw.Addresses = make(map[string]*model.Address)
	}

	if _, exists := fw.Addresses[req.Name]; exists {
		utils.SendError(w, http.StatusConflict, "ERR_NET_2006", "Resource already exists", "address "+req.Name)
		return
	}

	address := &model.Address{
		Name: req.Name,
	}

	if req.IPs != nil {
		var parsedIPs []net.IPNet
		for _, ipStr := range req.IPs {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				utils.SendError(w, http.StatusBadRequest, "ERR_NET_1001", "Invalid IP or CIDR format", "failed to parse: "+ipStr)
				return
			}
			parsedIPs = append(parsedIPs, *ipNet)
		}
		address.IPs = parsedIPs
	}

	fw.Addresses[req.Name] = address

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "address created in candidate"}`))
}

func EditAddressHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing address name in URL")
		return
	}

	address, exists := fw.Addresses[nameParam]
	if !exists {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "address "+nameParam)
		return
	}

	var req struct {
		IPs []string `json:"ips"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
		return
	}

	if req.IPs != nil {
		var parsedIPs []net.IPNet
		for _, ipStr := range req.IPs {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				utils.SendError(w, http.StatusBadRequest, "ERR_NET_1001", "Invalid IP or CIDR format", "failed to parse: "+ipStr)
				return
			}
			parsedIPs = append(parsedIPs, *ipNet)
		}
		address.IPs = parsedIPs
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "address updated in candidate"}`))
}

func DeleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	nameParam := chi.URLParam(r, "name")
	if nameParam == "" {
		utils.SendError(w, http.StatusBadRequest, "ERR_NET_1002", "Missing required field", "missing address name in URL")
		return
	}

	if _, exists := fw.Addresses[nameParam]; !exists {
		utils.SendError(w, http.StatusNotFound, "ERR_NET_1003", "Resource not found", "address "+nameParam)
		return
	}

	delete(fw.Addresses, nameParam)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "address deleted from candidate"}`))
}
