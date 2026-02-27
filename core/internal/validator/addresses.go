package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
)

func ValidateAddresses(fw *model.Firewall) error {

	for name, addr := range fw.Addresses {

		for _, ipnet := range addr.IPs {

			//validar IP
			if ipnet.IP == nil || ipnet.Mask == nil {
				return fmt.Errorf("invalid network in address %s", name)
			}

			//verificar mascara valida
			ones, bits := ipnet.Mask.Size()
			if ones < 0 || bits != 32 {
				return fmt.Errorf("invalid mask in address %s", name)
			}
		}
	}

	return nil
}