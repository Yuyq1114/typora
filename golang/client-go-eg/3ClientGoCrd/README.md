# 第三阶段：client-go 操作 CRD

目标：理解 client-go 默认没有 WebApp 的 typed client，需要通过 Dynamic Client 和 GVR 操作 CRD。

现有代码：

- `client.go`：使用 `dynamic.Interface` Get/List WebApp，返回 `Unstructured`。
- `../main.go`：将 `controllerChoice` 改成 `3` 后运行一次并退出。

关键入口：`Run()`。
