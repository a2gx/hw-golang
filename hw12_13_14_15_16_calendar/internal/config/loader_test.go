package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testConfig struct {
	Name  string `mapstructure:"name"`
	Port  int    `mapstructure:"port"`
	Debug bool   `mapstructure:"debug"`
}

const testContent = `
name: "test-app"
port: 8080
debug: true
`

func createTempConfig(content string) (string, error) {
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

func TestLoader(t *testing.T) {
	os.Setenv("PORT", "9090")
	defer os.Clearenv()

	configPath, err := createTempConfig(testContent)
	require.NoError(t, err)
	defer os.Remove(configPath)

	var conf testConfig
	err = loader(&conf, configPath)

	require.NoError(t, err)
	require.Equal(t, "test-app", conf.Name) // from file
	require.Equal(t, 9090, conf.Port)       // from env
	require.Equal(t, true, conf.Debug)      // from file
}
