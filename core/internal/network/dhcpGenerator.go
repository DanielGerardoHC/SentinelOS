package network

import (
	"fmt"
	"strings"

	"sentinelos/core/internal/model"
)

func GenerateDnsmasqConfig(dhcps []*model.DHCP) string {

	var sb strings.Builder

	sb.WriteString("log-dhcp\n\n")

	for _, d := range dhcps {

		sb.WriteString(fmt.Sprintf(
			"interface=%s\n",
			d.Interface,
		))

		sb.WriteString(fmt.Sprintf(
			"dhcp-range=%s,%s,%s,%d\n",
			d.Interface,
			d.StartIP,
			d.EndIP,
			d.LeaseTimeMin,
		))

		sb.WriteString(fmt.Sprintf(
			"dhcp-option=3,%s\n",
			d.Gateway,
		))

		if len(d.DNS) > 0 {
			sb.WriteString(fmt.Sprintf(
				"dhcp-option=6,%s\n",
				strings.Join(d.DNS, ","),
			))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
