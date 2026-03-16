package validator

import (
	"fmt"

	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateServices(fw *model.Firewall) error {
	seen := map[string]bool{}

	for _, s := range fw.Services {
		if seen[s.Name] {
			return &utils.APIError{Code: "ERR_NET_2006", Message: "Resource already exists", Details: fmt.Sprintf("duplicate service %s", s.Name)}
		}
		seen[s.Name] = true

		for _, p := range s.Ports {
			if p <= 0 || p > 65535 {
				return &utils.APIError{Code: "ERR_NET_1008", Message: "Invalid port number", Details: fmt.Sprintf("invalid port in service %s", s.Name)}
			}
		}
	}

	return nil
}
