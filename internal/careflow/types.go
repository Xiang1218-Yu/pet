package careflow

import "time"

type Dose struct {
	Medicine string
	At       time.Time
	Units    int
}
type Timeline struct{ Doses []Dose }
