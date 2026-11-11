# 1、介绍 AI 模型训练和推理的基本流程

标准答案：

训练的目标是通过数据优化模型参数，推理则是使用训练完成的参数计算预测结果。

训练流程通常是：

```
数据加载 → 前向传播 → Loss 计算 → 反向传播 → 梯度更新 → 保存 Checkpoint
```

推理流程通常是：

```
模型加载 → 输入预处理 → 前向计算 → 输出后处理
```

训练一般需要保留反向传播所需的中间激活和梯度，推理通常不需要这些信息，因此两者的显存构成、计算负载和性能优化方向不同。

追问：训练和推理在计算、显存、并发和调度方面有什么区别？

最佳回答：

| 维度 | 训练                                | 推理                         |
| ---- | ----------------------------------- | ---------------------------- |
| 计算 | 前向 + 反向传播 + 参数更新          | 主要是前向计算               |
| 显存 | 参数、梯度、优化器状态、激活        | 参数、部分中间激活、缓存     |
| 目标 | 训练吞吐、收敛效率、资源利用率      | 延迟、吞吐、成本             |
| 并行 | 数据、模型、流水线等并行            | 批处理、模型并行、请求并发   |
| 调度 | 长任务、Gang Scheduling、Checkpoint | 弹性扩缩容、低延迟、负载均衡 |

需要补充一点：推理不一定比训练更简单。大模型推理可能受到 KV Cache、显存带宽和动态批处理的影响。

对于 VLA 模型，推理还可能受到机器人控制频率和动作输出延迟的约束。



# 2、PyTorch 中 CPU Tensor 和 CUDA Tensor 有什么区别？

标准答案：

CPU Tensor 的数据存储和计算主要由 CPU 侧完成；CUDA Tensor 的数据位于 CUDA 设备内存中，相应计算通常由 GPU 执行。

示例：

```
import torchdevice = torch.device(    "cuda" if torch.cuda.is_available() else "cpu")model = MyModel().to(device)model.eval()x = torch.randn(16, 3, 224, 224).to(device)with torch.inference_mode():    output = model(x)
```

`model.to(device)` 将模型参数和相关缓冲区移动到设备；输入 Tensor 也需要位于兼容的设备上。

追问：如何将模型放到 GPU 上推理？数据传输有什么开销？

最佳回答：

模型和输入数据需要在正确的 CUDA 设备上。

CPU 与 GPU 之间的数据传输可能经过 PCIe 等互连，带来延迟和带宽成本。

优化方法包括：

- 避免每层计算前后都进行 CPU/GPU 数据拷贝。
- 尽可能将连续计算保留在 GPU。
- 使用适当的 Batch。
- 对满足条件的数据传输使用 pinned memory 和异步拷贝。
- 避免不必要的同步操作，例如频繁调用 `tensor.item()`。

还需要注意 CUDA 操作通常具有异步执行特征，简单使用 Python 的计时函数可能无法准确测量 GPU Kernel 执行时间。



# 3、常见的模型推理优化技术有哪些？



标准答案：

常见推理优化手段包括：

| 技术          | 核心作用                                   | 主要代价或限制         |
| ------------- | ------------------------------------------ | ---------------------- |
| FP16 / BF16   | 降低精度与内存带宽需求，利用低精度计算能力 | 数值精度、硬件支持     |
| INT8 量化     | 减少权重占用并可能加速计算                 | 精度损失、量化适配     |
| Batching      | 合并请求，提高硬件吞吐                     | 单请求等待延迟可能增加 |
| CUDA Graph    | 降低重复 Kernel Launch 的 CPU 开销         | 图捕获和形状等约束     |
| Kernel Fusion | 合并计算操作，减少内存访问和启动开销       | 实现复杂度             |
| TensorRT      | 进行图优化、算子选择与执行加速             | 模型和算子兼容性       |

追问：FP16、BF16、INT8、Batching 和 CUDA Graph 各自解决什么问题？

最佳回答：

