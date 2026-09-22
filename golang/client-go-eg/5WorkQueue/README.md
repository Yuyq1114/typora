# 第五阶段：WorkQueue

本阶段代码已经包含在第二阶段。

位置：`../2WebAppCrd/controller.go`

- `queue` 字段使用 `TypedRateLimitingInterface[string]`。
- `NewTypedRateLimitingQueue()` 创建队列。
- EventHandler 将 `namespace/name` 加入队列。
- `Run()` 循环调用 `queue.Get()`。
- 失败时调用 `AddRateLimited()` 重试，成功时调用 `Forget()`。
- 每个任务完成后调用 `Done()`。
