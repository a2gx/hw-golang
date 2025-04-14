package app

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	UserID      string
	StartTime   time.Time
	EndTime     time.Time
	NotifyTime  time.Time
}
