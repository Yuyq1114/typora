# 第九阶段：OwnerReference

本阶段代码已经包含在第二阶段。

位置：`../2WebAppCrd/controller.go`

- `ownerReference()`：为 Deployment 和 Service 设置 WebApp 为 controller owner。
- `ownedBy()`：更新资源前确认它属于当前 WebApp。
- `enqueueOwner()`：子资源变化时找到所属 WebApp 并触发 Reconcile。

删除 WebApp 后，Kubernetes Garbage Collector 会删除它拥有的 Deployment 和 Service：

```powershell
kubectl delete webapp webapp-demo -n test1
kubectl get deployment,service,pod -n test1
```
