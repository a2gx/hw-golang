package app

import "time"

type Event struct {
	ID          string
	Title       string
	Datetime    time.Time
	Duration    time.Duration
	Description string
	UserID      string
	NotifyTime  time.Duration
}
