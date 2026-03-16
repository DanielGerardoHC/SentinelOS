package validator

import (
	"net"
	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateDHCP(fw *model.Firewall) error {
	for _, d := range fw.DHCPConfigs {
		_, isIfc := fw.Interfaces[d.Interface]
		_, isVlan := fw.Vlans[d.Interface]

		if !isIfc && !isVlan {
			return &utils.APIError{Code: "ERR_NET_2003", Message: "Resource references unknown entity", Details: "dhcp references unknown interface " + d.Interface}
		}

		if net.ParseIP(d.StartIP) == nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: "invalid start-ip on dhcp " + d.Interface}
		}

		if net.ParseIP(d.EndIP) == nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: "invalid end-ip on dhcp " + d.Interface}
		}

		if net.ParseIP(d.Gateway) == nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: "invalid gateway on dhcp " + d.Interface}
		}
	}

	return nil
}
