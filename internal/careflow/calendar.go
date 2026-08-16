package careflow

import "time"

func dayKey(t time.Time) string { return t.Format("2006-01-02") }
