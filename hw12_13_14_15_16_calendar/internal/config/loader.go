package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SetDefaultEnv устанавливает значения по умолчанию в Viper из структуры
func setDefaultEnv(prefix string, conf interface{}) {
	valueOf, typeOf := reflect.ValueOf(conf), reflect.TypeOf(conf)

	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
		typeOf = typeOf.Elem()
	}

	for i := 0; i < typeOf.NumField(); i++ {
		fieldValue, fieldType := valueOf.Field(i), typeOf.Field(i)
		key := fieldType.Tag.Get("mapstructure")

		// генерируем ключ, например: "database.host"
		if prefix != "" {
			key = prefix + "." + strings.ToUpper(key)
		}

		if fieldValue.Kind() == reflect.Struct {
			setDefaultEnv(key, fieldValue.Interface())
		} else {
			viper.SetDefault(key, fieldValue.Interface())
		}
	}
}

// Loader загружает конфигурацию в переданную структуру `conf` из файла и переменных окружения
func loader[T any](conf *T, defaultConfigPath string) error {
	viper.Reset()

	pflag.String("config", defaultConfigPath, "config file path")
	pflag.Parse()

	configPath := viper.GetString("config")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// устанавливаем значения по умолчанию из структуры
	setDefaultEnv("", conf)

	// пытаемся прочитать конфигурацию из файла
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("warning: error reading config file: %v", err)
		}
	} else {
		log.Printf("info: loaded config from file: %s", configPath)
	}

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind command line flags: %w", err)
	}

	if err := viper.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
