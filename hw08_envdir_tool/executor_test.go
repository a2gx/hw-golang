package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env := Environment{
			"TEXT_ENV": {"value from env - success", false},
		}
		cmd := []string{"bash", "-c", "echo $TEXT_ENV"}
		exc := RunCmd(cmd, env)

		require.Equal(t, 0, exc)
	})

	t.Run("remove_env_var", func(t *testing.T) {
		_ = os.Setenv("TEXT_ENV", "value from env - success")

		env := Environment{
			"TEXT_ENV": {NeedRemove: false},
		}
		cmd := []string{"bash", "-c", "echo $TEXT_ENV"}
		exc := RunCmd(cmd, env)

		require.Equal(t, 0, exc)
	})

	t.Run("invalid_command", func(t *testing.T) {
		env := Environment{}
		cmd := []string{"nonexistent_command"}
		exc := RunCmd(cmd, env)

		require.Equal(t, 1, exc)
	})

	t.Run("empty_command", func(t *testing.T) {
		env := Environment{}
		cmd := []string{}
		exc := RunCmd(cmd, env)

		require.Equal(t, 1, exc)
	})
}
