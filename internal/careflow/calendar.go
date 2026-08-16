package careflow

import "time"

func sameDay(a, b time.Time) bool { return a.YearDay() == b.YearDay() }
