package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Printf("error: command is empty\n")
		return 1
	}

	// set environments
	for key, val := range env {
		if val.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				log.Printf("error: failed to unset env variable %s: %v\n", key, err)
				return 1
			}
		} else {
			if err := os.Setenv(key, val.Value); err != nil {
				log.Printf("error: failed to set env variable %s: %v\n", key, err)
				return 1
			}
		}
	}

	executable, err := exec.LookPath(cmd[0])
	if err != nil {
		log.Printf("error: failed to find executable: %v\n", cmd[0])
		return 1
	}

	command := exec.Command(executable, cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		exitCode := command.ProcessState.ExitCode()
		log.Printf("error: command %s exited with %d: %v\n", cmd, exitCode, err)
		return exitCode
	}

	return 0
}
