package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
)
func ValidateVlans(fw *model.Firewall) error {

	for _, v := range fw.Vlans {

		if _, ok := fw.Interfaces[v.Parent]; !ok {
			return fmt.Errorf("vlan %s parent interface %s not found", v.Name, v.Parent)
		}

		if v.Zone != "" {
			if _, ok := fw.Zones[v.Zone]; !ok {
				return fmt.Errorf("vlan %s references unknown zone %s", v.Name, v.Zone)
			}
		}

		if v.IP != "" {
			if _, _, err := net.ParseCIDR(v.IP); err != nil {
				return fmt.Errorf("invalid IP on vlan %s", v.Name)
			}
		}
	}

	return nil
}