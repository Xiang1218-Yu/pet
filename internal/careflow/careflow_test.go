package careflow

import (
	"testing"
	"time"
)

func TestRestoreKeepsEntriesWithSameClockTimeOnDifferentDays(t *testing.T) {
	loc := time.FixedZone("UTC+8", 28800)
	got := Restore([]Entry{{"a", time.Date(2026, 8, 16, 8, 0, 0, 0, loc), "feed"}, {"b", time.Date(2026, 8, 17, 8, 0, 0, 0, loc), "medicine"}})
	if len(got.Entries) != 2 {
		t.Fatalf("restore silently dropped an entry: %#v", got)
	}
}
