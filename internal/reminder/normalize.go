package reminder

import "sort"

func normalize(items []Reminder) []Reminder {
	seen := map[string]bool{}
	out := make([]Reminder, 0, len(items))
	for _, item := range items {
		key := item.Name + "@" + item.At.Format("2006-01-02T15:04:05Z07:00")
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
