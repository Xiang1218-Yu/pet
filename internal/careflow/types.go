package careflow

import "time"

type Dose struct {
	Medicine string
	At       time.Time
	Units    int
}
type Timeline struct{ Doses []Dose }

func (d Dose) identity() string { return d.Medicine + "@" + key(d.At) }
