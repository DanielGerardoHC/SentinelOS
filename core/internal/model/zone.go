package model

type ZoneType string

const (
	ZoneL2 ZoneType = "l2"
	ZoneL3 ZoneType = "l3"
)

type Zone struct {
	Name       string
	Type       ZoneType
	Interfaces []string
	Networks   []string
}
