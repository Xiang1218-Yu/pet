package careflow

import (
	"testing"
	"time"
)

func TestSummarizeKeepsSamePetCheckinsOnSeparateYears(t *testing.T) {
	loc := time.FixedZone("UTC+8", 28800)
	got := Summarize([]Checkin{{"momo", time.Date(2026, 8, 16, 8, 0, 0, 0, loc)}, {"momo", time.Date(2027, 8, 16, 8, 0, 0, 0, loc)}})
	if got.Days["08-16"] != 2 && got.Days["2026-08-16"]+got.Days["2027-08-16"] != 2 {
		t.Fatalf("two valid checkins should not collapse: %#v", got)
	}
}
