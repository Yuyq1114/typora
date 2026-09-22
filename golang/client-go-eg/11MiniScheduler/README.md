# 扩展一：Mini Scheduler

状态：已实现最小可运行版本，代码在 `scheduler.go`。

目标：理解 Scheduler 如何为 `spec.nodeName` 为空的 Pod 选择节点并执行 Bind。

最小版本：

1. Watch 未调度 Pod。
2. List Node。
3. Filter：检查 CPU、Memory、NodeSelector、Taint/Toleration。
4. Score：按剩余资源为候选节点打分。
5. Bind：调用 Pod Binding API，把 Pod 绑定到得分最高的节点。

对应 Scheduling Framework：

```text
PreFilter → Filter → PostFilter
PreScore → Score
Reserve → Permit → PreBind → Bind → PostBind
```

当前实现 `Filter → Score → Bind`。它只处理 `schedulerName: mini-scheduler` 的 Pod，不影响默认 Scheduler。

把 `../main.go` 的 `controllerChoice` 改成 `11`，启动后执行：

```powershell
kubectl apply -f .\11MiniScheduler\pod.yaml
kubectl get pod mini-scheduler-demo -n test1 -o wide
```
