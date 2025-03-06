package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("testdata", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")

		require.NotNil(t, env)
		require.NoError(t, err)

		expected := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", true},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {`"hello"`, false},
			"UNSET": {"", true},
		}

		assert.Equal(t, expected, env)
	})

	t.Run("tempdata", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "tempdata")
		assert.NoError(t, err)

		files := map[string]string{
			"VAR1":      "value1\n",
			"VAR2":      "value2",
			"EMPTY":     "",
			"WITH_NULL": "line\x00break",
			".HIDDEN":   "should_ignore",
			"SUBDIR":    "",
		}

		for k, v := range files {
			path := filepath.Join(tempDir, k)
			if k == "SUBDIR" {
				assert.NoError(t, os.Mkdir(path, 0o755))
			} else {
				assert.NoError(t, os.WriteFile(path, []byte(v), 0o755))
			}
		}

		env, err := ReadDir(tempDir)

		require.NotNil(t, env)
		require.NoError(t, err)

		expected := Environment{
			"VAR1":      {Value: "value1", NeedRemove: false},
			"VAR2":      {Value: "value2", NeedRemove: false},
			"EMPTY":     {Value: "", NeedRemove: true},
			"WITH_NULL": {Value: "line\nbreak", NeedRemove: false},
		}

		assert.Equal(t, expected, env)
	})
}
