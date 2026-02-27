
package validator

import (
	"net"
	"fmt"
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

	return nil
}