package tools

import (
	"testing"
	"time"
)

func TestGetDateInterval(t *testing.T) {
	testCases := []struct {
		name       string
		inputDate  time.Time
		addDays    int
		wantStart  time.Time
		wantFinish time.Time
	}{
		{
			name:       "Добавление 1 дня",
			inputDate:  time.Date(2024, 5, 15, 12, 30, 45, 0, time.UTC),
			addDays:    1,
			wantStart:  time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			wantFinish: time.Date(2024, 5, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Добавление 7 дней",
			inputDate:  time.Date(2024, 5, 15, 8, 0, 0, 0, time.UTC),
			addDays:    7,
			wantStart:  time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			wantFinish: time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Переход на следующий месяц",
			inputDate:  time.Date(2024, 5, 30, 22, 15, 0, 0, time.UTC),
			addDays:    3,
			wantStart:  time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC),
			wantFinish: time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Переход на следующий год",
			inputDate:  time.Date(2024, 12, 30, 10, 0, 0, 0, time.UTC),
			addDays:    5,
			wantStart:  time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC),
			wantFinish: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Отрицательное значение дней",
			inputDate:  time.Date(2024, 5, 15, 12, 30, 45, 0, time.UTC),
			addDays:    -5,
			wantStart:  time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			wantFinish: time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Нулевое значение дней",
			inputDate:  time.Date(2024, 5, 15, 12, 30, 45, 0, time.UTC),
			addDays:    0,
			wantStart:  time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			wantFinish: time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotStart, gotFinish := GetDateInterval(tc.inputDate, tc.addDays)

			if !gotStart.Equal(tc.wantStart) {
				t.Errorf("неправильная начальная дата: получено %v, ожидалось %v",
					gotStart.Format("2006-01-02 15:04:05"), tc.wantStart.Format("2006-01-02 15:04:05"))
			}

			if !gotFinish.Equal(tc.wantFinish) {
				t.Errorf("неправильная конечная дата: получено %v, ожидалось %v",
					gotFinish.Format("2006-01-02 15:04:05"), tc.wantFinish.Format("2006-01-02 15:04:05"))
			}
		})
	}
}
