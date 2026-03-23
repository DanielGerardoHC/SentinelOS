package firewall

import (
	"fmt"
	"strings"

	"sentinelos/core/internal/model"
)

var ManagementPortMap = map[string]string{
	"PING":  "icmp type echo-request",
	"SSH":   "tcp dport 22",
	"HTTP":  "tcp dport 80",
	"HTTPS": "tcp dport 443",
	"API":   "tcp dport 8080", // Puerto del backend de SentinelOS
}

func GenerateRules(fw *model.Firewall) string {

	var sb strings.Builder

	sb.WriteString("flush ruleset\n\n")

	sb.WriteString("table inet filter {\n")

	sb.WriteString("  chain input {\n")
	sb.WriteString("    type filter hook input priority 0; policy drop;\n\n")
	sb.WriteString("    ct state established,related accept\n")
	sb.WriteString("    ct state invalid drop\n")
	sb.WriteString("    iifname \"lo\" accept\n\n")

	type IfaceMeta struct {
		Name       string
		Management []string
	}
	var allIfaces []IfaceMeta

	for _, i := range fw.Interfaces {
		if i.State == "up" {
			allIfaces = append(allIfaces, IfaceMeta{Name: i.Name, Management: i.Management})
		}
	}
	for _, v := range fw.Vlans {
		if v.State == "up" {
			allIfaces = append(allIfaces, IfaceMeta{Name: v.Name, Management: v.Management})
		}
	}

	for _, iface := range allIfaces {
		if len(iface.Management) > 0 {
			sb.WriteString(fmt.Sprintf("    # Mgmt para %s\n", iface.Name))
			for _, mgt := range iface.Management {
				if ruleSyntax, exists := ManagementPortMap[strings.ToUpper(mgt)]; exists {
					sb.WriteString(fmt.Sprintf("    iifname \"%s\" %s accept\n", iface.Name, ruleSyntax))
				}
			}
		}
		for _, dhcp := range fw.DHCPConfigs {
			if dhcp.Interface == iface.Name {
				sb.WriteString(fmt.Sprintf("    iifname \"%s\" udp dport { 67, 68 } accept comment \"DHCP Server\"\n", iface.Name))
				break
			}
		}
	}

	sb.WriteString("  }\n\n")
	sb.WriteString("  chain forward {\n")
	sb.WriteString("    type filter hook forward priority 0; policy drop;\n\n")
	sb.WriteString("    ct state established,related accept\n")
	sb.WriteString("    ct state invalid drop\n\n")
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

	if len(p.Services) == 0 {
		rule := fmt.Sprintf(
			"    %s %s %s\n",
			srcMatch,
			dstMatch,
			actionToNft(p.Action),
		)
		rule = strings.ReplaceAll(rule, "  ", " ")
		rules = append(rules, rule)
	} else {
		for _, svc := range p.Services {
			rule := fmt.Sprintf(
				"    %s %s %s dport %s %s\n",
				srcMatch,
				dstMatch,
				svc.Protocol,
				portsToString(svc.Ports),
				actionToNft(p.Action),
			)
			rule = strings.ReplaceAll(rule, "  ", " ")
			rules = append(rules, rule)
		}
	}

	return strings.Join(rules, "")
}

func zoneSrcMatch(z *model.Zone) string {
	if z == nil {
		return ""
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
		return ""
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
