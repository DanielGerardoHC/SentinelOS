package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
)

func ValidateNAT(fw *model.Firewall) error {

	for _, n := range fw.NATRules {

		if n.SrcZone == nil {
			return fmt.Errorf("nat rule missing src zone")
		}

		if n.DstZone == nil {
			return fmt.Errorf("nat rule missing dst zone")
		}

		if _, ok := fw.Interfaces[n.OutInterface]; !ok {
			return fmt.Errorf("nat %d invalid outInterface", n.ID)
		}
	}

	return nil
}