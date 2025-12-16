package model

import "net"

type Interface struct {
	Name       string
	Zone       *Zone
	IP         *net.IPNet
	Management []string
}
