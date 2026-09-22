# 第六阶段：Lister / Cache

本阶段代码已经包含在第二阶段。

位置：`../2WebAppCrd/controller.go`

- `webInformer.Lister()`：读取 WebApp 缓存。
- `depInformer.Lister()`：读取 Deployment 缓存。
- `svcInformer.Lister()`：读取 Service 缓存。
- `reconcile()` 中通过这些 Lister 的 `Get()` 读取资源。

Lister 读取 Informer 的本地缓存；创建和更新资源时才通过 client-go 请求 API Server。
