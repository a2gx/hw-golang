package app

import "errors"

var (
	ErrIdRequired   = errors.New("event ID is required")
	ErrDateBusy     = errors.New("date is busy")
	ErrNotFound     = errors.New("event not found")
	ErrAlreadyExist = errors.New("event already exist")
)
