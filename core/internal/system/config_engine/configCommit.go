package config_engine

import (
	"errors"

	"sentinelos/core/internal/runtime"
	"sentinelos/core/internal/validator"
	"sentinelos/core/pkg/utils"
)

func Commit() error {

	configLock.Lock()
	defer configLock.Unlock()

	if candidate == nil {
		return errors.New("no candidate config")
	}

	// Validate
	if err := validator.ValidateFirewall(candidate); err != nil {
		abortConfigLocked()
		return err

	}

	// Backup
	backup, err := CloneFirewall(running)
	if err != nil {
		abortConfigLocked()

		return &utils.APIError{Code: "ERR_SYS_4003", Message: "Failed to create config backup", Details: err.Error()}
	}

	// Apply
	if err := runtime.ApplyFullRuntime(candidate); err != nil {
		_ = runtime.ApplyFullRuntime(backup)
		return &utils.APIError{Code: "ERR_SYS_4002", Message: "Failed to apply runtime configuration", Details: err.Error()}
	}

	// Success
	running = candidate

	if err := SaveYAML(running); err != nil {
		return err
	}

	abortConfigLocked()

	return nil
}
