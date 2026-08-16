package reminder

import "sort"

func normalize(items []Reminder) []Reminder {
	// Dedup on name + local calendar date, not on name alone: a daily rule
	// legitimately fires on consecutive days (e.g. "feed" each morning), and
	// collapsing by name would drop tomorrow's feed as a dup of today's.
	// Date alone (not the full instant) is used because two distinct rules
	// may legitimately share a wall time on the same day — e.g. feed and
	// water both at 00:05 — and must both be kept.
	seen := map[string]bool{}
	out := make([]Reminder, 0, len(items))
	for _, item := range items {
		key := item.Name + "|" + item.At.Format("2006-01-02")
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
	}
	// Stable: ties on At keep their input order, so feed (emitted before
	// water) stays first even when both land on the same instant.
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
