package model

type Vlan struct {
	Name       string
	Parent     string
	ID         int
	IP         string
	Zone       string
	State      string
	Management []string
}
