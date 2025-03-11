package config

import (
	"log"
	"sync"
)

type SchedulerConfig struct {
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`

	// TODO
}

func GetSchedulerConfig() *SchedulerConfig {
	var (
		once     sync.Once
		instance = &SchedulerConfig{}
		pathname = "./configs/scheduler.yaml"
	)

	once.Do(func() {
		if err := loader[SchedulerConfig](instance, pathname); err != nil {
			log.Fatalf("load SchedulerConfig is failed, err: %v", err)
		}
	})

	return instance
}
