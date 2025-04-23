package app

import "errors"

var (
	ErrIDRequired = errors.New("event ID is required")
	ErrNotFound   = errors.New("event not found")
)
