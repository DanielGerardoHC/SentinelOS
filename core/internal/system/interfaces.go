package system

import (
	"errors"
)

var ErrFirewallNotInitialized = errors.New("firewall not initialized")

type InterfaceInfo struct {
	Name       string   `json:"name"`
	IP         string   `json:"ip"`
	Zone       string   `json:"zone"`
	State      string   `json:"state"`
	Management []string `json:"management"`
}

func GetInterfaces() ([]InterfaceInfo, error) {

	fw := GetFirewall()
	if fw == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []InterfaceInfo

	for _, iface := range fw.Interfaces {
		out = append(out, InterfaceInfo{
			Name:       iface.Name,
			IP:         iface.IP,
			Zone:       iface.Zone,
			State:      iface.State,
			Management: iface.Management,
		})
	}

	return out, nil
}

func GetInterfacesL2() ([]InterfaceInfo, error) {
	fw := GetFirewall()
	if fw == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []InterfaceInfo
	for _, iface := range fw.Interfaces {
		if iface.IP == "" {
			out = append(out, InterfaceInfo{
				Name:       iface.Name,
				IP:         iface.IP,
				Zone:       iface.Zone,
				State:      iface.State,
				Management: iface.Management,
			})
		}
	}

	return out, nil
}

func GetInterfacesL3() ([]InterfaceInfo, error) {
	fw := GetFirewall()
	if fw == nil {
		return nil, ErrFirewallNotInitialized
	}

	var out []InterfaceInfo
	for _, iface := range fw.Interfaces {
		if iface.IP != "" {
			out = append(out, InterfaceInfo{
				Name:       iface.Name,
				IP:         iface.IP,
				Zone:       iface.Zone,
				State:      iface.State,
				Management: iface.Management,
			})
		}
	}

	return out, nil
}