FP16 和 BF16 主要降低计算数据精度与存储成本，两者数值范围和有效精度不同。

BF16 的指数范围接近 FP32，通常更不容易出现数值溢出，但它的尾数精度低于 FP16。

INT8 量化能进一步降低模型权重的存储成本，但不意味着所有模型量化后都会更快。是否加速取决于硬件、算子和具体量化方案。

Batching 主要提高吞吐量，尤其适合较小的推理任务，但如果 Batch 太大，可能增加单请求延迟或显存压力。

CUDA Graph 主要用于减少重复 CUDA 工作负载的 CPU 启动开销，适合执行结构较稳定的计算图。

加分回答：

在仿真评测场景，我还会优先考虑减少模型重复加载、复用 GPU 上的模型实例，并把仿真和推理的性能分别测量，因为整个评测慢不一定是模型推理慢。



# 4、如何分析 PyTorch 程序的性能瓶颈？

标准答案：

我会分四层排查：

第一层：整体运行指标

观察 GPU 利用率、显存、CPU、内存、数据读取速度和每秒处理样本数。

第二层：PyTorch 算子

通过 `torch.profiler` 分析耗时热点、算子调用和 CPU/GPU 活动。

第三层：CUDA 执行

通过 NVIDIA Nsight Systems 查看 CUDA Kernel、内存拷贝、CPU 调用和 GPU 活动时间线。

第四层：底层 Kernel

对于已定位的 GPU Kernel，再使用更细粒度的工具分析访存效率、占用率等。

追问：`torch.profiler`、Nsight Systems、`py-spy` 分别分析什么？

最佳回答：

| 工具             | 主要分析对象               | 典型用途                           |
| ---------------- | -------------------------- | ---------------------------------- |
| `torch.profiler` | PyTorch 算子、CPU/GPU 活动 | 发现慢算子和执行热点               |
| Nsight Systems   | CPU/GPU 全局时间线         | 发现 Kernel 间空隙、同步和拷贝开销 |
| `py-spy`         | Python 调用栈              | 排查 Python 侧 CPU 热点和阻塞      |

例如我发现 GPU 只有 30% 利用率。

先用 Nsight Systems 观察 GPU 是否出现大段空闲。

如果 GPU 空闲时 CPU 正在处理数据，可以通过 `py-spy` 定位 Python 侧热点。

如果 GPU 一直忙但推理依然慢，则结合 `torch.profiler` 识别主要算子，再决定是否深入 Kernel 分析。

面试加分点： GPU 利用率低只是现象，需要先定位 CPU、I/O、数据加载、同步还是 GPU 计算瓶颈，而不是直接增加并发。



# 5 HAMi



HAMi（Heterogeneous AI Computing Virtualization Middleware） 是一个开源的 Kubernetes GPU 虚拟化与异构算力调度项目，前身叫 `k8s-vGPU-scheduler`。

它的核心作用是：

让多个 Kubernetes Pod 共享同一张物理 GPU，并根据显存、算力等资源需求进行调度和限制。

它是 AI Infra、GPU 云平台和多租户模型推理平台中很值得了解的技术。

GitHub



## 一、HAMi 解决了什么问题？

假设一台服务器只有一张 24 GB 显存的 GPU。

三个推理服务分别需要 4 GB、6 GB、8 GB 显存。

传统独占式分配

物理 GPU · 24 GB

Pod A

使用 4 GB

其余 20 GB 未使用

B、C 无法申请独占 GPU

HAMi 共享分配

物理 GPU · 24 GB

Pod A · 4 GB

Pod B · 6 GB

Pod C · 8 GB

剩余 6 GB

3 个 Pod 共享同一张 GPU

传统 NVIDIA Device Plugin 默认以整卡为主要分配单位。一张 GPU 被分配给一个 Pod 后，其他 Pod 即使只需要少量显存，也不能再次请求这张完整 GPU。

HAMi 则引入更细粒度的共享资源管理。

例如：

