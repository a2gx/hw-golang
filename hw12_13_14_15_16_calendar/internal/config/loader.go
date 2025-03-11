package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// Loader загружает конфигурацию в переданную структуру `conf` из файла и переменных окружения
func loader[T any](conf *T, pathname string) error {
	pflag.String("config", pathname, "config file path")
	pflag.Parse()

	configPath := viper.GetString("config")
	if configPath == "" {
		configPath = pathname
	}

	// checking file
	if _, err := os.Stat(configPath); err == nil {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("warning: failed to read config file: %v\n", err)
		} else {
			log.Printf("success: loaded config from file: %s\n", configPath)
		}
	} else if !os.IsNotExist(err) {
		log.Printf("warning: config file [%s] not exists\n", configPath)
	}

	// environments
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind pflags: %w", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
