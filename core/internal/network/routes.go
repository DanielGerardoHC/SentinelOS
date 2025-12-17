package network

import (
	"fmt"
	"sentinelos/core/internal/model"
	"strings"
)

func GenerateRoutesConfig(routes []*model.Route) string {

	var sb strings.Builder

	for _, route := range routes {
		if route.Destination == "0.0.0.0/0" {
			route.Destination = "default"
		}
		sb.WriteString(fmt.Sprintf(
			"ip route replace %s via %s dev %s metric %d\n",
			route.Destination,
			route.Gateway,
			route.Interface,
			route.Metric,
		))

	}
	return sb.String()
}
