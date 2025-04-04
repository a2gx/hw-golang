package main

import (
	pkgconfig "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/config"
)

type Config struct {
	Logger struct {
		Level    string `mapstructure:"level"`
		Handler  string `mapstructure:"handler"`
		Filename string `mapstructure:"filename"`
		Source   bool   `mapstructure:"source"`
	} `mapstructure:"logger"`

	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Storage struct {
		Source string `mapstructure:"source"`
	} `mapstructure:"storage"`

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
