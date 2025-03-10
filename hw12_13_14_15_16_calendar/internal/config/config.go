package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`

	Database struct {
		Username string `yaml:"username" envconfig:"DB_USERNAME"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	} `yaml:"database"`

	Logger struct {
		Level string `yaml:"level" envconfig:"LOG_LEVEL"`
	}
}

func NewConfig(path string) (*Config, error) {
	var conf Config

	if err := readFromFile(&conf, path); err != nil {
		return nil, err
	}
	if err := readFromEnv(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func readFromFile(cfg *Config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("error decoding file: %w", err)
	}

	return nil
}

func readFromEnv(cfg *Config) error {
	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("error processing env: %w", err)
	}
	return nil
}
