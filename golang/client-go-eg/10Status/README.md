# 第十阶段：Status

本阶段尚未实现。

目标：controller 读取 Deployment 的实际状态，并写入 WebApp：

```yaml
status:
  readyReplicas: 1
  phase: Running
```

需要修改：

1. `../2WebAppCrd/crd.yaml`：增加 `status` schema 和 `subresources.status`。
2. `../2WebAppCrd/controller.go`：保存 Dynamic Client，在 Reconcile 末尾更新 WebApp 的 `/status`。
3. Deployment 更新事件继续触发所属 WebApp 的 Reconcile。

完成后查看：

```powershell
kubectl get webapp webapp-demo -n test1 -o yaml
```

`spec` 表示期望状态，`status` 表示 controller 观察到的实际状态。
