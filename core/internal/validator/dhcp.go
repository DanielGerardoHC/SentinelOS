package validator

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
)

func ValidateDHCP(fw *model.Firewall) error {

	for _, d := range fw.DHCPConfigs {

		_, isIfc := fw.Interfaces[d.Interface]
		_, isVlan := fw.Vlans[d.Interface]

		if !isIfc && !isVlan {
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
