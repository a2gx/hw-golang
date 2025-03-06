package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || strings.Contains(entry.Name(), "=") {
			log.Printf("Skipping %s\n", entry.Name())
			continue
		}

		path := filepath.Join(dir, entry.Name())
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		var line string
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			line = scanner.Text()
			line = strings.TrimRight(line, " \t\n\r")
			line = strings.ReplaceAll(line, "\x00", "\n")
		}

		_ = file.Close()

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		env[entry.Name()] = EnvValue{
			Value:      line,
			NeedRemove: line == "",
		}
	}

	return env, nil
}
