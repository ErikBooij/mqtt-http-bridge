package config

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"strings"
)

type Config struct {
	// Server
	BindAddress string `envconfig:"BIND" default:"127.0.0.1"`
	BindPort    int    `envconfig:"PORT" default:"1883"`

	// Auth
	OpenAuth bool   `envconfig:"OPEN_AUTH" default:"false"`
	Users    string `envconfig:"USERS" default:""`

	// Store
	StorageDriver string `envconfig:"STORAGE_DRIVER" default:"memory"`

	// Parsed
	UsersParsed []ParsedUser
}

type ParsedUser struct {
	Username string
	Password string
}

func Load() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	if !cfg.OpenAuth {
		for _, user := range strings.Split(cfg.Users, ",") {
			parts := strings.Split(user, ":")
			if len(parts) != 2 {
				return cfg, errors.New("invalid user format")
			}

			cfg.UsersParsed = append(cfg.UsersParsed, ParsedUser{
				Username: strings.TrimSpace(parts[0]),
				Password: strings.TrimSpace(parts[1]),
			})
		}

		if len(cfg.UsersParsed) == 0 {
			return cfg, errors.New("no users defined")
		}
	}

	return cfg, nil
}
