package reminder

import (
	"testing"
	"time"
)

func TestBuildKeepsDailyLocalCareRemindersAcrossMidnight(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*60*60)
	start := time.Date(2026, 8, 16, 23, 50, 0, 0, loc)
	got := Build([]Rule{{Name: "feed", Hour: 0, Minute: 5}, {Name: "water", Hour: 0, Minute: 5, Weekdays: []time.Weekday{time.Monday}}}, start, 1)
	if len(got) != 2 {
		t.Fatalf("got %d reminders, want 2: %#v", len(got), got)
	}
	if got[0].At.Location() != loc || got[0].At.Day() != 17 {
		t.Fatalf("first reminder should be 2026-08-17 in local zone: %#v", got[0])
	}
	if got[1].Name != "water" || got[1].At.Day() != 17 {
		t.Fatalf("weekday reminder lost or wrong: %#v", got)
	}
	next := Build([]Rule{{Name: "feed", Hour: 0, Minute: 5}}, start, 2)
	if len(next) != 2 {
		t.Fatalf("same daily rule on separate dates was deduplicated: %#v", next)
	}
}
