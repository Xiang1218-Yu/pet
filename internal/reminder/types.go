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

func (r Rule) at(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), r.Hour, r.Minute, 0, 0, day.Location())
}
