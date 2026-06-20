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

// AddressesHandler godoc
// @Summary List Addresses
// @Description Get a list of all addresses.
// @Tags networking, addresses
// @Produce json
// @Success 200 {object} map[string]model.Address "List of addresses"
// @Failure 500 {object} utils.APIError "ERR_SYS_4001 Internal server error"
// @Security ApiKeyAuth
// @Router /api/addresses [get]
func AddressesHandler(w http.ResponseWriter, r *http.Request) {
	addresses, err := system.GetAddresses()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "ERR_SYS_4001", "Internal server error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

// AddressCreateRequest represents the payload for creating an address
type AddressCreateRequest struct {
	Name string   `json:"name"`
	IPs  []string `json:"ips"`
}

// CreateAddressHandler godoc
// @Summary Create Address
// @Description Create a new address object.
// @Tags networking, addresses
// @Accept json
// @Produce json
// @Param request body AddressCreateRequest true "Address details"
// @Success 200 {object} map[string]string "message: address created in candidate"
// @Failure 400 {object} utils.APIError "ERR_NET_1004 Invalid JSON or Missing field"
// @Failure 409 {object} utils.APIError "ERR_NET_2006 Resource already exists"
// @Security ApiKeyAuth
// @Router /api/addresses [post]
func CreateAddressHandler(w http.ResponseWriter, r *http.Request) {
	fw := config_engine.GetCandidate()
	if fw == nil {
		utils.SendError(w, http.StatusBadRequest, "ERR_SYS_3001", "No active config session", "")
		return
	}

	var req AddressCreateRequest

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

// AddressEditRequest represents the payload for editing an address
type AddressEditRequest struct {
	IPs []string `json:"ips"`
}

// EditAddressHandler godoc
// @Summary Edit Address
// @Description Update an existing address object.
// @Tags networking, addresses
// @Accept json
// @Produce json
// @Param name path string true "Address Name"
// @Param request body AddressEditRequest true "Address details"
// @Success 200 {object} map[string]string "message: address updated in candidate"
// @Failure 400 {object} utils.APIError "Invalid JSON or Missing field"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/addresses/{name} [put]
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

	var req AddressEditRequest

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

// DeleteAddressHandler godoc
// @Summary Delete Address
// @Description Delete an existing address object.
// @Tags networking, addresses
// @Produce json
// @Param name path string true "Address Name"
// @Success 200 {object} map[string]string "message: address deleted from candidate"
// @Failure 400 {object} utils.APIError "Missing required field"
// @Failure 404 {object} utils.APIError "ERR_NET_1003 Resource not found"
// @Security ApiKeyAuth
// @Router /api/addresses/{name} [delete]
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
