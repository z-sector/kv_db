package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type NetworkConfig struct {
	Address        string `yaml:"address"`
	MaxConnections int    `yaml:"max_connections"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func Load(filename string) (*Config, error) {
	if filename == "" {
		return &Config{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
