package config_engine

import "errors"

func BeginConfig(username string) error {

	configLock.Lock()
	defer configLock.Unlock()

	if configLock.locked {
		return errors.New("config already locked by " + configLock.owner)
	}

	clone, err := CloneFirewall(running)
	if err != nil {
		return err
	}

	candidate = clone
	configLock.locked = true
	configLock.owner = username

	return nil
}