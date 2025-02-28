package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(RunCmd(cmd, env))
}
