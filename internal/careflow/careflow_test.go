package careflow

import (
	"testing"
	"time"
)

func TestBuildKeepsRepeatedCareAcrossDays(t *testing.T) {
	loc := time.FixedZone("UTC+8", 28800)
	from := time.Date(2026, 8, 16, 0, 0, 0, 0, loc)
	got := Build([]Visit{{"feed", time.Date(2026, 8, 16, 8, 0, 0, 0, loc)}, {"feed", time.Date(2026, 8, 17, 8, 0, 0, 0, loc)}, {"clean", time.Date(2026, 8, 17, 9, 0, 0, 0, loc)}}, from)
	if len(got) != 2 || len(got[0].Visits) != 1 || len(got[1].Visits) != 2 {
		t.Fatalf("daily care queue merged distinct visits: %#v", got)
	}
}
