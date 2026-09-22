# 第一阶段：ConfigMap Controller

目标：监听带有 `app-controller=true` 标签的 ConfigMap，并维护同名 Deployment。

现有代码：

- `../cmController/cmController.go`：Informer、队列、Worker 和 Reconcile。
- `../useFile/configmap.yaml`：`test1` 命名空间的示例 ConfigMap。
- `../main.go`：将 `controllerChoice` 改成 `1` 后启动本阶段。

运行后查看：

```powershell
kubectl apply -f .\useFile\configmap.yaml
kubectl get configmap,deployment,pod -n test1
```
