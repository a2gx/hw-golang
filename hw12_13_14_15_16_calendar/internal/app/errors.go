package app

import "errors"

var (
	ErrDateBusy     = errors.New("date is busy")
	ErrNotFound     = errors.New("event not found")
	ErrAlreadyExist = errors.New("event already exist")
)
