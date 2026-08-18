# Bug 复现说明

## 问题

按日期汇总签到时，同一宠物名被错误地作为全局唯一键；同一宠物在不同年份的合法签到会被后续去重逻辑跳过。日期分桶键还缺少年份，使跨年同月同日的签到落入同一个报告桶。

## 触发条件

准备同一宠物在 2026-08-16 和 2027-08-16 各有一条签到记录。两条记录都应被接受，并分别参与对应年份的汇总；含 Bug 的实现会在第一条记录之后跳过第二条记录。

## 复现命令

```bash
go test ./internal/careflow -run 'TestSummarizeKeepsSamePetCheckinsOnSeparateYears' -count=1
```

修复前实际输出：

```text
--- FAIL: TestSummarizeKeepsSamePetCheckinsOnSeparateYears (0.00s)
    careflow_test.go:12: two valid checkins should not collapse: careflow.Report{Days:map[string]int{"08-16":1}}
FAIL
FAIL	pet/internal/careflow	0.465s
FAIL
```

失败含义：两条跨年份的合法签到被压缩成一条，报告实际只有 `Days["08-16"] == 1`。重复运行 `-count=20` 可稳定复现该失败。
