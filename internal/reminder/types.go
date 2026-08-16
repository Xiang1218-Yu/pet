package reminder

import "time"

type Rule struct {
	Name         string
	Hour, Minute int
	Weekdays     []time.Weekday
}
type Reminder struct {
	Name string
	At   time.Time
}
