package metrics

import (
	"testing"
	"time"
)

func d(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
}

func TestBusinessDaysBetween(t *testing.T) {
	tests := []struct {
		name string
		from time.Time
		to   time.Time
		want float64
	}{
		{name: "same day weekday", from: d(2026, 3, 2), to: d(2026, 3, 2), want: 0},
		{name: "same day weekend", from: d(2026, 3, 1), to: d(2026, 3, 1), want: 0},
		{name: "monday to tuesday", from: d(2026, 3, 2), to: d(2026, 3, 3), want: 1},
		{name: "monday to friday", from: d(2026, 3, 2), to: d(2026, 3, 6), want: 4},
		{name: "friday to monday", from: d(2026, 3, 6), to: d(2026, 3, 9), want: 1},
		{name: "friday to next friday", from: d(2026, 3, 6), to: d(2026, 3, 13), want: 5},
		{name: "two full weeks", from: d(2026, 3, 2), to: d(2026, 3, 16), want: 10},
		{name: "saturday to monday", from: d(2026, 3, 7), to: d(2026, 3, 9), want: 1},
		{name: "saturday to sunday", from: d(2026, 3, 7), to: d(2026, 3, 8), want: 0},
		{name: "sunday to monday", from: d(2026, 3, 8), to: d(2026, 3, 9), want: 1},
		{name: "wednesday to next wednesday", from: d(2026, 3, 4), to: d(2026, 3, 11), want: 5},
		{name: "four week span", from: d(2026, 3, 2), to: d(2026, 3, 30), want: 20},
		{name: "thursday to tuesday across weekend", from: d(2026, 3, 5), to: d(2026, 3, 10), want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BusinessDaysBetween(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("BusinessDaysBetween(%s, %s) = %v, want %v",
					tt.from.Format("Mon 2006-01-02"),
					tt.to.Format("Mon 2006-01-02"),
					got, tt.want)
			}
		})
	}
}
