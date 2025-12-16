package model

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

type Service struct {
	Name     string
	Protocol Protocol
	Ports    []int
}
