package firewall

import (
	"fmt"
	"strings"

	"sentinelos/core/internal/model"
)

func GenerateRules(fw *model.Firewall) string {

	var sb strings.Builder

	sb.WriteString("flush ruleset\n\n")

	sb.WriteString("table inet filter {\n")

	sb.WriteString("  chain input {\n")
	sb.WriteString("    type filter hook input priority 0;\n")
	sb.WriteString("    udp dport 67 accept\n")
	sb.WriteString("    udp dport 68 accept\n")
	sb.WriteString("    policy drop;\n")

	sb.WriteString("    ct state established,related accept\n")
	sb.WriteString("    iifname \"lo\" accept\n")
	sb.WriteString("    tcp dport 22 accept\n") // SSH
	sb.WriteString("    ip protocol icmp accept\n")

	sb.WriteString("  }\n\n")

	sb.WriteString("  chain forward {\n")
	sb.WriteString("    type filter hook forward priority 0;\n")
	sb.WriteString("    udp dport 67 accept\n")
	sb.WriteString("    udp dport 68 accept\n")
	sb.WriteString("    policy drop;\n\n")
	sb.WriteString("    ct state established,related accept\n")

	for _, p := range fw.Policies {
		sb.WriteString(generatePolicyRule(p, fw))
	}

	sb.WriteString("  }\n")
	sb.WriteString("}\n")

	return sb.String()
}

func generatePolicyRule(p *model.Policy, fw *model.Firewall) string {

	srcMatch := zoneSrcMatch(p.SrcZone)
	dstMatch := zoneDstMatch(p.DstZone)

	var rules []string

	for _, svc := range p.Services {

		rule := fmt.Sprintf(
			"    %s %s %s dport %s %s\n",
			srcMatch,
			dstMatch,
			svc.Protocol,
			portsToString(svc.Ports),
			actionToNft(p.Action),
		)

		rules = append(rules, rule)
	}

	return strings.Join(rules, "")
}

func zoneSrcMatch(z *model.Zone) string {
	if z == nil {
		return "" // ANY
	}

	if len(z.Networks) > 0 {
		return fmt.Sprintf("ip saddr %s", cidrSet(z.Networks))
	}

	if len(z.Interfaces) > 0 {
		return fmt.Sprintf("iifname %s", ifaceSet(z.Interfaces))
	}

	return ""
}

func zoneDstMatch(z *model.Zone) string {
	if z == nil {
		return "" // ANY
	}

	if len(z.Networks) > 0 {
		return fmt.Sprintf("ip daddr %s", cidrSet(z.Networks))
	}

	if len(z.Interfaces) > 0 {
		return fmt.Sprintf("oifname %s", ifaceSet(z.Interfaces))
	}

	return ""
}

func portsToString(ports []int) string {
	if len(ports) == 1 {
		return fmt.Sprintf("%d", ports[0])
	}

	var s []string
	for _, p := range ports {
		s = append(s, fmt.Sprintf("%d", p))
	}
	return "{ " + strings.Join(s, ", ") + " }"
}

func actionToNft(a model.Action) string {
	switch a {
	case model.Allow:
		return "accept"
	case model.Deny:
		return "drop"
	default:
		return "drop"
	}
}

func cidrSet(nets []string) string {
	if len(nets) == 1 {
		return nets[0]
	}
	return "{ " + strings.Join(nets, ", ") + " }"
}

func ifaceSet(ifaces []string) string {
	if len(ifaces) == 1 {
		return fmt.Sprintf("\"%s\"", ifaces[0])
	}

	var quoted []string
	for _, i := range ifaces {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", i))
	}

	return "{ " + strings.Join(quoted, ", ") + " }"
}
