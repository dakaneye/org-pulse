// Package metrics provides pure computation functions for delivery health metrics.
package metrics

import "time"

// BusinessDaysBetween returns fractional business days between from and to.
// Only weekday (Mon-Fri) hours count. Returns 0 if from >= to.
func BusinessDaysBetween(from, to time.Time) float64 {
	if !from.Before(to) {
		return 0
	}

	var weekdayHours float64

	current := from
	for current.Before(to) {
		if isWeekday(current) {
			// Advance to end of this calendar day or to 'to', whichever is sooner
			endOfDay := startOfNextDay(current)
			if endOfDay.After(to) {
				weekdayHours += to.Sub(current).Hours()
			} else {
				weekdayHours += endOfDay.Sub(current).Hours()
			}
			current = endOfDay
		} else {
			// Skip entire weekend day
			current = startOfNextDay(current)
		}
	}

	return weekdayHours / 24
}

func isWeekday(t time.Time) bool {
	day := t.Weekday()
	return day != time.Saturday && day != time.Sunday
}

func startOfNextDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
}
