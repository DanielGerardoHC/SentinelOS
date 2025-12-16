package model

type NATAction string

const (
	Masquerade NATAction = "masquerade"
	SNAT       NATAction = "snat"
	DNAT       NATAction = "dnat"
)

type NATRule struct {
	ID           int
	Type         string
	SrcZone      *Zone
	DstZone      *Zone
	OutInterface string
	Action       NATAction
	Description  string
}
