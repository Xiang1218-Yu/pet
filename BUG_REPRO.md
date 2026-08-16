# Bug reproduction

Restore two archive entries that have different identifiers and dates but both occur at 08:00. The second entry is silently discarded. Run `go test ./internal/careflow -run TestRestoreKeepsEntriesWithSameClockTimeOnDifferentDays -count=20` to reproduce it.
