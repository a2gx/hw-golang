package main

import (
	pkgconfig "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/config"
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

func NewConfig() (*Config, error) {
	instance := &Config{}
	pathname := "./configs/config.yaml"

	if err := pkgconfig.LoadConfig(instance, pathname); err != nil {
		return nil, err
	}

	return instance, nil
}
