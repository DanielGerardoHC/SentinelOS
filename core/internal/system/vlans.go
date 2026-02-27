package system

type vlanInfo struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Parent    string   `json:"parent"`
	IP        string   `json:"ip"`
	Zone      string   `json:"zone"`
	State     string   `json:"state"`
	Management []string `json:"management"`
}

func GetVlans() ([]vlanInfo, error) {

fw := GetFirewall()
if fw == nil {
	return nil, ErrFirewallNotInitialized
}

	var out []vlanInfo

	for _, vlan := range fw.Vlans {
		out = append(out, vlanInfo{
			ID:        vlan.ID,
			Name:      vlan.Name,
			Parent:    vlan.Parent,
			IP:        vlan.IP,
			Zone:      vlan.Zone,
			State:     vlan.State,
			Management: vlan.Management,
		})
	}

	return out, nil
}