```
resources:
  limits:
    nvidia.com/gpu: 1
    nvidia.com/gpumem: 4096
    nvidia.com/gpucores: 30
```

在 HAMi 的资源语义下，这表示当前容器请求一张逻辑 GPU 分配、约 4 GiB 显存配额和 30% 的计算资源限额。

注意：这并不是把物理 GPU 真正切割成几个硬件设备。

## 二、HAMi 的整体架构

HAMi 的经典 Device Plugin 模式主要包含四部分：

| 组件               | 作用                                |
| ------------------ | ----------------------------------- |
| Mutating Webhook   | 修改需要 HAMi 管理的 Pod 调度配置   |
| Scheduler Extender | 根据显存、算力和设备状态选择 GPU    |
| HAMi Device Plugin | 将设备和虚拟化组件注入容器          |
| HAMi-Core          | 在用户态拦截 CUDA API，限制资源使用 |

用户创建 GPU Pod

声明显存和算力需求

Mutating Webhook

配置 HAMi 调度入口

Kubernetes Scheduler + HAMi Extender

选择 Node 和具体物理 GPU

Node 上的 Kubelet + Device Plugin

分配设备、注入库及资源限制配置

容器内的 HAMi-Core / libvgpu.so

拦截 GPU API，执行显存与算力限制

多个 Pod 共享同一物理 GPU

HAMi

## 三、核心原理一：HAMi 如何实现显存隔离？

这是理解 HAMi 最重要的一部分。

HAMi 并不是直接修改 NVIDIA 驱动，而是采用用户态 CUDA API 拦截（Hook）。

例如，应用程序通过 CUDA API 申请显存：

```
cuMemAlloc(&ptr, size);
```

正常情况下，该调用最终进入 NVIDIA CUDA 驱动接口。

而 HAMi 会在容器中加载 `libvgpu.so`，拦截指定的 CUDA 函数。

PyTorch / CUDA Application

申请 2 GB GPU 显存

libvgpu.so — CUDA API Hook

检查当前用量 + 本次申请量

未超出配额

转发实际 CUDA 调用

超出配额

返回 CUDA OOM 错误

假设 Pod A 的显存配额是 4096 MiB：

- 已使用 3000 MiB。
- 再申请 500 MiB，允许。
- 再申请 1500 MiB，会超过配额，因此拒绝。

同时，HAMi 也会拦截部分 NVML 查询接口，使容器内查询到的显存信息符合其虚拟化视图。

### libvgpu.so 是怎么自动进入容器的？

HAMi Device Plugin 在分配设备时，会为容器注入相应的库文件和动态链接配置。

经典实现利用 `/etc/ld.so.preload`，让 Linux 动态链接器优先加载 `libvgpu.so`。

该库通过符号拦截机制，将指定的 CUDA API 调用导向自己的逻辑。

这属于软件层面的隔离，并非硬件显存分区。

所以它有一个重要限制：不经过被拦截 CUDA API 路径的访问，可能无法被 HAMi 完整限制。

HAMi



## 四、核心原理二：HAMi 如何限制 GPU 算力？

显存隔离比较容易理解，但算力限制更加复杂。

假设某个 Pod 申请：

```
nvidia.com/gpucores: 30
```

这不代表物理 GPU 的 30% CUDA Core 被永久划分给这个 Pod。

HAMi 采用软件方式控制 CUDA Kernel 的提交节奏。

基本思路是：

1. 拦截 CUDA Kernel Launch 相关调用。
2. 根据使用情况与配置的算力限额进行节流。
3. 在达到限制时延迟后续 Kernel 提交。
4. 通过周期性调整，使工作负载的计算使用量趋近设定配额。

这是一种近似的算力使用控制机制，不是严格的 SM 硬件分区，也不保证任何时刻都恰好占用 30% 的 GPU 算力。

HAMi



## 五、HAMi、MIG、Time-Slicing 有什么区别？

这是 AI Infra 面试很可能出现的追问。

