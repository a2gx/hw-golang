package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testYAML = `
server:
  host: "localhost"
  port: "8080"
`

func createTempConfigFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func TestNewConfig(t *testing.T) {
	t.Run("loading from file", func(t *testing.T) {
		tmpFile, err := createTempConfigFile(testYAML)
		require.NoError(t, err)
		defer os.Remove(tmpFile)

		var conf Config
		err = readFromFile(&conf, tmpFile)
		require.NoError(t, err)

		require.Equal(t, conf.Server.Host, "localhost")
		require.Equal(t, conf.Server.Port, "8080")
	})

	t.Run("loading from env", func(t *testing.T) {
		_ = os.Setenv("SERVER_HOST", "envhost")
		_ = os.Setenv("SERVER_PORT", "9090")
		defer os.Clearenv()

		var conf Config
		err := readFromEnv(&conf)
		require.NoError(t, err)

		require.Equal(t, conf.Server.Host, "envhost")
		require.Equal(t, conf.Server.Port, "9090")
	})

	t.Run("loading full config", func(t *testing.T) {
		tmpFile, err := createTempConfigFile(testYAML)
		require.NoError(t, err)
		defer os.Remove(tmpFile)

		_ = os.Setenv("SERVER_HOST", "overridehost")
		defer os.Clearenv()

		conf, err := NewConfig(tmpFile)
		require.NoError(t, err)

		// переменные окружения должны перегружать файл YAML
		require.Equal(t, "overridehost", conf.Server.Host) // from Env
		require.Equal(t, "8080", conf.Server.Port)         // from Yaml
	})
}
