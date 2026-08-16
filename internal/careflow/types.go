package careflow

import "time"

type Checkin struct {
	Pet string
	At  time.Time
}
type Report struct{ Days map[string]int }

func (c Checkin) identity() string { return c.Pet + "@" + c.At.Format("2006-01-02T15:04Z07:00") }
