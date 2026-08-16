package careflow

import "time"

type Visit struct {
	Kind string
	At   time.Time
}
type Day struct {
	Date   time.Time
	Visits []Visit
}
