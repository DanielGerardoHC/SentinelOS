package system

type Address struct {
	Name string   `json:"name"`
	Ips  []string `json:"ips"`
}

func GetAddresses() ([]Address, error) {
	fw := GetFirewall()
	if fw == nil {
		return nil, ErrFirewallNotInitialized
	}

	out := make([]Address, 0, len(fw.Addresses))

	for _, address := range fw.Addresses {
		var ipsStr []string

		for _, ipNet := range address.IPs {
			ipsStr = append(ipsStr, ipNet.String())
		}

		out = append(out, Address{
			Name: address.Name,
			Ips:  ipsStr,
		})
	}

	return out, nil
}
