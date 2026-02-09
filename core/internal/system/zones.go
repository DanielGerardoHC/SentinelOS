package system

type ZoneInfo struct {
	Name        string   `json:"name"`
	Networks    []string   `json:"networks"`
	Interfaces  []string `json:"interfaces"`
	Type	    string   `json:"type"`
}

func GetZones() ([]ZoneInfo, error) {

	if firewall == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []ZoneInfo

	for _, zone := range firewall.Zones {
		out = append(out, ZoneInfo{
			Name:        zone.Name,
			Networks:    zone.Networks,
			Interfaces:  zone.Interfaces,
			Type:        string(zone.Type),
		})
	}

	return out, nil
}		