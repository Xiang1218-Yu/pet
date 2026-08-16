# Bug reproduction

Build a daily reminder list from 23:50 in UTC+8 with the next-day feed and water rules. The default daily rule is filtered out, and same-name reminders on different dates can be collapsed. Run `go test ./internal/reminder -run TestBuildKeepsDailyLocalCareRemindersAcrossMidnight -count=20` to reproduce the stable failure.
