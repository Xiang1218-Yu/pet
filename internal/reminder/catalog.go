package reminder

import "time"

func activeOn(rule Rule, day time.Time) bool {
	for _, weekday := range rule.Weekdays {
		if weekday == day.Weekday() {
			return true
		}
	}
	return false
}
