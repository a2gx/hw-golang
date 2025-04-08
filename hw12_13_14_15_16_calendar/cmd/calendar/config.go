package main

import (
	"flag"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/config"
)

type Config struct {
	App struct {
		Storage string `toml:"storage"`
		Server  string `toml:"server"`
	}

	Logger struct {
		Level       string `mapstructure:"level"`
		HandlerType string `mapstructure:"handler_type"`
		Filename    string `mapstructure:"filename"`
		AddSource   bool   `mapstructure:"add_source"`
	} `mapstructure:"logger"`

	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"database"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func NewConfig() (*Config, error) {
	instance := &Config{}

	if err := config.LoadConfig(instance, configFile); err != nil {
		return nil, err
	}

	return instance, nil
}
