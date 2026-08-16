package careflow

import "time"

// key identifies a single treatment slot per calendar day so that a dose
// scheduled at the same wall-clock time on different days (e.g. a morning
// vitamin each day) is kept as separate entries rather than collapsed.
func key(t time.Time) string { return t.Format("2006-01-02 15:04") }
