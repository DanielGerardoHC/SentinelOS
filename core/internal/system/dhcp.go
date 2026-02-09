package system

type DhcpInfo struct {
	Interfaces   string   `json:"interface"`
	StartIP      string   `json:"start_ip"`
	EndIP        string   `json:"end_ip"`
	Gateway      string   `json:"gateway"`
	DNS          []string `json:"dns"`
	LeaseTimeMin int      `json:"lease_time_min"`
}

func GetDhcpInfo() ([]DhcpInfo, error) {

	if firewall == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []DhcpInfo

	for _, dhcp := range firewall.DHCPConfigs {
		out = append(out, DhcpInfo{
			Interfaces:   dhcp.Interface,
			StartIP:      dhcp.StartIP,
			EndIP:        dhcp.EndIP,
			Gateway:      dhcp.Gateway,
			DNS:          dhcp.DNS,
			LeaseTimeMin: dhcp.LeaseTimeMin,
		})
	}

	return out, nil
}