| 维度     | HAMi                      | NVIDIA MIG         | NVIDIA Time-Slicing |
| -------- | ------------------------- | ------------------ | ------------------- |
| 实现方式 | 软件 API 拦截             | GPU 硬件分区       | 分时共享            |
| 显存配额 | 支持软件限制              | 硬件隔离           | 无独立硬限制        |
| 算力管理 | 软件限速                  | 硬件资源分区       | 分时执行            |
| 硬件要求 | 更广泛的受支持 GPU        | 特定数据中心 GPU   | 较广泛的 GPU        |
| 隔离强度 | 依赖 API 拦截机制         | 较强               | 相对较弱            |
| 适用场景 | 多租户推理、轻量 GPU 任务 | 强隔离、多租户计算 | 简单共享 GPU        |

如何选择？

如果你的目标是让多个轻量模型服务共享一张 GPU，同时希望对每个 Pod 的显存进行限制，HAMi 比较适合。

如果业务需要更强的硬件隔离，且使用支持 MIG 的 GPU，则可以考虑 MIG。

如果只需要简单共享 GPU，且不需要细粒度显存隔离，可以考虑 NVIDIA Time-Slicing。

## 六、HAMi 与 Kubernetes Scheduler 究竟是什么关系？

这里很容易产生误解。

HAMi 不只是一个 Scheduler，也不只是一个 Device Plugin，而是将调度、设备分配和运行时限制结合起来。

在经典模式中：

- Kubernetes Scheduler 负责 Pod 的整体调度流程。
- HAMi Scheduler Extender 负责 GPU 层面的细粒度调度决策。
- HAMi Device Plugin 负责节点设备分配与容器注入。
- HAMi-Core 负责运行时 API 限制。

假设一张 24 GB GPU 已经分配了 8 GB + 6 GB。

新的 Pod 请求 4 GB。

HAMi 会查询它维护的设备资源视图，判断这张 GPU 是否满足资源要求，再结合调度策略决定是否分配。

它可以支持 Binpack 或 Spread 等策略。

这里也要注意一个较新的变化：HAMi 还在发展基于 Kubernetes DRA（Dynamic Resource Allocation）的实现。在 DRA 模式下，具体设备调度可以更多地交由 Kubernetes 原生资源分配机制承担，不再完全沿用经典 Scheduler Extender 的路径。

HAMi



## 七、放在具身智能仿真评测平台中，HAMi 有什么价值？

这个问题与你刚才准备的岗位高度相关。

假设公司部署了 200 台 GPU 服务器，运行 Isaac Sim、MuJoCo 和多个 VLA 模型推理服务。

这些工作负载的资源需求不同：

| 工作负载             | GPU 使用特点                | HAMi 的价值        |
| -------------------- | --------------------------- | ------------------ |
| 小型 VLA 推理服务    | 显存占用有限、GPU 间歇计算  | 适合评估共享       |
| 轻量模型测试         | 单个任务利用率低            | 提高资源密度       |
| Isaac Sim 大规模仿真 | GPU、显存、渲染资源需求较大 | 需要实际测试       |
| 多 GPU 模型训练      | 依赖通信拓扑和计算性能      | 未必适合细粒度共享 |

值得注意：不是所有 GPU 任务都应该使用 HAMi。

大型仿真、重负载训练可能更适合整卡调度，以减少资源争用和性能不确定性。

一个更成熟的设计是将 GPU 资源分池：

- 共享 GPU 池：服务轻量推理、小型测试。
- 独占 GPU 池：服务重负载仿真和训练。
- MIG GPU 池：服务需要硬件级资源分区的工作负载。

## 八、面试官会怎么问 HAMi？

面试官：请介绍 HAMi 的实现原理。

可以直接这样回答：

