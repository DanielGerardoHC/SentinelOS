package config_engine

import (
	"os"

	"gopkg.in/yaml.v3"
	"sentinelos/core/internal/config"
	"sentinelos/core/internal/model"
)

func SaveYAML(fw *model.Firewall) error {

	var raw config.RawConfig

	// Zones
	for _, z := range fw.Zones {
		raw.Zones = append(raw.Zones, struct {
			Name       string   `yaml:"name"`
			Type       string   `yaml:"type"`
			Interfaces []string `yaml:"interfaces"`
			Networks   []string `yaml:"networks"`
		}{
			Name:       z.Name,
			Type:       string(z.Type),
			Interfaces: z.Interfaces,
			Networks:   z.Networks,
		})
	}

	// Interfaces
	for _, i := range fw.Interfaces {
		raw.Interfaces = append(raw.Interfaces, struct {
			Name       string   `yaml:"name"`
			IP         string   `yaml:"ip"`
			Zone       string   `yaml:"zone"`
			State      string   `yaml:"state"`
			Management []string `yaml:"management"`
		}{
			Name:       i.Name,
			IP:         i.IP,
			Zone:       i.Zone,
			State:      i.State,
			Management: i.Management,
		})
	}

	// Vlans
	for _, v := range fw.Vlans {
		raw.Vlans = append(raw.Vlans, struct {
			Name       string   `yaml:"name"`
			Parent     string   `yaml:"parent"`
			ID         int      `yaml:"id"`
			IP         string   `yaml:"ip"`
			Zone       string   `yaml:"zone"`
			State      string   `yaml:"state"`
			Management []string `yaml:"management"`
		}{
			Name:       v.Name,
			Parent:     v.Parent,
			ID:         v.ID,
			IP:         v.IP,
			Zone:       v.Zone,
			State:      v.State,
			Management: v.Management,
		})
	}

	// Addresses
	for _, a := range fw.Addresses {
		var ips []string
		for _, ip := range a.IPs {
			ips = append(ips, ip.String())
		}

		raw.Addresses = append(raw.Addresses, struct {
			Name string   `yaml:"name"`
			IPs  []string `yaml:"ips"`
		}{
			Name: a.Name,
			IPs:  ips,
		})
	}

	// Services
	for _, s := range fw.Services {
		raw.Services = append(raw.Services, struct {
			Name     string `yaml:"name"`
			Protocol string `yaml:"protocol"`
			Ports    []int  `yaml:"ports"`
		}{
			Name:     s.Name,
			Protocol: string(s.Protocol),
			Ports:    s.Ports,
		})
	}

	// Routes
	for _, r := range fw.Routes {
		raw.Route = append(raw.Route, struct {
			ID          int    `yaml:"id"`
			Destination string `yaml:"destination"`
			Gateway     string `yaml:"gateway"`
			Interface   string `yaml:"interface"`
			Metric      int    `yaml:"metric"`
			Description string `yaml:"description"`
		}{
			ID:          r.ID,
			Destination: r.Destination,
			Gateway:     r.Gateway,
			Interface:   r.Interface,
			Metric:      r.Metric,
			Description: r.Description,
		})
	}

	// Policies
	for _, p := range fw.Policies {

		var services []string
		for _, s := range p.Services {
			services = append(services, s.Name)
		}

		raw.Policies = append(raw.Policies, struct {
			ID       int      `yaml:"id"`
			SrcZone  string   `yaml:"src-zone"`
			DstZone  string   `yaml:"dst-zone"`
			SrcAddr  string   `yaml:"src-addr"`
			DstAddr  string   `yaml:"dst-addr"`
			Services []string `yaml:"services"`
			Action   string   `yaml:"action"`
			Log      bool     `yaml:"log"`
		}{
			ID:       p.ID,
			SrcZone:  zoneNameOrAny(p.SrcZone),
			DstZone:  zoneNameOrAny(p.DstZone),
			SrcAddr:  addrNameOrAny(p.SrcAddr),
			DstAddr:  addrNameOrAny(p.DstAddr),
			Services: services,
			Action:   string(p.Action),
			Log:      p.Log,
		})
	}

	// NAT
	for _, n := range fw.NATRules {
		raw.NATRules = append(raw.NATRules, struct {
			ID           int    `yaml:"id"`
			Type         string `yaml:"type"`
			SrcZone      string `yaml:"src-zone"`
			DstZone      string `yaml:"dst-zone"`
			Action       string `yaml:"action"`
			OutInterface string `yaml:"outInterface"`
			Description  string `yaml:"description"`
		}{
			ID:           n.ID,
			Type:         n.Type,
			SrcZone:      zoneNameOrAny(n.SrcZone),
			DstZone:      zoneNameOrAny(n.DstZone),
			Action:       string(n.Action),
			OutInterface: n.OutInterface,
			Description:  n.Description,
		})
	}

	// DHCP
	for _, d := range fw.DHCPConfigs {
		raw.DHCP = append(raw.DHCP, struct {
			Interface    string   `yaml:"interface"`
			StartIP      string   `yaml:"start-ip"`
			EndIP        string   `yaml:"end-ip"`
			SubnetMask   string   `yaml:"subnet-mask"`
			Gateway      string   `yaml:"gateway"`
			DNS          []string `yaml:"dns"`
			LeaseTimeMin int      `yaml:"lease-time-min"`
		}{
			Interface:    d.Interface,
			StartIP:      d.StartIP,
			EndIP:        d.EndIP,
			SubnetMask:   d.SubnetMask,
			Gateway:      d.Gateway,
			DNS:          d.DNS,
			LeaseTimeMin: d.LeaseTimeMin,
		})
	}

	data, err := yaml.Marshal(&raw)
	if err != nil {
		return err
	}

	return os.WriteFile(
		"/srv/sentinelos/core/internal/config/config.yml",
		data,
		0644,
	)
}

func zoneName(z *model.Zone) string {
	if z == nil {
		return ""
	}
	return z.Name
}

func addrName(a *model.Address) string {
	if a == nil {
		return ""
	}
	return a.Name
}

func zoneNameOrAny(z *model.Zone) string {
	if z == nil {
		return "any"
	}
	return z.Name
}

func addrNameOrAny(a *model.Address) string {
	if a == nil {
		return "any"
	}
	return a.Name
}
