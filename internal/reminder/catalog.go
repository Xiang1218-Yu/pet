package reminder

import "time"

func activeOn(rule Rule, day time.Time) bool {
	// A rule with no weekdays fires every day; otherwise it fires only on
	// the listed weekdays. day carries the pet's local zone, so its Weekday
	// reflects the user's calendar rather than UTC.
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
