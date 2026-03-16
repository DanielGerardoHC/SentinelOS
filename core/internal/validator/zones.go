package validator

import (
	"fmt"

	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateZones(fw *model.Firewall) error {
	seen := map[string]bool{}

	for _, z := range fw.Zones {
		if z.Name == "" {
			return &utils.APIError{Code: "ERR_NET_1002", Message: "Missing required field", Details: "zone with empty name"}
		}

		if seen[z.Name] {
			return &utils.APIError{Code: "ERR_NET_2006", Message: "Resource already exists", Details: fmt.Sprintf("duplicate zone %s", z.Name)}
		}
		seen[z.Name] = true

		for _, iface := range z.Interfaces {
			if _, ok := fw.Interfaces[iface]; !ok {
				return &utils.APIError{Code: "ERR_NET_2003", Message: "Resource references unknown entity", Details: fmt.Sprintf("zone %s references unknown interface %s", z.Name, iface)}
			}
		}
	}

	return nil
}
