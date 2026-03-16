package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidateNAT(fw *model.Firewall) error {
	for _, n := range fw.NATRules {
		if n.SrcZone == nil {
			return &utils.APIError{Code: "ERR_NET_1002", Message: "Missing required field", Details: "nat rule missing src zone"}
		}

		if n.DstZone == nil {
			return &utils.APIError{Code: "ERR_NET_1002", Message: "Missing required field", Details: "nat rule missing dst zone"}
		}

		if _, ok := fw.Interfaces[n.OutInterface]; !ok {
			return &utils.APIError{Code: "ERR_NET_2003", Message: "Resource references unknown entity", Details: fmt.Sprintf("nat %d invalid outInterface", n.ID)}
		}
	}

	return nil
}
