package tools

import "time"

func GetDateInterval(date time.Time, add int) (start, finish time.Time) {
	// Нормализуем дату, сравниваем без учета времени
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return date, date.AddDate(0, 0, add)
}
