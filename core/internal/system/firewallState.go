package system

import "sentinelos/core/internal/model"

var firewall *model.Firewall

func SetFirewall(fw *model.Firewall) {
	firewall = fw
}

func GetFirewall() *model.Firewall {
	return firewall
}
