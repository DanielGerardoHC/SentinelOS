package config_engine

func AbortConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	abortConfigLocked()
}

func abortConfigLocked() {
	candidate = nil
	configLock.locked = false
	configLock.owner = ""
}
