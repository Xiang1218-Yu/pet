package reminder

import "time"

func Build(rules []Rule, start time.Time, days int) []Reminder {
	var result []Reminder
	for offset := 0; offset <= days; offset++ {
		day := start.AddDate(0, 0, offset)
		for _, rule := range rules {
			if !activeOn(rule, day) {
				continue
			}
			at := time.Date(day.Year(), day.Month(), day.Day(), rule.Hour, rule.Minute, 0, 0, time.UTC)
			if at.Before(start) {
				continue
			}
			result = append(result, Reminder{Name: rule.Name, At: at})
		}
	}
	return normalize(result)
}
