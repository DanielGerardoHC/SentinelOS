package validator

import (
	"fmt"
	"net"
	"strings"

	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateRoutes(fw *model.Firewall) error {
	for _, r := range fw.Routes {
		dest := strings.TrimSpace(r.Destination)

		if dest == "default" {
			continue
		}

		if _, _, err := net.ParseCIDR(dest); err != nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: fmt.Sprintf("invalid destination in route %d", r.ID)}
		}

		if _, _, err := net.ParseCIDR(r.Destination); err != nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: fmt.Sprintf("invalid destination in route %d", r.ID)}
		}

		if net.ParseIP(r.Gateway) == nil {
			return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: fmt.Sprintf("invalid gateway in route %d", r.ID)}
		}

		if _, ok := fw.Interfaces[r.Interface]; !ok {
			return &utils.APIError{Code: "ERR_NET_2003", Message: "Resource references unknown entity", Details: fmt.Sprintf("route %d references unknown interface %s", r.ID, r.Interface)}
		}
	}

	return nil
}
