package model

type Action string

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

type Policy struct {
	ID       int
	SrcZone  *Zone
	DstZone  *Zone
	SrcAddr  *Address
	DstAddr  *Address
	Services []*Service
	Action   Action
	Log      bool
}
