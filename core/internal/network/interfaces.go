package network

import (
	"fmt"
	"sentinelos/core/internal/model"
	"strings"
)

func GenerateInterfacesConfig(interfaces map[string]*model.Interface) string {

	var sb strings.Builder
	for _, iface := range interfaces {
		sb.WriteString(fmt.Sprintf(
			"ip addr flush dev %s\n",
			iface.Name,
		))
		sb.WriteString(fmt.Sprintf(
			"ip addr add %s dev %s\n",
			iface.IP,
			iface.Name,
		))

		if iface.State == "up" {
			sb.WriteString(fmt.Sprintf(
				"ip link set %s up\n",
				iface.Name,
			))
		} else if iface.State == "down" {
			sb.WriteString(fmt.Sprintf(
				"ip link set %s down\n",
				iface.Name,
			))
		}
	}
	return sb.String()
}
