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

	sb.WriteString("  chain prerouting {\n")
	sb.WriteString("    type nat hook prerouting priority dstnat; policy accept;\n")

	for _, r := range fw.NATRules {
		if r.Type == model.TypeDNATIP || r.Type == model.TypeDNATPort {
			sb.WriteString(generateDNATRule(r))
		}
	}
	sb.WriteString("  }\n\n")

	sb.WriteString("  chain postrouting {\n")
	sb.WriteString("    type nat hook postrouting priority srcnat; policy accept;\n")

	for _, r := range fw.NATRules {
		if r.Type == model.TypeSNAT {
			sb.WriteString(generateSNATRule(r))
		}
	}

	sb.WriteString("  }\n")
	sb.WriteString("}\n")

	return sb.String()
}


func formatSet(items []string, isInterface bool) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		if isInterface {
			return fmt.Sprintf("\"%s\"", items[0])
		}
		return items[0]
	}

	var formattedItems []string
	for _, item := range items {
		if isInterface {
			formattedItems = append(formattedItems, fmt.Sprintf("\"%s\"", item))
		} else {
			formattedItems = append(formattedItems, item)
		}
	}
	return fmt.Sprintf("{ %s }", strings.Join(formattedItems, ", "))
}

func generateSNATRule(r *model.NATRule) string {
	var matchParts []string

	if r.SrcZone != nil {
		if len(r.SrcZone.Interfaces) > 0 {
			matchParts = append(matchParts, fmt.Sprintf("iifname %s", formatSet(r.SrcZone.Interfaces, true)))
		}
		if len(r.SrcZone.Networks) > 0 {
			matchParts = append(matchParts, fmt.Sprintf("ip saddr %s", formatSet(r.SrcZone.Networks, false)))
		}
	}

	if r.SrcAddress != "" && r.SrcAddress != "any" {
		matchParts = append(matchParts, fmt.Sprintf("ip saddr %s", r.SrcAddress)) 
	}

	if r.DstZone != nil {
		if r.OutInterface == "" || r.OutInterface == "any" {
			if len(r.DstZone.Interfaces) > 0 {
				matchParts = append(matchParts, fmt.Sprintf("oifname %s", formatSet(r.DstZone.Interfaces, true)))
			}
		}
		if len(r.DstZone.Networks) > 0 {
			matchParts = append(matchParts, fmt.Sprintf("ip daddr %s", formatSet(r.DstZone.Networks, false)))
		}
	}

	if r.DstAddress != "" && r.DstAddress != "any" {
		matchParts = append(matchParts, fmt.Sprintf("ip daddr %s", r.DstAddress))
	}

	if r.OutInterface != "" && r.OutInterface != "any" {
		matchParts = append(matchParts, fmt.Sprintf("oifname \"%s\"", r.OutInterface))
	}

	match := strings.Join(matchParts, " ")
	if match != "" {
		match += " "
	}

	if r.TranslatedIP == "" {
		return fmt.Sprintf("    %smasquerade comment \"%s\"\n", match, r.Description)
	}

	return fmt.Sprintf("    %ssnat to %s comment \"%s\"\n", match, r.TranslatedIP, r.Description)
}

func generateDNATRule(r *model.NATRule) string {
	var matchParts []string

	if r.DstAddress != "" && r.DstAddress != "any" {
		matchParts = append(matchParts, fmt.Sprintf("ip daddr %s", r.DstAddress))
	}

	if r.Type == model.TypeDNATIP {
		match := strings.Join(matchParts, " ")
		if match != "" {
			match += " "
		}
		return fmt.Sprintf("    %sdnat to %s comment \"%s\"\n", match, r.TranslatedIP, r.Description)
		
	} else if r.Type == model.TypeDNATPort {

		if r.Service != "" && r.Service != "any" {
			matchParts = append(matchParts, fmt.Sprintf("tcp dport %s", r.Service))
		}

		match := strings.Join(matchParts, " ")
		if match != "" {
			match += " "
		}

		portPart := ""
		if r.TranslatedPort != "" {
			portPart = ":" + r.TranslatedPort
		}

		return fmt.Sprintf("    %sdnat to %s%s comment \"%s\"\n", match, r.TranslatedIP, portPart, r.Description)
	}

	return ""
}