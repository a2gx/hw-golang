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
		Port     int    `mapstructure:"port"`
	} `mapstructure:"database"`

	DatabaseDns      string
	MigrationCommand string
	MigrationName    string
	MigrationDir     string
}

var configPath string
var command, name string

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

	dbDns := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		instance.Database.Username,
		instance.Database.Password,
		instance.Database.Dbname,
		"localhost",
		instance.Database.Port,
	)

	return &Config{
		MigrationDir:     "./migrations",
		MigrationCommand: command,
		MigrationName:    name,
		DatabaseDns:      dbDns,
		Logger:           instance.Logger,
	}, nil
}
