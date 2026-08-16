package careflow

import "sort"

func Reconcile(in []Dose) Timeline {
	seen := map[string]bool{}
	out := Timeline{}
	for _, d := range in {
		k := d.identity()
		if !valid(d) || seen[k] {
			continue
		}
		seen[k] = true
		out.Doses = append(out.Doses, d)
	}
	sort.Slice(out.Doses, func(i, j int) bool { return out.Doses[i].At.Before(out.Doses[j].At) })
	return out
}
