package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
)

func ValidateRoutes(fw *model.Firewall) error {

	for _, r := range fw.Routes {

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