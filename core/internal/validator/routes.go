package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
	"strings"
)

func ValidateRoutes(fw *model.Firewall) error {

	for _, r := range fw.Routes {

		dest := strings.TrimSpace(r.Destination)

		if dest == "default" {
			continue
		}

		if _, _, err := net.ParseCIDR(dest); err != nil {
			return fmt.Errorf("invalid destination in route %d", r.ID)
		}

		if _, _, err := net.ParseCIDR(r.Destination); err != nil {
			return fmt.Errorf("invalid destination in route %d", r.ID)
		}

		if net.ParseIP(r.Gateway) == nil {
			return fmt.Errorf("invalid gateway in route %d", r.ID)
		}

		if _, ok := fw.Interfaces[r.Interface]; !ok {
			return fmt.Errorf("route %d references unknown interface %s", r.ID, r.Interface)
		}
	}

	return nil
}
