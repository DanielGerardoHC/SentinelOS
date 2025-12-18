package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RawConfig struct {
	Zones []struct {
		Name       string   `yaml:"name"`
		Type       string   `yaml:"type"`
		Interfaces []string `yaml:"interfaces"`
		Networks   []string `yaml:"networks"`
	} `yaml:"zones"`

	Interfaces []struct {
		Name       string   `yaml:"name"`
		IP         string   `yaml:"ip"`
		Zone       string   `yaml:"zone"`
		State      string   `yaml:"state"`
		Management []string `yaml:"management"`
	} `yaml:"interfaces"`

	Vlans []struct {
		Name       string   `yaml:"name"`
		Parent     string   `yaml:"parent"`
		ID         int      `yaml:"id"`
		IP         string   `yaml:"ip"`
		Zone       string   `yaml:"zone"`
		State      string   `yaml:"state"`
		Management []string `yaml:"management"`
	} `yaml:"vlans"`

	Addresses []struct {
		Name string   `yaml:"name"`
		IPs  []string `yaml:"ips"`
	} `yaml:"addresses"`

	Services []struct {
		Name     string `yaml:"name"`
		Protocol string `yaml:"protocol"`
		Ports    []int  `yaml:"ports"`
	} `yaml:"services"`

	Route []struct {
		ID          int    `yaml:"id"`
		Destination string `yaml:"destination"`
		Gateway     string `yaml:"gateway"`
		Interface   string `yaml:"interface"`
		Metric      int    `yaml:"metric"`
		Description string `yaml:"description"`
	} `yaml:"routes"`

	Policies []struct {
		ID       int      `yaml:"id"`
		SrcZone  string   `yaml:"src-zone"`
		DstZone  string   `yaml:"dst-zone"`
		SrcAddr  string   `yaml:"src-addr"`
		DstAddr  string   `yaml:"dst-addr"`
		Services []string `yaml:"services"`
		Action   string   `yaml:"action"`
		Log      bool     `yaml:"log"`
	} `yaml:"policies"`

	NATRules []struct {
		ID           int    `yaml:"id"`
		Type         string `yaml:"type"`
		SrcZone      string `yaml:"src-zone"`
		DstZone      string `yaml:"dst-zone"`
		Action       string `yaml:"action"`
		OutInterface string `yaml:"outInterface"`
		Description  string `yaml:"description"`
	} `yaml:"nat"`
}

func LoadRawConfig(path string) (*RawConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg RawConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
