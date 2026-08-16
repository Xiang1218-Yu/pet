package reminder

import "time"

func Build(rules []Rule, start time.Time, days int) []Reminder {
	var result []Reminder
	if days < 0 {
		return result
	}
	for offset := 0; offset <= days; offset++ {
		day := start.AddDate(0, 0, offset)
		for _, rule := range rules {
			if !activeOn(rule, day) {
				continue
			}
			at := rule.at(day)
			if at.Before(start) {
				continue
			}
			result = append(result, Reminder{Name: rule.Name, At: at})
		}
	}
	return normalize(result)
}
