package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

// SetDefaultEnv - устанавливает default в Env из структуры
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
			key = prefix + "." + key
		}

		if fieldValue.Kind() == reflect.Struct {
			setDefaultEnv(key, fieldValue.Interface())
		} else {
			viper.SetDefault(key, fieldValue.Interface())
		}
	}
}

func loader[T any](conf *T, defaultConfigPath string) error {
	viper.Reset()

	pflag.String("config", defaultConfigPath, "config file path")
	pflag.Parse()

	configPath := viper.GetString("config")
	if configPath == "" {
		configPath = defaultConfigPath // TODO need check
	}

	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// дефолт для переменных окружения
	setDefaultEnv("", conf)

	if err := viper.ReadInConfig(); err != nil {
		// read successfully
	}

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind pflags: %w", err)
	}

	return viper.Unmarshal(conf)
}
