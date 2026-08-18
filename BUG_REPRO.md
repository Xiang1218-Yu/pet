# Bug 复现说明

## 问题

跨午夜构建本地护理提醒时，默认每日规则可能被错误过滤，且跨日期的同名提醒可能被当成同一条记录合并，导致应保留的次日提醒缺失。

## 触发条件

使用固定 UTC+8 时区和 23:50 的启动时间，输入次日 00:05 的 feed/water 规则，并保留连续两天同名的 feed 规则。

## 复现命令

```bash
go test ./internal/reminder -run 'TestBuildKeepsDailyLocalCareRemindersAcrossMidnight' -count=1
```

修复前实际输出：

```text
--- FAIL: TestBuildKeepsDailyLocalCareRemindersAcrossMidnight (0.00s)
    reminder_test.go:13: got 1 reminders, want 2: []reminder.Reminder{reminder.Reminder{Name:"water", At:time.Date(2026, time.August, 17, 0, 5, 0, 0, time.UTC)}}
FAIL
FAIL	pet/internal/reminder	0.801s
FAIL
```

失败含义：期望保留 2 条提醒，但实际只得到 1 条；次日 00:05 的另一条提醒被过滤或在跨日期同名去重时丢失。重复运行 `-count=20` 可稳定复现该失败。
