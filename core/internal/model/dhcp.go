package model

type DHCP struct {
	Interface    string
	StartIP      string
	EndIP        string
	SubnetMask   string
	Gateway      string
	DNS          []string
	LeaseTimeMin int
}
