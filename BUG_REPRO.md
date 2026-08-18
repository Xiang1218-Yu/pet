# Bug 复现说明

## 问题

重建治疗时间线时，去重键只保留时分，导致不同日期但同为 08:00 的有效剂量发生碰撞；同时，零剂量记录被错误视为有效记录。

## 触发条件

准备两条有效的维生素剂量记录：日期分别为 2026-08-16 和 2026-08-17，时间都为 08:00；再加入一条 `Units == 0` 的药物记录。正确结果应保留两条有效维生素剂量，并拒绝零剂量记录。

## 复现命令

公开目标测试：

```bash
go test ./internal/careflow -run 'TestReconcilePreservesSeparateTreatmentDays' -count=1
```

修复前公开测试实际输出：

```text
ok  	pet/internal/careflow	0.463s
```

该测试会出现**假通过**：它只断言最终数量为 2，而 Bug 实现实际返回第一天的 vitamin 和第二天的零剂量 medicine，第二天的有效 vitamin 已被丢弃。

用同一输入增加内容校验后的诊断输出为：

```text
--- FAIL: TestBugReproTemporary (0.00s)
    bug_repro_tmp_test.go:14: reconcile mismatch: got 2 doses: careflow.Timeline{Doses:[]careflow.Dose{careflow.Dose{Medicine:"vitamin", At:time.Date(2026, time.August, 16, 8, 0, 0, 0, time.Location("UTC+8")), Units:1}, careflow.Dose{Medicine:"medicine", At:time.Date(2026, time.August, 17, 9, 0, 0, 0, time.Location("UTC+8")), Units:0}}}
FAIL
FAIL	pet/internal/careflow	0.762s
FAIL
```

因此该题的错误现象是“数量恰好为 2 但内容错误”：第二天有效剂量丢失，零剂量记录错误进入结果。
