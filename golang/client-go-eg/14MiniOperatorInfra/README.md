# 扩展四：Mini Operator、网络、存储与 Linux

状态：已实现最小实验版本。

这部分把前面的 Controller 放到完整的 Kubernetes 运行链路中理解。

## Mini Operator

`operator.go` 增加 Lease Leader Election；`Dockerfile` 和 `operator.yaml` 提供镜像、ServiceAccount、RBAC 和双副本 Deployment。Status、Conditions 和 Finalizer继续放在第十阶段实现。

```powershell
docker build -f .\14MiniOperatorInfra\Dockerfile -t client-go-eg:local .
kubectl apply -f .\2WebAppCrd\crd.yaml
kubectl apply -f .\14MiniOperatorInfra\operator.yaml
kubectl logs -n mini-operator-system deployment/mini-operator
```

如果集群不共享 Docker 本地镜像，需要先将镜像导入当前集群的 containerd 或推送到镜像仓库。

## 网络

```text
Pod → kubelet → CRI → Container Runtime → CNI
    → veth → network namespace → route / iptables / eBPF
```

实验文件：`network-pod.yaml`。

## 存储

```text
PVC → PV → CSI Controller → Attach
    → CSI Node → Mount → Pod
```

实验文件：`storage.yaml`。

## Linux 资源隔离

```text
Kubernetes resources
        ↓
Container Runtime
        ↓
runc
        ↓
cgroup / namespace / mount / capabilities / seccomp
```

实验文件：`linux-pod.yaml`，包含资源限制、seccomp、禁止提权和删除 capabilities。
