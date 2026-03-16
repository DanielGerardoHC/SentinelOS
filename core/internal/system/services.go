package system

import "sentinelos/core/internal/model"

type Service struct {
	Name     string         `json:"name"`
	Protocol model.Protocol `json:"protocol"`
	Ports    []int          `json:"ports"`
}

func GetServices() ([]Service, error) {
	fw := GetFirewall()
	if fw == nil {
		return nil, ErrFirewallNotInitialized
	}

	out := make([]Service, 0, len(fw.Services))

	for _, service := range fw.Services {
		out = append(out, Service{
			Name:     service.Name,
			Protocol: service.Protocol,
			Ports:    service.Ports,
		})
	}

	return out, nil
}
