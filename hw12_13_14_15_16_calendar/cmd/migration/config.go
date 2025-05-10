package main

import (
	"flag"
	"fmt"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/config"
)

type Config struct {
	Logger struct {
		Level     string `mapstructure:"level"`
		Handler   string `mapstructure:"handler"`
		Filename  string `mapstructure:"filename"`
		AddSource bool   `mapstructure:"add_source"`
	} `mapstructure:"logger"`

	Database struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Dbname   string `mapstructure:"dbname"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
	} `mapstructure:"database"`

	DatabaseDNS string

	MigrationCommand string
	MigrationName    string
	MigrationDir     string
}

var (
	configPath    string
	command, name string
)

func init() {
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "Path to configuration file")
	flag.StringVar(&command, "command", "status", "Command for migration")
	flag.StringVar(&name, "name", "", "Name for new migration")
}

func NewConfig() (*Config, error) {
	instance := &Config{}
	if err := config.LoadConfig(instance, configPath); err != nil {
		return nil, err
	}

	instance.DatabaseDNS = fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		instance.Database.Username,
		instance.Database.Password,
		instance.Database.Dbname,
		instance.Database.Host,
		instance.Database.Port,
	)

	instance.MigrationDir = "./migrations"
	instance.MigrationCommand = command
	instance.MigrationName = name

	return instance, nil
}
