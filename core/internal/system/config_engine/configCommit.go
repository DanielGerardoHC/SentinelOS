package config_engine

import (
	"errors"

	"sentinelos/core/internal/runtime"
	"sentinelos/core/internal/validator"
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
		return err
	}

	// Apply
	if err := runtime.ApplyFullRuntime(candidate); err != nil {
		_ = runtime.ApplyFullRuntime(backup)
		return err
	}

	// Success
	running = candidate

	if err := SaveYAML(running); err != nil {
		return err
	}

	abortConfigLocked()

	return nil
}
