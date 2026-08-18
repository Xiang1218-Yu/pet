# Bug 复现说明

## 问题

恢复归档记录时，系统使用只包含时分的键进行去重。不同日期但同为 08:00 的两条有效记录会被误判为重复；后到记录直接被静默跳过，没有错误、日志或丢弃计数。

## 触发条件

准备两条不同标识、不同日期但同为 08:00 的归档记录：记录 a 的时间为 2026-08-16 08:00，记录 b 的时间为 2026-08-17 08:00。两条记录都应恢复，但含 Bug 的实现会让输入顺序中的第一条胜出。

## 复现命令

```bash
go test ./internal/careflow -run 'TestRestoreKeepsEntriesWithSameClockTimeOnDifferentDays' -count=1
```

修复前实际输出：

```text
--- FAIL: TestRestoreKeepsEntriesWithSameClockTimeOnDifferentDays (0.00s)
    careflow_test.go:12: restore silently dropped an entry: careflow.Archive{Entries:[]careflow.Entry{careflow.Entry{ID:"a", At:time.Date(2026, time.August, 16, 8, 0, 0, 0, time.Location("UTC+8")), Note:"feed"}}}
FAIL
FAIL	pet/internal/careflow	0.461s
FAIL
```

失败含义：输入包含 a、b 两条记录，但恢复结果只保留 a；记录 b 因 `slot("15:04")` 只生成 `08:00` 键而被无声丢弃。重复运行 `-count=20` 可稳定复现该失败。
