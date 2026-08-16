package careflow

import "time"

func key(t time.Time) string { return t.Format("15:04") }
