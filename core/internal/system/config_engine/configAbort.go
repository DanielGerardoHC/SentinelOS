package config_engine

func AbortConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	// logica interna
	abortConfigLocked()
}

// quien llame ya debe tener configLock.Lock() activado.
func abortConfigLocked() {
	candidate = nil
	configLock.locked = false
	configLock.owner = ""
}
