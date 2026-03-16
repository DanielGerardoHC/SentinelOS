package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateInterfaces(fw *model.Firewall) error {
	seenIPs := map[string]bool{}

	for name, iface := range fw.Interfaces {
		if iface.Zone != "" {
			if _, ok := fw.Zones[iface.Zone]; !ok {
				return &utils.APIError{Code: "ERR_NET_2003", Message: "Resource references unknown entity", Details: fmt.Sprintf("interface %s references unknown zone %s", name, iface.Zone)}
			}
		}

		if iface.IP == "" {
			continue
		}

		_, _, err := net.ParseCIDR(iface.IP)
		if err != nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: fmt.Sprintf("invalid IP on interface %s", name)}
		}

		if seenIPs[iface.IP] {
			return &utils.APIError{Code: "ERR_NET_2004", Message: "Duplicate IP address detected", Details: iface.IP}
		}
		seenIPs[iface.IP] = true
	}

	for vlanName, vlan := range fw.Vlans {
		parent := vlan.Parent

		parentIface, ok := fw.Interfaces[parent]
		if !ok {
			return &utils.APIError{Code: "ERR_NET_2005", Message: "VLAN references unknown parent interface", Details: fmt.Sprintf("vlan %s -> parent %s", vlanName, parent)}
		}

		if parentIface.IP != "" {
			return &utils.APIError{Code: "ERR_NET_2001", Message: "Interface acts as a parent for VLANs and cannot have a direct IP assigned", Details: parent}
		}

		if parentIface.State == "down" && vlan.State == "up" {
			return &utils.APIError{Code: "ERR_NET_2002", Message: "Parent interface cannot be down while child VLAN is up", Details: fmt.Sprintf("parent: %s, vlan: %s", parent, vlanName)}
		}
	}

	return nil
}
