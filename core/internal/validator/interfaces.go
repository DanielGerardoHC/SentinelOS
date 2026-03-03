package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
)

func ValidateInterfaces(fw *model.Firewall) error {

	seenIPs := map[string]bool{}

	for name, iface := range fw.Interfaces {

		if iface.Zone != "" {
			if _, ok := fw.Zones[iface.Zone]; !ok {
				return fmt.Errorf("interface %s references unknown zone %s", name, iface.Zone)
			}
		}

		if iface.IP == "" {
			continue
		}

		_, _, err := net.ParseCIDR(iface.IP)
		if err != nil {
			return fmt.Errorf("invalid IP on interface %s", name)
		}

		if seenIPs[iface.IP] {
			return fmt.Errorf("duplicate IP detected: %s", iface.IP)
		}
		seenIPs[iface.IP] = true
	}
	// Regla enterprise:
	// Si una interfaz tiene VLANs asociadas, no puede tener IP
	for _, vlan := range fw.Vlans {
		parent := vlan.Parent

		if parentIface, ok := fw.Interfaces[parent]; ok {
			if parentIface.IP != "" {
				return fmt.Errorf(
					"interface %s has VLANs configured and cannot have IP assigned",
					parent,
				)
			}
		}
	}

	return nil
}