> HAMi 是 Kubernetes 上的 GPU 共享与虚拟化方案，主要解决传统 Device Plugin 整卡分配造成的资源浪费问题。
>
> 它在经典模式下通过 Mutating Webhook、Scheduler Extender、Device Plugin 和 HAMi-Core 协同工作。
>
> Scheduler Extender 根据 GPU 显存和算力资源选择设备；Device Plugin 将设备和虚拟化组件注入容器；HAMi-Core 使用用户态 CUDA API Hook 实现显存申请限制和计算节流。
>
> 它的特点是不需要修改应用代码，也通常不需要修改 GPU 驱动。
>
> 不过它属于软件层面的虚拟化，隔离能力依赖拦截的 API 路径，不能等同于 MIG 这种硬件分区。
>
> 如果用于 AI Infra 平台，我会根据任务的资源特点，将 GPU 分成共享池和独占池，以兼顾利用率和性能稳定性。

我建议你继续深入的两个知识点是：

1. Device Plugin 的 `Register`、`ListAndWatch`、`Allocate` 到底如何工作？ 这是 Kubernetes 设备管理的核心链路。
2. HAMi 的 `libvgpu.so` 为什么能拦截 CUDA API？`LD_PRELOAD`、动态链接、`dlsym` 是怎么工作的？ 这是理解 HAMi 真正底层原理的关键。

官方资料：[HAMi GPU 虚拟化原理](https://project-hami.io/docs/core-concepts/gpu-virtualization?utm_source=chatgpt.com) · [HAMi GitHub 源码](https://github.com/Project-HAMi/HAMi?utm_source=chatgpt.com)



# 6 如果要基于 Kubernetes 设计一个 AI 模型推理平台，你会怎么做？



我会分为控制面和数据面。

控制面负责模型注册、版本管理、部署声明、资源调度、自动扩缩容和生命周期管理。

数据面负责实际模型加载、推理执行以及请求处理。

例如用户提交一个模型部署请求，可以将模型名称、版本、GPU 需求、副本数、运行时配置定义为自定义资源。

Controller 负责把它转换为 Kubernetes Workload，再由 Scheduler 和设备管理组件分配合适的计算资源。

推理侧可以根据实际需求选择 vLLM、Triton 等推理运行时。

这里最大的挑战之一是 GPU 资源管理：

- 如何实现 GPU 设备发现和分配。
- 如何减少 GPU 显存碎片和资源浪费。
- 如何支持模型的弹性扩缩容。
- 如何优化模型加载与首 Token 延迟。
- 如何对不同租户实现资源隔离。

如果涉及高并发推理，还要进一步考虑请求批处理、KV Cache、模型并行和推理实例调度。

但需要区分：Kubernetes Device Plugin 提供设备资源接入，并不意味着 Kubernetes 自动具备细粒度 GPU 共享能力。 这需要硬件能力和额外的设备管理或调度方案。



# 7、Agent Infra 和普通 AI Infra 有什么区别？如何为 Agent 设计安全执行环境？

一般的 AI Infra 更多关注模型训练、推理、资源调度和计算基础设施。

Agent Infra 则进一步关注 Agent 执行任务所需的运行环境，例如代码执行、工具调用、任务调度、上下文状态管理和执行权限。

如果让我设计一个 Agent 执行平台，我会优先采用控制面与执行环境分离的方案。

控制面负责 Agent 任务的创建、调度、生命周期、状态管理和审计。

执行面则通过容器或者更强隔离的沙箱运行具体任务。

需要解决几个关键问题：

第一是隔离。不同用户的 Agent 不应互相访问文件、进程和敏感信息。对于执行不可信代码的场景，可以考虑 gVisor、Kata Containers 或微虚拟机等更强隔离方案。

第二是资源控制。限制 CPU、内存、执行时间、并发数以及网络访问范围。

第三是生命周期管理。支持超时回收、失败重试、任务取消、状态恢复。

第四是可观测性与审计。记录任务状态、执行日志、资源使用以及关键工具调用。

第五是安全访问。使用最小权限、短期凭证、网络出口限制和密钥隔离，避免 Agent 任意访问宿主机或内部服务。

如果要求高并发和低延迟，还需要沙箱池、资源预热和弹性调度。





