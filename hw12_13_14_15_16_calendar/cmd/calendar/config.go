package main

import (
	pkgconfig "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/config"
)

type Config struct {
	Logger struct {
		Level    string `mapstructure:"level"`
		Filename string `mapstructure:"filename"`
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
	var instance = &Config{}
	var pathname = "./configs/config.yaml"

	if err := pkgconfig.LoadConfig(instance, pathname); err != nil {
		return nil, err
	}

	return instance, nil
}
