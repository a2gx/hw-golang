package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type TestConfigNestedData struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type TestConfigNested struct {
	Server   TestConfigNestedData `mapstructure:"server"`
	Database TestConfigNestedData `mapstructure:"database"`
}

const nestedFileContent = `
server:
 host: "server.host"
 port: "8081"

database:
 host: "database.host"
 port: "5432"
`

func TestLoading(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		env      map[string]string
		expected TestConfig
		wantErr  bool
	}{
		{
			name:   "loading file",
			config: "host: one.text\nport: 1234\ndatabase: test_1",
			expected: TestConfig{
				Host:     "one.text",
				Port:     1234,
				Database: "test_1",
			},
		},
		{
			name: "loading envelopment",
			env: map[string]string{
				"HOST":     "two.test",
				"PORT":     "4321",
				"DATABASE": "test_2",
			},
			expected: TestConfig{
				Host:     "two.test",
				Port:     4321,
				Database: "test_2",
			},
		},
		{
			name:   "combined config data",
			config: "host: three.text\nport: 1234\ndatabase: test",
			env: map[string]string{
				"PORT":     "5432",
				"DATABASE": "db_test",
			},
			expected: TestConfig{
				Host:     "three.text",
				Port:     5432,
				Database: "db_test",
			},
		},
		{
			name:     "invalid file",
			config:   "host: localhost\nport: non-a-number\ndatabase: test",
			expected: TestConfig{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// очищаем флаги перед тестом
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

			// создаем файл
			var configFile string
			if tt.config != "" {
				tmpFile, err := os.CreateTemp("", "config-*.yaml")
				require.NoError(t, err)
				defer os.Remove(tmpFile.Name())

				_, err = tmpFile.Write([]byte(tt.config))
				require.NoError(t, err)
				configFile = tmpFile.Name()
			}

			// переменные окружения
			for k, v := range tt.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// загружаем конфигурацию
			var conf TestConfig
			err := loader(&conf, configFile)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, conf)
		})
	}
}

func TestLoadingNestedStruct(t *testing.T) {
	// очищаем флаги перед тестом
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// создаем файл
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(nestedFileContent))
	require.NoError(t, err)

	os.Setenv("SERVER_HOST", "server.env")
	os.Setenv("DATABASE_HOST", "database.env")

	var conf TestConfigNested
	err = loader[TestConfigNested](&conf, tmpFile.Name())

	expected := TestConfigNested{
		Server: TestConfigNestedData{
			Host: "server.env",
			Port: 8081,
		},
		Database: TestConfigNestedData{
			Host: "database.env",
			Port: 5432,
		},
	}

	require.NoError(t, err)
	require.Equal(t, expected, conf)
}
