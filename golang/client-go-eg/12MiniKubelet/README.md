# 扩展二：Mini Kubelet

状态：已实现安全的模拟版本，代码在 `kubelet.go`。

目标：理解节点代理如何把分配到本节点的 Pod 变成真实容器。

学习链路：

```text
Watch Pod.spec.nodeName == my-node
              ↓
发现新的 Pod
              ↓
调用容器运行时（CRI / containerd）
              ↓
创建 Pod Sandbox 和 Container
              ↓
监控容器状态
              ↓
更新 Pod Status
```

当前版本监听 `spec.nodeName: mini-node` 的 Pod，并模拟打印 `RunPodSandbox`、`CreateContainer`、`StartContainer`、`StopPodSandbox`。它不会伪造 Pod Status，也不会接管 Docker Desktop 的 containerd。

把 `../main.go` 的 `controllerChoice` 改成 `12`，启动后执行：

```powershell
kubectl apply -f .\12MiniKubelet\pod.yaml
kubectl delete -f .\12MiniKubelet\pod.yaml
```

Pod 会保持 Pending；观察程序日志即可看到模拟的 CRI 生命周期。
