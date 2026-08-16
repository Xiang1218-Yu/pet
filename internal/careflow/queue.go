package careflow

func merge(visits []Visit) []Visit {
	seen := map[string]bool{}
	out := []Visit{}
	for _, v := range visits {
		// Dedupe by kind AND calendar day, so the same kind of care on
		// consecutive days (e.g. feed on day 1 and feed on day 2) is kept.
		key := v.Kind + v.At.Format("2006-01-02")
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, v)
	}
	return out
}
