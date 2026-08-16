package careflow

import "time"

type Checkin struct {
	Pet string
	At  time.Time
}
type Report struct{ Days map[string]int }
