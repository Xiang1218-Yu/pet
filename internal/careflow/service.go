package careflow

func Restore(in []Entry) Archive {
	a := Archive{}
	seen := map[string]bool{}
	for _, e := range in {
		if !usable(e) || seen[slot(e.At)] {
			continue
		}
		seen[slot(e.At)] = true
		a.Entries = append(a.Entries, e)
	}
	return a
}
