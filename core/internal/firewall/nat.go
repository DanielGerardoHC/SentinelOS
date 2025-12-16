package firewall

import (
	"fmt"
	"strings"

	"sentinelos/core/internal/model"
)

func GenerateNATRules(fw *model.Firewall) string {

	if len(fw.NATRules) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("table ip nat {\n")
	sb.WriteString("  chain postrouting {\n")
	sb.WriteString("    type nat hook postrouting priority 100;\n")

	for _, r := range fw.NATRules {
		sb.WriteString(generateSNATRule(r))
	}

	sb.WriteString("  }\n")
	sb.WriteString("}\n")

	return sb.String()
}

func generateSNATRule(r *model.NATRule) string {

	if r.Action != model.Masquerade {
		return ""
	}

	return fmt.Sprintf(
		"    oifname \"%s\" masquerade\n",
		r.OutInterface,
	)
}
