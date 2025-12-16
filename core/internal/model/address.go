package model

import "net"

type Address struct {
	Name string
	IPs  []net.IPNet
}
