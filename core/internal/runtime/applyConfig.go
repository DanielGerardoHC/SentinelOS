package runtime

import (
	"sentinelos/core/internal/firewall"
	"sentinelos/core/internal/network"
	"sentinelos/core/internal/model"
)

func ApplyFullRuntime(fw *model.Firewall) error {

	// Interfaces
	interfaces := network.GenerateInterfacesConfig(fw.Interfaces)
	if err := firewall.ApplyInterfacesConfig(interfaces); err != nil {
		return err
	}

	// VLANs
	vlans := network.GenerateVlansConfig(fw.Vlans)
	if err := firewall.ApplyVlansConfig(vlans); err != nil {
		return err
	}

	// Routes
	routes := network.GenerateRoutesConfig(fw.Routes)
	if err := firewall.ApplyRoutes(routes); err != nil {
		return err
	}

	// Sysctl
	var ifaceNames []string
	for name := range fw.Interfaces {
		ifaceNames = append(ifaceNames, name)
	}
if err := ApplySysctl(ifaceNames); err != nil {
		return err
	}

	// Firewall rules
	rules := firewall.GenerateRules(fw)
	rules += firewall.GenerateNATRules(fw)

	if err := firewall.ApplyRules(rules); err != nil {
		return err
	}

	// DHCP
	dnsmasq := network.GenerateDnsmasqConfig(fw.DHCPConfigs)
	if err := firewall.ApplyDHCP(dnsmasq); err != nil {
		return err
	}

	return nil
}