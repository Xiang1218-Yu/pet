# Bug reproduction

Create feed visits on two consecutive days and a clean visit on the second day. The queue silently retains only one feed visit. Run `go test ./internal/careflow -run TestBuildKeepsRepeatedCareAcrossDays -count=20` to reproduce it.
