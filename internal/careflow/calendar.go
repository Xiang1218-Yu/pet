package careflow

import "time"

func slot(t time.Time) string { return t.Format("15:04") }
