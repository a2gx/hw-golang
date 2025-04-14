package main

import (
	"flag"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/config"
)

type Config struct {
	Server struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"server"`
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
