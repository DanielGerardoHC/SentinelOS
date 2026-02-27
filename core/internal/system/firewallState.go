/*package system

import "sentinelos/core/internal/model"

var firewall *model.Firewall

func SetFirewall(fw *model.Firewall) {
	firewall = fw
}

func GetFirewall() *model.Firewall {
	return firewall
}
*/

package system

import (
	"sentinelos/core/internal/model"
	"sentinelos/core/internal/system/config_engine"
)

func SetFirewall(fw *model.Firewall) {
	config_engine.SetFirewall(fw)
}

func GetFirewall() *model.Firewall {
	return config_engine.GetRunning()
}
