package config

import (
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"mqtt-http-bridge/src/processor"
	"os"
	"slices"
	"strings"
)

type Config struct {
	// Meta
	AppEnv      string `envconfig:"APP_ENV" default:"production"`
	PrepareData bool   `envconfig:"PREPARE_DATA" default:"false"`

	Broker          BrokerConfig                    `yaml:"broker"`
	ExternalBrokers map[string]ExternalBrokerConfig `yaml:"external-brokers"`
	Server          ServerConfig                    `yaml:"server"`
	Storage         StorageConfig                   `yaml:"storage"`

	// Internal Options
	Silent bool
}

type BrokerConfig struct {
	Address  string       `yaml:"bind-address" default:"0.0.0.0"`
	Port     int          `yaml:"port" default:"1883"`
	OpenAuth bool         `yaml:"open-auth" default:"false"`
	Users    []BrokerUser `yaml:"users" default:""`
}

type BrokerUser struct {
	Username string
	Password string
}

type ExternalBrokerConfig struct {
	Name     string   `yaml:"name"`
	ClientID string   `yaml:"client-id"`
	Host     string   `yaml:"host"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Topics   []string `yaml:"topics"`
}

type ServerConfig struct {
	Address string `yaml:"bind-address" default:"0.0.0.0"`
	Port    int    `yaml:"port" default:"8080"`
}

var supportedStorageDrivers = []string{"memory", "file"}

type StorageConfig struct {
	Driver  string                 `yaml:"driver"`
	Options map[string]interface{} `yaml:"options"`
}

type StorageConfigFile struct {
	File string `yaml:"file"`
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "dev"
}

type ParsedUser struct {
	Username string
	Password string
}

func Load() (*Config, error) {
	var cfg Config

	confFile := "mqtt-http.conf.yaml"

	if f := os.Getenv("CONFIG_FILE"); f != "" {
		confFile = f
	}

	fileContents, err := os.ReadFile(confFile)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(fileContents, &cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	if !cfg.Broker.OpenAuth {
		for idx, user := range cfg.Broker.Users {
			username := strings.TrimSpace(user.Username)
			password := strings.TrimSpace(user.Password)

			if username == "" || password == "" {
				return nil, fmt.Errorf("invalid user configured for built-in broker at index #%d", idx)
			}

			cfg.Broker.Users[idx].Username = username
			cfg.Broker.Users[idx].Password = password
		}

		if len(cfg.Broker.Users) == 0 {
			return nil, errors.New("no users defined")
		}
	}

	if !slices.Contains(supportedStorageDrivers, cfg.Storage.Driver) {
		return nil, fmt.Errorf("invalid storage driver: %s (should be one of %s)", cfg.Storage.Driver, strings.Join(supportedStorageDrivers, "/"))
	}

	if _, ok := cfg.ExternalBrokers[processor.InternalBroker]; ok {
		return nil, fmt.Errorf("the name %s cannot be used for an external broker", processor.InternalBroker)
	}

	return &cfg, nil
}

func (c *Config) StorageConfigFile() (StorageConfigFile, error) {
	if c.Storage.Driver != "file" {
		return StorageConfigFile{}, errors.New("storage driver is not 'file'")
	}

	var scf StorageConfigFile

	if err := mapstructure.Decode(c.Storage.Options, &scf); err != nil {
		return StorageConfigFile{}, fmt.Errorf("unable to decode storage options: %w", err)
	}

	return scf, nil
}
