package validator

import (
	"fmt"
	"net"

	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateVlans(fw *model.Firewall) error {
	for _, v := range fw.Vlans {
		if _, ok := fw.Interfaces[v.Parent]; !ok {
			return &utils.APIError{Code: "ERR_NET_2005", Message: "VLAN references unknown parent interface", Details: fmt.Sprintf("vlan %s parent interface %s not found", v.Name, v.Parent)}
		}

		if v.Zone != "" {
			if _, ok := fw.Zones[v.Zone]; !ok {
				return &utils.APIError{Code: "ERR_NET_2003", Message: "Resource references unknown entity", Details: fmt.Sprintf("vlan %s references unknown zone %s", v.Name, v.Zone)}
			}
		}

		if v.IP != "" {
			if _, _, err := net.ParseCIDR(v.IP); err != nil {
				return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: fmt.Sprintf("invalid IP on vlan %s", v.Name)}
			}
		}
	}

	return nil
}
