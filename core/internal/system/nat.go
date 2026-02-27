package system

type NatInfo struct {
	ID 		     int    `json:"id"`
	SrcZone      string `json:"src-zone"`
	DstZone      string `json:"dst-zone"`
	OutInterface string `json:"out-interface"`
	Action       string `json:"action"`
	Description  string `json:"description"`
}

func GetNatRules() ([]NatInfo, error) {

fw := GetFirewall()
if fw == nil {
	return nil, ErrFirewallNotInitialized
}

	var out []NatInfo

	for _, nat := range fw.NATRules {
		out = append(out, NatInfo{
			ID: nat.ID,

			SrcZone: func() string {
				if nat.SrcZone != nil {
					return nat.SrcZone.Name
				}
				return ""
			}(),

			DstZone: func() string {
				if nat.DstZone != nil {
					return nat.DstZone.Name
				}
				return ""
			}(),

			OutInterface: nat.OutInterface,
			Action:       string(nat.Action),
			Description:  nat.Description,
		})
	}

	return out, nil
}	