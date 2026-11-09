# 1、Linux CPU 达到 100%，怎么定位？

### 标准回答

首先区分 CPU 使用率高究竟来自用户态、内核态、中断处理，还是虚拟化环境中的资源争抢。

```
top
mpstat -P ALL 1
pidstat -u -t 1
vmstat 1
```

主要关注：

| 指标      | 含义                            |
| --------- | ------------------------------- |
| `%usr`    | 用户态计算                      |
| `%sys`    | 内核态执行                      |
| `%soft`   | 软中断处理                      |
| `%iowait` | CPU 空闲且存在等待 I/O 的时间   |
| `%steal`  | 虚拟机等待宿主机分配 CPU 的时间 |

如果 `%usr` 高，优先检查应用热点，例如 Go 服务的 pprof CPU Profile。

如果 `%sys` 高，检查系统调用、网络收发、锁竞争和频繁上下文切换。

如果 `%soft` 高，则重点检查网络中断、网卡队列以及软中断负载。

### 面试官追问

Q：CPU 使用率不高，Load Average 为什么很高？

因为 Linux Load Average 不只是统计正在使用 CPU 的任务，还包含不可中断睡眠状态等任务。

例如大量进程正在等待磁盘 I/O，可能导致 Load Average 很高，但 CPU 使用率并不高。

排查：

```
vmstat 1
ps -eo pid,stat,wchan:32,comm
iostat -xz 1
```

Q：上下文切换为什么会影响性能？

上下文切换涉及线程调度状态保存与恢复，还可能带来 CPU Cache、TLB 等局部性损失。

如果大量线程频繁竞争锁，就可能出现上下文切换很高，但业务吞吐量很低的现象。

对于 Kubernetes，还要检查 CFS CPU Throttling。CPU 使用率看起来不高，也可能因为容器持续受到 CPU Limit 限流而产生高延迟。



# 2、Linux Namespace 和 Cgroups 有什么区别？

### 标准回答

Namespace 负责资源视图隔离，Cgroups 负责资源控制和统计。

| 机制              | 作用                   | 示例                     |
| ----------------- | ---------------------- | ------------------------ |
| PID Namespace     | 隔离进程编号空间       | 容器内 PID 1             |
| Network Namespace | 隔离网络栈             | 独立网卡、路由、端口空间 |
| Mount Namespace   | 隔离挂载视图           | 容器文件系统             |
| UTS Namespace     | 隔离主机名             | 独立 hostname            |
| IPC Namespace     | 隔离 IPC 资源          | 共享内存、消息队列       |
| User Namespace    | 隔离用户和组 ID        | 容器内 Root 映射         |
| Cgroups           | 管理 CPU、内存、I/O 等 | CPU Limit、Memory Limit  |

### 追问：容器为什么不是虚拟机？

容器通常共享宿主机 Linux 内核，通过 Namespace、Cgroups、文件系统和安全机制实现隔离。

虚拟机则运行独立的 Guest OS 内核，通过 Hypervisor 提供硬件虚拟化。

### 追问：容器内存超出 Limit 会怎样？

对于 Cgroups v2，`memory.max` 是内存使用硬限制。当内存回收无法使使用量降到限制以下时，就可能触发该 Cgroup 内的 OOM 处理。

Kubernetes 中通常表现为容器进程被杀死，随后容器状态出现 `OOMKilled`。

CPU Limit 与 Memory Limit 的处理不同：

- CPU 超过配额通常会被限流。
- 内存达到硬限制可能触发 OOM Kill。

这是一个常见面试陷阱。

# 3、Linux 常见内核网络参数有哪些？如何调优？



### 标准回答

我会根据实际瓶颈分类调优，而不是套用通用参数模板。

| 参数                             | 作用                         |
| -------------------------------- | ---------------------------- |
| `net.core.somaxconn`             | 限制 listen backlog 上限     |
| `net.ipv4.tcp_max_syn_backlog`   | 控制未完成握手的请求队列上限 |
| `net.ipv4.ip_local_port_range`   | 本地自动分配临时端口的范围   |
| `net.netfilter.nf_conntrack_max` | 连接跟踪表容量上限           |
| `net.core.rmem_max`              | Socket 接收缓冲区可配置上限  |
| `net.core.wmem_max`              | Socket 发送缓冲区可配置上限  |

### 面试官追问

Q：为什么不能直接把参数调得很大？

例如增大 `nf_conntrack_max` 需要更多内核内存，也可能增加连接跟踪维护成本。

扩大 TCP 队列也不会自动提升应用消费能力。如果应用处理速度不足，只会让请求排队时间变长。

因此我的调优流程是：

1. 通过指标确定瓶颈。
2. 找出相关内核参数及其约束。
3. 在测试环境验证效果。
4. 逐步调整并监控 P99、吞吐量、丢包率和资源使用。
5. 保留回滚方案。

面试加分回答：

> 我认为内核调优首先是容量和瓶颈分析问题。只有确定了瓶颈，参数调整才有意义。否则很容易通过扩大队列掩盖上游或应用层问题。