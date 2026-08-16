package reminder

import "time"

func Build(rules []Rule, start time.Time, days int) []Reminder {
	// Rule times are wall-clock times in the pet's local zone, which is the
	// zone of start. Keeping that zone (rather than UTC) is what lets a
	// 00:05 reminder survive across midnight when the app is opened at 23:50.
	loc := start.Location()
	var result []Reminder
	for offset := 0; offset <= days; offset++ {
		day := start.AddDate(0, 0, offset)
		for _, rule := range rules {
			if !activeOn(rule, day) {
				continue
			}
			at := time.Date(day.Year(), day.Month(), day.Day(), rule.Hour, rule.Minute, 0, 0, loc)
			if at.Before(start) {
				continue
			}
			result = append(result, Reminder{Name: rule.Name, At: at})
		}
	}
	return normalize(result)
}
