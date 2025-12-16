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
		Addresses:  make(map[string]*model.Address),
		Services:   make(map[string]*model.Service),
		Policies:   []*model.Policy{},
	}

	// pasos:
	// 1. zonas
	for _, z := range raw.Zones {
		if z.Name == "" {
			return nil, fmt.Errorf("zona sin nombre")
		}

		zone := &model.Zone{
			Name: z.Name,
			Type: model.ZoneType(z.Type),
		}

		fw.Zones[z.Name] = zone
	}

	// 2. interfaces
	for _, i := range raw.Interfaces {

		zone, ok := fw.Zones[i.Zone]
		if !ok {
			return nil, fmt.Errorf("la interfaz %s referencia zona inexistente %s", i.Name, i.Zone)
		}

		_, ipnet, err := net.ParseCIDR(i.IP)
		if err != nil {
			return nil, fmt.Errorf("IP inválida en interfaz %s", i.Name)
		}

		iface := &model.Interface{
			Name:       i.Name,
			Zone:       zone,
			IP:         ipnet,
			Management: i.Management,
		}

		fw.Interfaces[i.Name] = iface
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

	// 5. policies

	for _, p := range raw.Policies {

		srcZone := fw.Zones[p.SrcZone]
		dstZone := fw.Zones[p.DstZone]
		srcAddr := fw.Addresses[p.SrcAddr]
		dstAddr := fw.Addresses[p.DstAddr]

		if srcZone == nil || dstZone == nil {
			return nil, fmt.Errorf("policy %d referencia zonas inválidas", p.ID)
		}

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
	return fw, nil
}
