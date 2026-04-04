package system

import "strings"

type NatInfo struct {
    ID             int    `json:"id"`
    Type           string `json:"type"`
    SrcZone        string `json:"src-zone"`
    DstZone        string `json:"dst-zone"`
    SrcAddress     string `json:"src-addr"`
    DstAddress     string `json:"dst-addr"`
    Service        string `json:"service"`
    OutInterface   string `json:"out-interface"`
    TranslatedIP   string `json:"translated-ip"`
    TranslatedPort string `json:"translated-port"`
    Description    string `json:"description"`
}

func GetNatRules(typeFilter string) ([]NatInfo, error) {
    fw := GetFirewall()
    if fw == nil {
        return nil, ErrFirewallNotInitialized
    }

    var out []NatInfo = []NatInfo{}

    for _, nat := range fw.NATRules {
        if typeFilter != "" && !strings.HasPrefix(string(nat.Type), typeFilter) {
            continue
        }
        out = append(out, NatInfo{
            ID:   nat.ID,
            Type: string(nat.Type),
            SrcZone: func() string {
                if nat.SrcZone != nil { return nat.SrcZone.Name }
                return ""
            }(),
            DstZone: func() string {
                if nat.DstZone != nil { return nat.DstZone.Name }
                return ""
            }(),
            SrcAddress:     nat.SrcAddress,
            DstAddress:     nat.DstAddress,
            Service:        nat.Service,
            OutInterface:   nat.OutInterface,
            TranslatedIP:   nat.TranslatedIP,
            TranslatedPort: nat.TranslatedPort,
            Description:    nat.Description,
        })
    }

    return out, nil
}