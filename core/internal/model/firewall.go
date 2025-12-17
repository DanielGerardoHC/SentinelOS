package model

type Firewall struct {
	Zones      map[string]*Zone
	Interfaces map[string]*Interface
	Addresses  map[string]*Address
	Services   map[string]*Service
	Routes     []*Route
	Policies   []*Policy
	NATRules   []*NATRule
	// Routes     []*Route
}
