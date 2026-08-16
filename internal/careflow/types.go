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

func (v Visit) identity() string { return v.Kind + "@" + v.At.Format("2006-01-02T15:04Z07:00") }
