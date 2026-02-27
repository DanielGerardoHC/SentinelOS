package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
)


func ValidateDHCP(fw *model.Firewall) error {

	for _, d := range fw.DHCPConfigs {

		if _, ok := fw.Interfaces[d.Interface]; !ok {
			return fmt.Errorf("dhcp references unknown interface %s", d.Interface)
		}

		if net.ParseIP(d.StartIP) == nil {
			return fmt.Errorf("invalid start-ip on dhcp %s", d.Interface)
		}

		if net.ParseIP(d.EndIP) == nil {
			return fmt.Errorf("invalid end-ip on dhcp %s", d.Interface)
		}

		if net.ParseIP(d.Gateway) == nil {
			return fmt.Errorf("invalid gateway on dhcp %s", d.Interface)
		}
	}

	return nil
}