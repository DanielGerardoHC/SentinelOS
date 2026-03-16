package validator

import (
	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateAddresses(fw *model.Firewall) error {
	for name, addr := range fw.Addresses {
		for _, ipnet := range addr.IPs {
			if ipnet.IP == nil || ipnet.Mask == nil {
				return &utils.APIError{Code: "ERR_NET_1001", Message: "Invalid IP or CIDR format", Details: "invalid network in address " + name}
			}

			ones, bits := ipnet.Mask.Size()
			if ones < 0 || bits != 32 {
				return &utils.APIError{Code: "ERR_NET_1007", Message: "Invalid subnet mask", Details: "invalid mask in address " + name}
			}
		}
	}

	return nil
}
