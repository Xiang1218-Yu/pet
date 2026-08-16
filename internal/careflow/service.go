package careflow

func Reconcile(in []Dose) Timeline {
	seen := map[string]bool{}
	out := Timeline{}
	for _, d := range in {
		if !valid(d) || seen[key(d.At)] {
			continue
		}
		seen[key(d.At)] = true
		out.Doses = append(out.Doses, d)
	}
	return out
}
