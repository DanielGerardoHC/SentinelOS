package auth

import (
	"os"

	"gopkg.in/yaml.v3"
)

type UsersFile struct {
	Users []User          `yaml:"users"`
	Roles map[string]Role `yaml:"roles"`
}

type User struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
	Role         string `yaml:"role"`
	Enabled      bool   `yaml:"enabled"`
}

type Role struct {
	Permissions []string `yaml:"permissions"`
}

func LoadUsers(path string) (*UsersFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var users UsersFile
	if err := yaml.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return &users, nil
}
