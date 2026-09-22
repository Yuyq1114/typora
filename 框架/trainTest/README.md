# Kubeflow 单卡训练示例

这个项目用 Kubeflow Trainer 在本机 k3s 上训练一个最简单的线性回归模型，并申请一张 NVIDIA GPU 来演示 GPU 调度。

为避免首次下载 3GB 以上的 PyTorch 镜像，训练算法只使用 Python 标准库在 CPU 上计算；训练 Pod 会实际占用 GPU，并用 `nvidia-smi` 验证 GPU 已注入。等网络条件合适时，再把镜像和代码替换为 PyTorch，就能让训练计算本身运行在 GPU 上。

## 文件作用

- `train.py`：实际执行的线性回归训练代码，并验证 Pod 内的 GPU。
- `trainjob.yaml`：告诉 Kubeflow 使用哪个运行环境、镜像和 GPU 数量。
- `run.ps1`：更新代码、提交任务、等待完成并显示日志。
- `cleanup.ps1`：删除示例任务和代码 ConfigMap。

本地代码会被保存到 Kubernetes ConfigMap，再挂载到训练 Pod 的 `/workspace/train.py`。因此修改 `train.py` 后不需要重新制作容器镜像。

## 运行

在 PowerShell 中进入本目录，然后执行：

```powershell
.\run.ps1
```

成功时日志最后应接近：

```text
result: y = 3.0000x + 2.0000
training completed
```

查看 Kubeflow 生成的资源：

```powershell
kubectl get trainjob,jobset,pods
kubectl describe trainjob kubeflow-train-demo
kubectl logs -l jobset.sigs.k8s.io/jobset-name=kubeflow-train-demo --all-containers=true
```

修改 `train.py` 后再次运行 `.\run.ps1`，脚本会删除旧任务并用新代码重新训练。

清理示例：

```powershell
.\cleanup.ps1
```

## 执行链路

```text
run.ps1
  -> ConfigMap（保存 train.py）
  -> TrainJob（Kubeflow API）
  -> JobSet
  -> Pod
  -> NVIDIA GPU
```
