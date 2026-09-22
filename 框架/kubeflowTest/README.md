# Kubeflow Trainer 开发工作流

这个项目不讲模型原理，而是用同一份训练代码逐层理解 Kubeflow Trainer：

```text
训练代码
  ↓
Kubernetes Job
  ↓
Kubeflow TrainJob
  ↓
自定义 TrainingRuntime
  ↓
Go/client-go 提交 TrainJob
```

当前项目学习的是“训练任务开发工作流”，不是 Kubeflow Pipelines。

项目使用预加载的 `python:3.11-slim`，Pod 启动时不在线安装依赖，因此各阶段可以稳定重复运行。训练代码故意只用标准库；你可以在理解控制链路后把镜像替换成自己的 PyTorch CUDA 镜像。

## 先记住三个角色

| 角色 | 负责的内容 |
| --- | --- |
| 算法开发者 | 镜像、命令、训练参数、资源数量 |
| 平台开发者 | TrainingRuntime、调度、存储、网络、默认策略 |
| Kubeflow Controller | 把 TrainJob 与 Runtime 合成为 JobSet，并同步状态 |

## Stage 1：原生 Kubernetes Job

文件：`manifests/01-k8s-job.yaml`

```powershell
.\run.ps1 -Stage 1
```

学习目标：先确认没有 Kubeflow 时，需要自己填写完整 Pod、GPU、挂载和重试策略。

观察：

```powershell
kubectl get job stage1-k8s-job -o yaml
kubectl get pods -l job-name=stage1-k8s-job
```

## Stage 2：使用 TrainJob

文件：`manifests/02-trainjob.yaml`

```powershell
.\run.ps1 -Stage 2
```

学习目标：TrainJob 只描述“我要怎样训练”；Kubeflow 把它和 `torch-distributed` Runtime 合成为 JobSet。

观察控制链：

```powershell
kubectl get trainjob stage2-trainjob -o yaml
kubectl get jobset stage2-trainjob -o yaml
kubectl get jobs,pods -l jobset.sigs.k8s.io/jobset-name=stage2-trainjob
```

## Stage 3：平台开发者自定义 Runtime

文件：

- `manifests/03-runtime.yaml`
- `manifests/03-trainjob.yaml`

```powershell
.\run.ps1 -Stage 3
```

学习目标：理解 Runtime 为什么是 Kubeflow Trainer 的扩展边界。平台团队把 NVIDIA RuntimeClass、JobSet 模板等能力放进 Runtime，算法开发者只提交 TrainJob。

比较：

```powershell
kubectl get trainingruntime single-gpu-runtime -o yaml
kubectl get trainjob stage3-custom-runtime -o yaml
kubectl get jobset stage3-custom-runtime -o yaml
```

## Stage 4：使用 Go 提交和观察任务

文件：`go-client/main.go`

```powershell
.\run.ps1 -Stage 4
```

Go 程序演示：

1. 使用当前 kubeconfig 创建 dynamic client。
2. 构造并创建 TrainJob CR。
3. 轮询 `status.conditions`。
4. 根据 JobSet 标签找到训练 Pod。
5. 读取训练日志。

这里使用 dynamic client，是因为 Kubeflow Trainer 的 Go API 没有放进 Kubernetes 官方 client-go。

## 统一检查与清理

```powershell
.\inspect.ps1
.\cleanup.ps1
```

## 每完成一阶段应该回答的问题

1. 当前提交的顶层资源是什么？
2. 谁创建了下一层资源？
3. PodSpec 最终由哪些对象共同决定？
4. Pod 完成后，状态怎样回到 TrainJob？
5. 哪些字段应该由算法开发者填写，哪些应该由平台开发者控制？

## 自定义 Runtime 的两个必要约定

本项目的 Runtime 包含两个容易忽略但非常重要的标签：

```yaml
metadata:
  labels:
    trainer.kubeflow.org/framework: torch

template:
  metadata:
    labels:
      trainer.kubeflow.org/trainjob-ancestor-step: trainer
```

- `framework: torch` 选择 Torch MLPolicy 插件。
- `trainjob-ancestor-step: trainer` 标出哪个 ReplicatedJob 接收 TrainJob 中的 image、command、env 和 resources。

漏掉第二个标签时，JobSet 仍会创建并显示成功，但用户提交的训练命令不会进入 Pod。这是阅读和二次开发 Runtime 时必须理解的契约。

## 下一轮扩展

完成四个阶段后，再按这个顺序扩展：

1. 把 ConfigMap 代码改成自己构建的训练镜像。
2. 给 Runtime 增加 PVC，保存模型 checkpoint。
3. 给 TrainJob 增加 suspend、deadline 和失败重试实验。
4. 阅读 Kubeflow Trainer 的 TrainJob Reconciler。
5. 写一个插件，在 Build 阶段自动向 Pod 注入标签或环境变量。
