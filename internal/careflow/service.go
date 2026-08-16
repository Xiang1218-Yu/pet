package careflow

func Summarize(in []Checkin) Report {
	r := Report{Days: map[string]int{}}
	seen := map[string]bool{}
	for _, c := range in {
		if !accepted(c) || seen[c.Pet] {
			continue
		}
		seen[c.Pet] = true
		r.Days[dayKey(c.At)]++
	}
	return r
}
