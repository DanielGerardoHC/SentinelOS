package model

type NATType string

const (
    TypeSNAT     NATType = "snat"
    TypeDNATIP   NATType = "dnat-ip"
    TypeDNATPort NATType = "dnat-port"
)

type NATRule struct {
    ID             int
    Type           NATType
    SrcZone        *Zone
    DstZone        *Zone
    SrcAddress     string
    DstAddress     string
    Service        string
    OutInterface   string
    TranslatedIP   string 
    TranslatedPort string 
    Description    string
}