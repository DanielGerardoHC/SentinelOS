package system

type RouteInfo struct {
	ID          int    `json:"id"`
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Interface   string `json:"interface"`
	Metric      int    `json:"metric"`
	Description string `json:"description"`
}

func GetRoutes() ([]RouteInfo, error) {

	fw := GetFirewall()
	if fw == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []RouteInfo

	for _, route := range fw.Routes {
		out = append(out, RouteInfo{
			ID:          route.ID,
			Destination: route.Destination,
			Gateway:     route.Gateway,
			Interface:   route.Interface,
			Metric:      route.Metric,
			Description: route.Description,
		})
	}

	return out, nil
}
