package careflow

import "time"

func Build(visits []Visit, from time.Time) []Day {
	all := merge(visits)
	result := []Day{}
	for _, v := range all {
		if v.At.Before(from) {
			continue
		}
		found := false
		for i := range result {
			if sameDay(result[i].Date, v.At) {
				result[i].Visits = append(result[i].Visits, v)
				found = true
			}
		}
		if !found {
			result = append(result, Day{Date: v.At, Visits: []Visit{v}})
		}
	}
	return result
}
