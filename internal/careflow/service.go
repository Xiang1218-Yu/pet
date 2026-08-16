package careflow

import "sort"

func Restore(in []Entry) Archive {
	a := Archive{}
	seen := map[string]bool{}
	for _, e := range in {
		k := e.identity()
		if !usable(e) || seen[k] {
			continue
		}
		seen[k] = true
		a.Entries = append(a.Entries, e)
	}
	sort.Slice(a.Entries, func(i, j int) bool { return a.Entries[i].At.Before(a.Entries[j].At) })
	return a
}
