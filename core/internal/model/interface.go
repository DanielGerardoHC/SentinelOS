package model

import "net"

type Interface struct {
	Name       string
	IP         *net.IPNet
	Management []string
}
