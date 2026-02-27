package config_engine

func AbortConfig() {

	configLock.Lock()
	defer configLock.Unlock()

	candidate = nil
	configLock.locked = false
	configLock.owner = ""
}