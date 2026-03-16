package validator

import (
	"sentinelos/core/internal/model"
	"sentinelos/core/pkg/utils"
)

func ValidatePolicies(fw *model.Firewall) error {
	for _, p := range fw.Policies {
		/*
		   if p.SrcZone == nil {
		      return &utils.APIError{Code: "ERR_NET_1002", Message: "Missing required field", Details: "policy missing src zone"}
		   }

		   if p.DstZone == "" {
		      return &utils.APIError{Code: "ERR_NET_1002", Message: "Missing required field", Details: "policy missing dst zone"}
		   }
		*/
		for _, svc := range p.Services {
			if svc == nil {
				return &utils.APIError{Code: "ERR_SEC_2001", Message: "Invalid policy reference", Details: "policy has nil service"}
			}
		}
	}

	return nil
}
