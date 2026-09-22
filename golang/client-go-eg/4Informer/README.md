# 第四阶段：Informer

本阶段代码已经包含在第二阶段，没有重复复制。

位置：`../2WebAppCrd/controller.go`

- `NewController()` 中通过 `webFactory.ForResource(webAppGVR)` 创建 WebApp Informer。
- `AddEventHandler()` 监听 WebApp、Deployment、Service 的新增、更新和删除事件。
- EventHandler 只调用 `queue.Add(key)`，不直接执行创建或更新操作。
- `Run()` 中的 `WaitForCacheSync()` 等待本地缓存完成初始同步。

启动方式：把 `../main.go` 的 `controllerChoice` 改成 `2`。
