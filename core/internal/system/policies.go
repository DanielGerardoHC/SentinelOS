package system

type PolicyInfo struct {
	ID      int      `json:"id"`
	Service []string `json:"services"`
	SrcZone string   `json:"src-zone"`
	DstZone string   `json:"dst-zone"`
	SrcAddr string   `json:"src-addr"`
	DstAddr string   `json:"dst-addr"`
	Action  string   `json:"action"`
	Log     bool     `json:"log"`
}

func GetPolicies() ([]PolicyInfo, error) {

	if firewall == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []PolicyInfo

	for _, policy := range firewall.Policies {

		// Convertir services 
		var services []string
		for _, s := range policy.Services {
			if s != nil {
				services = append(services, s.Name)
			}
		}

		// Manejo seguro de punteros
		srcZone := ""
		if policy.SrcZone != nil {
			srcZone = policy.SrcZone.Name
		}

		dstZone := ""
		if policy.DstZone != nil {
			dstZone = policy.DstZone.Name
		}

		srcAddr := ""
		if policy.SrcAddr != nil {
			srcAddr = policy.SrcAddr.Name
		}

		dstAddr := ""
		if policy.DstAddr != nil {
			dstAddr = policy.DstAddr.Name
		}

		out = append(out, PolicyInfo{
			ID:      policy.ID,
			Service: services,
			SrcZone: srcZone,
			DstZone: dstZone,
			SrcAddr: srcAddr,
			DstAddr: dstAddr,
			Action:  string(policy.Action),
			Log:     policy.Log,
		})
	}

	return out, nil
}
