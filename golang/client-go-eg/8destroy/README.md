# 第八阶段：故意破坏

本阶段不需要新增 controller 代码，用实验验证 Reconcile 是否持续维护期望状态。

1. 把 `../main.go` 的 `controllerChoice` 改成 `2` 并启动。
2. 手动修改副本数：

```powershell
kubectl scale deployment webapp-demo -n test1 --replicas=3
kubectl get deployment webapp-demo -n test1
```

WebApp 声明的是 1，controller 应恢复为 1。

3. 手动删除 Deployment：

```powershell
kubectl delete deployment webapp-demo -n test1
kubectl get deployment webapp-demo -n test1
```

controller 应重新创建它。对应代码位于 `../2WebAppCrd/controller.go` 的 Deployment EventHandler 和 `ensureDeployment()`。
