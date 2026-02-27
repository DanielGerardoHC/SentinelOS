package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
)

func ValidatePolicies(fw *model.Firewall) error {

	for _, p := range fw.Policies {

		if p.SrcZone == nil {
			return fmt.Errorf("policy missing src zone")
		}

		if p.DstZone == nil {
			return fmt.Errorf("policy missing dst zone")
		}

		for _, svc := range p.Services {
			if svc == nil {
				return fmt.Errorf("policy has nil service")
			}
		}
	}

	return nil
}