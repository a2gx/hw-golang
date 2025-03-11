package config

import (
	"log"
	"sync"
)

type CalendarConfig struct {
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`

	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
}

func GetCalendarConfig() *CalendarConfig {
	var (
		once     sync.Once
		instance = &CalendarConfig{}
		pathname = "./configs/calendar.yaml"
	)

	once.Do(func() {
		if err := loader[CalendarConfig](instance, pathname); err != nil {
			log.Fatalf("load CalendarConfig is failed, err: %v", err)
		}
	})

	return instance
}
