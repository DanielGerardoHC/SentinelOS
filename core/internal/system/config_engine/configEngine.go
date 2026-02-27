package config_engine

import (
	"sync"

	"sentinelos/core/internal/model"
)

var running *model.Firewall
var candidate *model.Firewall

var configLock struct {
	sync.Mutex
	locked bool
	owner  string
}

func SetFirewall(fw *model.Firewall) {
	running = fw
}

func GetRunning() *model.Firewall {
	return running
}

func GetCandidate() *model.Firewall {
	return candidate
}