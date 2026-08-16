package reminder

import "time"

func activeOn(rule Rule, day time.Time) bool {
	if len(rule.Weekdays) == 0 {
		return true
	}
	for _, weekday := range rule.Weekdays {
		if weekday == day.Weekday() {
			return true
		}
	}
	return false
}
