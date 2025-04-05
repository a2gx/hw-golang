package app

import "time"

type Event struct {
	Id          string
	Title       string
	Datetime    time.Time
	Duration    time.Duration
	Description string
	UserId      string
	NotifyTime  time.Duration
}
