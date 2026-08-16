package careflow

import "time"

type Entry struct {
	ID   string
	At   time.Time
	Note string
}
type Archive struct{ Entries []Entry }
