package network

import (
	"fmt"
	"sentinelos/core/internal/model"
	"strings"
)

func GenerateVlansConfig(vlans map[string]*model.Vlan) string {

	var sb strings.Builder
	for _, vlan := range vlans {
		sb.WriteString(fmt.Sprintf(
			"ip link replace link %s name %s type vlan id %d\n",
			vlan.Parent,
			vlan.Name,
			vlan.ID,
		))
		sb.WriteString(fmt.Sprintf(
			"ip addr replace %s dev %s\n",
			vlan.IP,
			vlan.Name,
		))

		if vlan.State == "up" {
			sb.WriteString(fmt.Sprintf(
				"ip link set %s up\n",
				vlan.Name,
			))
		} else if vlan.State == "down" {
			sb.WriteString(fmt.Sprintf(
				"ip link set %s down\n",
				vlan.Name,
			))
		}
	}
	return sb.String()
}
