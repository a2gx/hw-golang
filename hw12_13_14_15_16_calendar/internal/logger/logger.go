package logger

import "fmt"

type Logger struct { // TODO
}

func GetLogger(level string) *Logger {
	return &Logger{}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

func (l Logger) Close() {
	// TODO
}

// TODO
