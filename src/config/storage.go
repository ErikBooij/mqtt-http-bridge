package config

import "github.com/kelseyhightower/envconfig"

type FileStorageConfig struct {
	Filename string `envconfig:"STORAGE_FILENAME" default:"storage.json"`
}

func LoadFileStorageConfig() (FileStorageConfig, error) {
	var cfg FileStorageConfig

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
