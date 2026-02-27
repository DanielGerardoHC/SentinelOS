package config_engine

import (
	"os"
	"gopkg.in/yaml.v3"
	"sentinelos/core/internal/model"
)

func SaveYAML(fw *model.Firewall) error {

	data, err := yaml.Marshal(fw)
	if err != nil {
		return err
	}

	return os.WriteFile("/srv/sentinelos/core/internal/config/config.yml", data, 0644)
}