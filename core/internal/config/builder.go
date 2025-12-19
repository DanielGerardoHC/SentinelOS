package config

import (
	"fmt"
	"net"
	"sentinelos/core/internal/model"
)

func BuildFirewall(raw *RawConfig) (*model.Firewall, error) {
	fw := &model.Firewall{
		Zones:      make(map[string]*model.Zone),
		Interfaces: make(map[string]*model.Interface),
		Vlans:      make(map[string]*model.Vlan),
		Addresses:  make(map[string]*model.Address),
		Services:   make(map[string]*model.Service),
		Routes:     []*model.Route{},
		Policies:   []*model.Policy{},
		NATRules:   []*model.NATRule{},
	}

	// pasos:
	// 1. zonas
	for _, z := range raw.Zones {
		if z.Name == "" {
			return nil, fmt.Errorf("zona sin nombre")
		}

		zone := &model.Zone{
			Name:       z.Name,
			Type:       model.ZoneType(z.Type),
			Interfaces: z.Interfaces,
			Networks:   z.Networks,
		}

		fw.Zones[z.Name] = zone
	}
	// 2. interfaces
	for _, i := range raw.Interfaces {

		iface := &model.Interface{
			Name:       i.Name,
			IP:         i.IP,
			Zone:       i.Zone,
			State:      i.State,
			Management: i.Management,
		}

		fw.Interfaces[i.Name] = iface
	}

	// 2.5 vlans
	for _, v := range raw.Vlans {

		vlan := &model.Vlan{
			Name:       v.Name,
			Parent:     v.Parent,
			ID:         v.ID,
			IP:         v.IP,
			Zone:       v.Zone,
			State:      v.State,
			Management: v.Management,
		}

		fw.Vlans[v.Name] = vlan
	}
	// 3. addresses

	for _, a := range raw.Addresses {

		addr := &model.Address{
			Name: a.Name,
		}

		for _, ip := range a.IPs {
			_, ipnet, err := net.ParseCIDR(ip)
			if err != nil {
				return nil, fmt.Errorf("IP inválida en address %s", a.Name)
			}
			addr.IPs = append(addr.IPs, *ipnet)
		}

		fw.Addresses[a.Name] = addr
	}
	// 4. services

	for _, s := range raw.Services {

		svc := &model.Service{
			Name:     s.Name,
			Protocol: model.Protocol(s.Protocol),
			Ports:    s.Ports,
		}

		fw.Services[s.Name] = svc
	}

	// Routes

	for _, r := range raw.Route {
		route := &model.Route{
			ID:          r.ID,
			Destination: r.Destination,
			Gateway:     r.Gateway,
			Interface:   r.Interface,
			Metric:      r.Metric,
			Description: r.Description,
		}
		fw.Routes = append(fw.Routes, route)
	}

	// 5. policies

	for _, p := range raw.Policies {

		srcZone := fw.Zones[p.SrcZone]
		dstZone := fw.Zones[p.DstZone]
		srcAddr := fw.Addresses[p.SrcAddr]
		dstAddr := fw.Addresses[p.DstAddr]

		/*	if srcZone == nil || dstZone == nil {
				return nil, fmt.Errorf("policy %d referencia zonas inválidas", p.ID)
			}
		*/
		var services []*model.Service
		for _, s := range p.Services {
			svc := fw.Services[s]
			if svc == nil {
				return nil, fmt.Errorf("policy %d referencia servicio inexistente %s", p.ID, s)
			}
			services = append(services, svc)
		}

		policy := &model.Policy{
			ID:       p.ID,
			SrcZone:  srcZone,
			DstZone:  dstZone,
			SrcAddr:  srcAddr,
			DstAddr:  dstAddr,
			Services: services,
			Action:   model.Action(p.Action),
			Log:      p.Log,
		}

		fw.Policies = append(fw.Policies, policy)
	}

	// Nat rules
	for _, n := range raw.NATRules {
		srcZone := fw.Zones[n.SrcZone]
		dstZone := fw.Zones[n.DstZone]

		if srcZone == nil || dstZone == nil {
			return nil, fmt.Errorf("NAT rule %d referencia zonas inválidas", n.ID)
		}

		natRule := &model.NATRule{
			ID:           n.ID,
			Type:         n.Type,
			SrcZone:      srcZone,
			DstZone:      dstZone,
			Action:       model.NATAction(n.Action),
			OutInterface: n.OutInterface,
			Description:  n.Description,
		}

		fw.NATRules = append(fw.NATRules, natRule)
	}

	// dhcp rules
	for _, d := range raw.DHCP {
		dhcp := &model.DHCP{
			Interface:    d.Interface,
			StartIP:      d.StartIP,
			EndIP:        d.EndIP,
			SubnetMask:   d.SubnetMask,
			Gateway:      d.Gateway,
			DNS:          d.DNS,
			LeaseTimeMin: d.LeaseTimeMin,
		}

		fw.DHCPConfigs = append(fw.DHCPConfigs, dhcp)
	}

	return fw, nil

}
