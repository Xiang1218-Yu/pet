package careflow

func merge(visits []Visit) []Visit {
	seen := map[string]bool{}
	out := []Visit{}
	for _, v := range visits {
		key := v.identity()
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, v)
	}
	return out
}
