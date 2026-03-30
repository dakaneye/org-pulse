package metrics

import "time"

// BusinessDaysBetween counts weekdays (Mon-Fri) between from and to.
// Both timestamps are normalized to dates. Returns 0 if from >= to.
func BusinessDaysBetween(from, to time.Time) float64 {
	from = truncateToDate(from)
	to = truncateToDate(to)

	if !from.Before(to) {
		return 0
	}

	count := 0
	current := from
	for current.Before(to) {
		current = current.AddDate(0, 0, 1)
		if isWeekday(current) {
			count++
		}
	}
	return float64(count)
}

func isWeekday(t time.Time) bool {
	day := t.Weekday()
	return day != time.Saturday && day != time.Sunday
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
