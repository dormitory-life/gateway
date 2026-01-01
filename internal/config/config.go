package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env         string        `yaml:"env"`
	Server      ServerConfig  `yaml:"server"`
	AuthService ServiceConfig `yaml:"auth_service"`
	CoreService ServiceConfig `yaml:"core_service"`
	JWT         JWTConfig     `yaml:"jwt"`
}

type ServerConfig struct {
	Port uint16 `yaml:"port"`
}

type ServiceConfig struct {
	Url string `yaml:"url"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

func ParseConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
