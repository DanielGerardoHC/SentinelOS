package validator

import (
	"sentinelos/core/internal/model"
)

func ValidateFirewall(fw *model.Firewall) error {
	if err := ValidateZones(fw); err != nil {
		return err
	}

	if err := ValidateInterfaces(fw); err != nil {
		return err
	}

	if err := ValidateVlans(fw); err != nil {
		return err
	}

	if err := ValidateAddresses(fw); err != nil {
		return err
	}

	if err := ValidateServices(fw); err != nil {
		return err
	}

	if err := ValidateDHCP(fw); err != nil {
		return err
	}

	if err := ValidateRoutes(fw); err != nil {
		return err
	}

	if err := ValidatePolicies(fw); err != nil {
		return err
	}

	if err := ValidateNAT(fw); err != nil {
		return err
	}
	return nil
}
