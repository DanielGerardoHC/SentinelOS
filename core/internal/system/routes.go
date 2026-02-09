package system

type RouteInfo struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Interface   string `json:"interface"`
	Metric      int    `json:"metric"`
	Description string `json:"description"`
}

func GetRoutes() ([]RouteInfo, error) {

	if firewall == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []RouteInfo

	for _, route := range firewall.Routes {
		out = append(out, RouteInfo{
			Destination: route.Destination,
			Gateway: route.Gateway,
			Interface: route.Interface,
			Metric: route.Metric,
			Description: route.Description,
		})
	}

	return out, nil
}
