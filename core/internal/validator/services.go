package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
)
func ValidateServices(fw *model.Firewall) error {

	seen := map[string]bool{}

	for _, s := range fw.Services {

		if seen[s.Name] {
			return fmt.Errorf("duplicate service %s", s.Name)
		}
		seen[s.Name] = true

		for _, p := range s.Ports {
			if p <= 0 || p > 65535 {
				return fmt.Errorf("invalid port in service %s", s.Name)
			}
		}
	}

	return nil
}