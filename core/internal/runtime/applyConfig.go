package runtime

import (
	"sentinelos/core/internal/firewall"
	"sentinelos/core/internal/network"
	"sentinelos/core/internal/model"
)

func ApplyFullRuntime(fw *model.Firewall) error {

	interfaces := network.GenerateInterfacesConfig(fw.Interfaces)
	if err := firewall.ApplyInterfacesConfig(interfaces); err != nil {
		return err
	}

	vlans := network.GenerateVlansConfig(fw.Vlans)
	if err := firewall.ApplyVlansConfig(vlans); err != nil {
		return err
	}

	routes := network.GenerateRoutesConfig(fw.Routes)
	if err := firewall.ApplyRoutes(routes); err != nil {
		return err
	}

	var ifaceNames []string
	for name := range fw.Interfaces {
		ifaceNames = append(ifaceNames, name)
	}
if err := ApplySysctl(ifaceNames); err != nil {
		return err
	}

	rules := firewall.GenerateRules(fw)
	rules += firewall.GenerateNATRules(fw)

	if err := firewall.ApplyRules(rules); err != nil {
		return err
	}

	dnsmasq := network.GenerateDnsmasqConfig(fw.DHCPConfigs)
	if err := firewall.ApplyDHCP(dnsmasq); err != nil {
		return err
	}

	return nil
}