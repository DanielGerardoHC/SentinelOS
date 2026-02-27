package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
)

func ValidateZones(fw *model.Firewall) error {

	seen := map[string]bool{}

	for _, z := range fw.Zones {

		if z.Name == "" {
			return fmt.Errorf("zone with empty name")
		}

		if seen[z.Name] {
			return fmt.Errorf("duplicate zone %s", z.Name)
		}
		seen[z.Name] = true

		for _, iface := range z.Interfaces {
			if _, ok := fw.Interfaces[iface]; !ok {
				return fmt.Errorf("zone %s references unknown interface %s", z.Name, iface)
			}
		}
	}

	return nil
}