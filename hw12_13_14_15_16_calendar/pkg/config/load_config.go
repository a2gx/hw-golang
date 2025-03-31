package config

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func setDefaultEnv(prefix string, conf interface{}) {
	valueOf, typeOf := reflect.ValueOf(conf), reflect.TypeOf(conf)

	if valueOf.Kind() == reflect.Ptr {
		valueOf, typeOf = valueOf.Elem(), typeOf.Elem()
	}

	for i := 0; i < typeOf.NumField(); i++ {
		fieldValue, fieldType := valueOf.Field(i), typeOf.Field(i)
		key := fieldType.Tag.Get("mapstructure")

		if prefix != "" {
			key = prefix + "." + key
		}

		if fieldValue.Kind() == reflect.Struct {
			setDefaultEnv(key, fieldValue.Addr().Interface())
		} else {
			viper.SetDefault(key, fieldValue.Interface())
		}
	}
}

func LoadConfig[T any](conf *T, filepath string) error {
	viper.Reset()

	viper.SetConfigFile(filepath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// устанавливаем пустые переменные окружения
	setDefaultEnv("", conf)

	// пытаемся получить конфиг из файла
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Printf("warning: error reading config file, %s", err)
		}
	}

	// привязываем флаги к конфигурации
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind command line flags: %w", err)
	}

	// переносим конфиг в структуру
	if err := viper.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
