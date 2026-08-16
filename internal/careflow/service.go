package careflow

func Summarize(in []Checkin) Report {
	r := Report{Days: map[string]int{}}
	seen := map[string]bool{}
	for _, c := range in {
		k := c.identity()
		if !accepted(c) || seen[k] {
			continue
		}
		seen[k] = true
		r.Days[dayKey(c.At)]++
	}
	return r
}
