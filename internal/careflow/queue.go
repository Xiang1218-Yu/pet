package careflow

func merge(visits []Visit) []Visit {
	seen := map[string]bool{}
	out := []Visit{}
	for _, v := range visits {
		if seen[v.Kind] {
			continue
		}
		seen[v.Kind] = true
		out = append(out, v)
	}
	return out
}
