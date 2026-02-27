package system

func FirewallRunning() bool {
	
	fw := GetFirewall()
	return fw != nil
}

func InterfacesCount() int {
	fw := GetFirewall()
	if fw == nil {
		return 0
	}
	return len(fw.Interfaces)
}

func RoutesCount() int {
	fw := GetFirewall()
	if fw == nil {
		return 0
	}
	return len(fw.Routes)
}

func DHCPRunning() bool {
	fw := GetFirewall()
	if fw == nil {
		return false
	}
	return len(fw.DHCPConfigs) > 0
}
