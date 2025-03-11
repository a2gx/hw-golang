package config

import (
	"log"
	"sync"
)

type SenderConfig struct {
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`

	// TODO
}

func GetSenderConfig() *SenderConfig {
	var (
		once     sync.Once
		instance = &SenderConfig{}
		pathname = "./configs/sender.yaml"
	)

	once.Do(func() {
		if err := loader[SenderConfig](instance, pathname); err != nil {
			log.Fatalf("load SenderConfig is failed, err: %v", err)
		}
	})

	return instance
}
