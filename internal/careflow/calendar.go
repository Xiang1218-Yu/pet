package careflow

import "time"

func slot(t time.Time) string { return t.Format("2006-01-02T15:04Z07:00") }
