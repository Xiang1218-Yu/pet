# Bug reproduction

Reconcile two positive vitamin doses at 08:00 on consecutive days plus a zero-unit medicine entry. The valid second dose is lost while the invalid entry can be accepted. Run `go test ./internal/careflow -run TestReconcilePreservesSeparateTreatmentDays -count=20` to reproduce it.
