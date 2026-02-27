package validator

import (
	"fmt"
	"sentinelos/core/internal/model"
)

func ValidateFirewall(fw *model.Firewall) error {

	if err := ValidateZones(fw); err != nil {
		return fmt.Errorf("zones: %w", err)
	}

	if err := ValidateInterfaces(fw); err != nil {
		return fmt.Errorf("interfaces: %w", err)
	}

	if err := ValidateVlans(fw); err != nil {
		return fmt.Errorf("vlans: %w", err)
	}

	if err := ValidateAddresses(fw); err != nil {
		return fmt.Errorf("addresses: %w", err)
	}

	if err := ValidateServices(fw); err != nil {
		return fmt.Errorf("services: %w", err)
	}

	if err := ValidateDHCP(fw); err != nil {
		return fmt.Errorf("dhcp: %w", err)
	}

	if err := ValidateRoutes(fw); err != nil {
		return fmt.Errorf("routes: %w", err)
	}

	if err := ValidatePolicies(fw); err != nil {
		return fmt.Errorf("policies: %w", err)
	}

	if err := ValidateNAT(fw); err != nil {
		return fmt.Errorf("nat: %w", err)
	}

	return nil
}