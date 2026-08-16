package careflow

import (
	"testing"
	"time"
)

func TestReconcilePreservesSeparateTreatmentDays(t *testing.T) {
	loc := time.FixedZone("UTC+8", 28800)
	got := Reconcile([]Dose{{"vitamin", time.Date(2026, 8, 16, 8, 0, 0, 0, loc), 1}, {"vitamin", time.Date(2026, 8, 17, 8, 0, 0, 0, loc), 1}, {"medicine", time.Date(2026, 8, 17, 9, 0, 0, 0, loc), 0}})
	if len(got.Doses) != 2 {
		t.Fatalf("expected two valid daily doses, got %#v", got)
	}
}
