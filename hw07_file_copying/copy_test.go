package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	listExpect := []struct {
		toPath string
		offset int64
		limit  int64
	}{
		{"testdata/out_offset0_limit0.txt", 0, 0},
		{"testdata/out_offset0_limit10.txt", 0, 10},
		{"testdata/out_offset0_limit1000.txt", 0, 1000},
		{"testdata/out_offset0_limit10000.txt", 0, 10000},
		{"testdata/out_offset100_limit1000.txt", 100, 1000},
		{"testdata/out_offset6000_limit1000.txt", 6000, 1000},
	}

	for _, val := range listExpect {
		err := Copy("testdata/input.txt", "out.txt", val.offset, val.limit)
		require.NoError(t, err)
		file, err := os.ReadFile("out.txt")
		require.NoError(t, err)
		expected, err := os.ReadFile(val.toPath)
		require.NoError(t, err)

		require.Equal(t, expected, file)

		err = os.Remove("out.txt")
		require.NoError(t, err)
	}
}
