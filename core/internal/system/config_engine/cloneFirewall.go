package config_engine

import "sentinelos/core/internal/model"

func CloneFirewall(src *model.Firewall) (*model.Firewall, error) {

	if src == nil {
		return nil, nil
	}

	clone := &model.Firewall{
		Zones:       make(map[string]*model.Zone),
		Interfaces:  make(map[string]*model.Interface),
		Vlans:       make(map[string]*model.Vlan),
		Addresses:   make(map[string]*model.Address),
		Services:    make(map[string]*model.Service),
		Routes:      make([]*model.Route, len(src.Routes)),
		Policies:    make([]*model.Policy, len(src.Policies)),
		NATRules:    make([]*model.NATRule, len(src.NATRules)),
		DHCPConfigs: make([]*model.DHCP, len(src.DHCPConfigs)),
	}

	// Clonar maps
	for k, v := range src.Zones {
		zoneCopy := *v
		clone.Zones[k] = &zoneCopy
	}

	for k, v := range src.Interfaces {
		ifaceCopy := *v
		clone.Interfaces[k] = &ifaceCopy
	}

	for k, v := range src.Vlans {
		vlanCopy := *v
		clone.Vlans[k] = &vlanCopy
	}

	for k, v := range src.Addresses {
		addrCopy := *v
		clone.Addresses[k] = &addrCopy
	}

	for k, v := range src.Services {
		svcCopy := *v
		clone.Services[k] = &svcCopy
	}

	// Clonar slices
	copy(clone.Routes, src.Routes)
	copy(clone.Policies, src.Policies)
	copy(clone.NATRules, src.NATRules)
	copy(clone.DHCPConfigs, src.DHCPConfigs)

	return clone, nil
}