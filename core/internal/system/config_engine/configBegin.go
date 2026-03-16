package config_engine

import (
	"fmt"

	"sentinelos/core/pkg/utils"
)

func BeginConfig(username string) error {
	configLock.Lock()
	defer configLock.Unlock()

	if configLock.locked {
		return &utils.APIError{
			Code:    "ERR_SYS_3002",
			Message: "Config already locked by another user",
			Details: fmt.Sprintf("owner: %s", configLock.owner),
		}
	}

	clone, err := CloneFirewall(running)
	if err != nil {
		return &utils.APIError{
			Code:    "ERR_SYS_4003",
			Message: "Failed to create config backup",
			Details: err.Error(),
		}
	}

	candidate = clone
	configLock.locked = true
	configLock.owner = username

	return nil
}
