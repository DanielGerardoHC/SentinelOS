package system

func FirewallRunning() bool {
	return firewall != nil
}

func InterfacesCount() int {
	if firewall == nil {
		return 0
	}
	return len(firewall.Interfaces)
}

func RoutesCount() int {
	if firewall == nil {
		return 0
	}
	return len(firewall.Routes)
}

func DHCPRunning() bool {
	if firewall == nil {
		return false
	}
	return len(firewall.DHCPConfigs) > 0
}
