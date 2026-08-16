# Bug reproduction

Summarize check-ins for the same pet on August 16 in 2026 and 2027. Only one check-in appears in the report. Run `go test ./internal/careflow -run TestSummarizeKeepsSamePetCheckinsOnSeparateYears -count=20` to reproduce it.
