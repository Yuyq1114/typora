![](images\k8s_architecture.png)



https://kubernetes.io/zh-cn/docs/tutorials/kubernetes-basics/

https://kubernetes.io/zh-cn/docs/concepts/

https://kubernetes.io/docs/concepts/overview/



## 思考和学习路径

k8s本身作为系统是怎么实现可用 正确  性能的。。。

Kubernetes 自身
│
├── 1. 可用性 Availability
├── 2. 正确性 Correctness
├── 3. 性能 Performance
├── 4. 安全 Security
├── 5. 可维护/可观测/可演进 Operability
└── 6. 成本/资源效率/可扩展性 Efficiency & Scalability

k8s作为基础设施 是怎么保证其管理的组件的 可用 正确 性能。。。以及



## 实现过程

### 1、第一版做一个 **ConfigMap Controller**：

```
API Server
    ↓
Reflector
    ↓
DeltaFIFO
    ↓
Informer
    ↓
Indexer / Local Cache
    ↓
EventHandler
    ↓
WorkQueue
    ↓
Worker
    ↓
SyncHandler / Reconcile
    ↓
client-go
    ↓
API Server
```

### 2、第二阶段：升级成 WebApp CRD



### 3、第三阶段：client-go 操作 CRD

Kubernetes API 本质上并不在乎 Go struct。



### 4、第四阶段：Informer    Watch

```
kube-apiserver

      │ WATCH

      ▼

Reflector

      │

      ▼

DeltaFIFO

      │

      ▼

Informer

   ┌──┴───┐
   ↓      ↓

Cache    Handler
          │
          ▼
       WorkQueue
```

### 5、第五阶段：WorkQueue

```
Event
 ↓
Queue
 ↓
Worker
 ↓
Reconcile
```

Burst Event
Retry
Rate Limit
Failure
Concurrency
Deduplication

### 6、第六阶段：Lister / Cache

```
API Server
    │
    │ WATCH
    ▼
Informer
    │
    ▼
Cache
    │
    ▼
Lister
    │
    ▼
Controller
```

### 7、第七阶段：真正写 Reconcile

```
Desired State

WebApp:
replicas = 3

       VS

Actual State

Deployment:
replicas = 2

       ↓

diff

       ↓

Update Deployment → 3
```

### 8、第八阶段：故意搞破坏

Controller 并不是“执行用户命令”，而是在持续维护系统不变量。

### 9、第九阶段：OwnerReference

```
WebApp DELETE
      ↓
Deployment DELETE
Service DELETE
```

OwnerReference
Garbage Collector
Controller

### 10、第十阶段：Status

```
Spec
=
用户声明的 Desired State

Status
=
Controller 观察到的 Actual State
```

### 11、最后

```
type WebApp struct {

    metav1.TypeMeta

    metav1.ObjectMeta

    Spec WebAppSpec

    Status WebAppStatus
}
```

Scheme
GroupVersion
TypeMeta
ObjectMeta
DeepCopy
Code Generator
Typed Client
Informer
Lister

### 12、扩展

#### 第一块扩展：自己写一个 Mini Scheduler

Controller 做的是：

> **应该存在什么？**

Scheduler 解决的是：

> **这个 Pod 应该放在哪里？**

```
CPU / Memory Filter
        ↓
NodeSelector
        ↓
Taint / Toleration
        ↓
Affinity
        ↓
Score
        ↓
Bind
```

Scheduling Framework

PreFilter
Filter
PostFilter
PreScore
Score
Reserve
Permit
PreBind
Bind
PostBind

#### 第二块扩展：自己写一个 Mini Kubelet

谁真正把 container 跑起来？

```
Mini Kubelet

Watch:
Pod.spec.nodeName == my-node

        ↓

发现新 Pod

        ↓

调用 containerd

        ↓

启动 Container

        ↓

监控 Container

        ↓

更新 Pod Status
```

#### 第三块扩展：Admission Webhook  ，API Server

```
Request
   ↓
Authentication
   ↓
Authorization
   ↓
Admission
   ↓
Validation
   ↓
Conversion
   ↓
Storage
   ↓
etcd
```

#### 第四块：Mini Operator，网络、存储和 Linux

```
Pod
 ↓
kubelet
 ↓
CRI
 ↓
Container Runtime
 ↓
CNI
 ↓
veth
 ↓
network namespace
 ↓
route / iptables / eBPF
```

```
PVC
 ↓
PV
 ↓
CSI Controller
 ↓
Attach
 ↓
CSI Node
 ↓
Mount
 ↓
Pod
```

```
resources:
  limits:
    memory: 1Gi
    cpu: "1"
```

```
Kubernetes
 ↓
Container Runtime
 ↓
runc
 ↓
cgroup
namespace
mount
capabilities
seccomp
```



## 基本概念

### 1. 集群（Cluster）

- **集群**是 Kubernetes 的基本工作单元，包含多个物理或虚拟机器，用于运行容器化的应用程序。集群由多个节点组成，包括主节点（Master Node）和工作节点（Worker Node）。

### 2. 主节点（Master Node）

- **主节点**是 Kubernetes 集群的控制平面，负责管理和调度整个集群的工作负载。主节点包含以下关键组件：
  - **API Server**：提供 Kubernetes API 的端点，用于管理集群的各种操作。
  - **Scheduler**：负责将新创建的 Pod 调度到集群中的节点上，考虑资源需求和约束。
  - **Controller Manager**：运行控制器，负责处理集群中的节点、副本控制、端点等资源的状态。
  - **etcd**：分布式键值存储，用于保存集群的配置信息、状态和元数据。

### 3. 工作节点（Worker Node）

- **工作节点**是 Kubernetes 集群中的计算资源节点，用于运行应用程序的容器实例。每个工作节点包含以下主要组件：
  - **Kubelet**：在节点上运行并管理容器的代理服务，与主节点的 API Server 通信，接收并执行 Pod 的创建、修改和删除等操作。
  - **Container Runtime**：负责运行容器，如 Docker、containerd 等。
  - **kube-proxy**：负责管理节点上的网络代理和负载均衡，维护网络规则和服务的转发。

### 4. Pod

- **Pod** 是 Kubernetes 中最小的调度单位，是一个或多个容器的组合。Pod 提供了一个独立的、共享网络和存储空间的运行环境，它包含了应用程序的容器以及共享的存储卷、网络 IP 和配置设置。

### 5. 控制器（Controller）

- **控制器**是 Kubernetes 中的一个核心概念，用于管理工作负载和资源对象的状态。常见的控制器包括：
  - **ReplicaSet**：确保指定数量的 Pod 副本在任何时候都在运行。
  - **Deployment**：用于定义和管理应用程序的发布版本，支持滚动更新和回滚操作。
  - **StatefulSet**：用于管理有状态应用程序的部署，如数据库。

### 6. 服务（Service）

- **服务**是 Kubernetes 中定义的一种抽象，用于定义一组 Pod 的逻辑集合和访问方式。服务提供了稳定的 DNS 名称和统一的入口点，用于在应用程序之间或外部网络之间进行通信。

### 7. 命名空间（Namespace）

- **命名空间**是用于在 Kubernetes 集群中对资源进行逻辑分组和隔离的一种方式。通过命名空间，可以将集群内的资源划分为不同的逻辑单元，以实现资源隔离和管理。

### 8. 存储卷（Volume）

- **存储卷**是 Kubernetes 中的一种抽象，用于持久化应用程序的数据。存储卷可以附加到 Pod 中的容器，使得数据在容器重新调度时不会丢失。

- Kubernetes（K8s）中的**存储卷（Volume）** 概念用于为Pod提供持久化和临时存储。Kubernetes 中的卷提供了对不同存储后端的抽象，允许应用程序将其数据持久化，即使 Pod 被删除或重新调度，存储数据也不会丢失。

  存储卷主要分为两类：

  1. **短期存储卷**（Pod 生命周期内存储）
  2. **持久存储卷**（跨越Pod生命周期的持久化存储）

  下面详细介绍 Kubernetes 中涉及到的存储卷的相关概念：

  #### 1. **Volume（卷）**

  Kubernetes 的卷（Volume）是用于挂载到 Pod 的目录，可以在 Pod 的多个容器之间共享文件系统。每个容器都可以通过挂载卷来访问共享数据。卷的生命周期与 Pod 绑定，Pod 结束时卷也会被清理，但不同的卷类型有不同的存储行为。

  ##### 关键特性：

  - **容器共享**：卷可以被多个容器共享，即 Pod 中的多个容器可以挂载同一个卷。
  - **数据持久性**：默认情况下，卷的生命周期与 Pod 绑定，Pod 删除时卷也会被删除，但某些卷类型（如 Persistent Volume）可以实现数据持久化。

  ##### 常见的 **短期存储卷** 类型：

  1. **emptyDir**：Pod 创建时临时分配的空目录，Pod 删除时也删除。
  2. **hostPath**：直接挂载主机文件系统的某个目录或文件到Pod中，适用于集群节点上有共享文件的场景。
  3. **configMap**：将Kubernetes的 `ConfigMap` 内容作为文件挂载到容器中，适合配置管理。
  4. **secret**：将 Kubernetes 的 `Secret` 内容作为文件挂载到容器中，适合传递敏感信息（如密码、令牌等）。
  5. **downwardAPI**：将 Pod 的元数据信息（如标签、注解等）挂载为文件，容器可以通过文件读取这些信息。

  ------

  #### 2. **Persistent Volume（PV）持久卷**

  **Persistent Volume（PV）** 是 Kubernetes 中的一种集群级别的资源，它代表了集群中的存储资源，通常由管理员预先配置。PV 是底层存储的抽象，存储类型可以是物理磁盘、云存储卷（如 AWS EBS、GCE Persistent Disk）、网络文件系统（如 NFS）、分布式存储（如 Ceph、GlusterFS）等。

  #### 关键特性：

  - **存储后端**：PV 可以是各种存储后端的抽象，如云存储、NFS、Ceph 等。

  - **容量定义**：PV 定义了可以提供的存储容量（如 10GiB）。

  - 访问模式

    ：

    - `ReadWriteOnce`：卷只能被一个节点以读写模式挂载。
    - `ReadOnlyMany`：卷可以被多个节点以只读模式挂载。
    - `ReadWriteMany`：卷可以被多个节点以读写模式挂载。

  - **生命周期**：PV 的生命周期独立于 Pod，它是集群管理员创建的资源，且在使用者不再需要时，集群管理员可以回收或重用 PV。

  - 回收策略

    ：

    - `Retain`：PV 释放后保留数据，需管理员手动清理。
    - `Recycle`：PV 释放后，Kubernetes 自动清理数据（仅支持某些存储后端）。
    - `Delete`：PVC 释放后自动删除 PV（适用于云存储）。

  #### PV 示例：

  ```
  apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: pv-nfs
  spec:
    capacity:
      storage: 10Gi
    accessModes:
      - ReadWriteOnce
    persistentVolumeReclaimPolicy: Retain
    nfs:
      path: /exported/path
      server: nfs-server.example.com
  ```

  ------

  #### 3. **Persistent Volume Claim（PVC）持久卷声明**

  **Persistent Volume Claim（PVC）** 是用户向 Kubernetes 请求存储的方式。PVC 声明了存储需求，如容量、访问模式等，Kubernetes 会根据 PVC 的要求寻找合适的 PV 并进行绑定。PVC 和 PV 之间的关系类似于计算机系统中的“请求-分配”机制。

  #### 关键特性：

  - **存储请求**：PVC 允许用户根据需要声明存储需求（如10Gi，访问模式 `ReadWriteOnce`）。
  - **与 PV 的绑定**：PVC 创建时，Kubernetes 会尝试找到符合条件的 PV，并将它们绑定在一起。
  - **动态存储供应**：如果没有符合要求的 PV，且 PVC 使用了 `StorageClass`，Kubernetes 可以通过 `StorageClass` 动态创建 PV。

  #### PVC 示例：

  ```
  apiVersion: v1
  kind: PersistentVolumeClaim
  metadata:
    name: pvc-nfs
  spec:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 10Gi
  ```

  ------

  #### 4. **StorageClass（SC）存储类**

  **StorageClass（SC）** 是 Kubernetes 中的一种资源，它定义了存储卷的动态供应策略。集群管理员可以通过 StorageClass 定义如何动态地为 PVC 创建 PV（即存储卷）。不同的存储后端可以有不同的 `StorageClass`，如 AWS EBS、GCE PD、NFS 等。

  ##### 关键特性：

  - **Provisioner（供应器）**：`StorageClass` 通过 Provisioner 指定具体的存储供应后端（如 `kubernetes.io/aws-ebs`、`kubernetes.io/gce-pd`）。
  - **动态存储供应**：PVC 请求时，如果没有现成的 PV 可用，Kubernetes 会根据 `StorageClass` 动态创建 PV。
  - **参数**：`StorageClass` 可以包含与存储后端相关的配置参数（如存储类型、区域等）。
  - **回收策略**：指定存储卷释放后的回收方式（如 `Retain`, `Delete`）。

  ##### StorageClass 示例：

  ```
  apiVersion: storage.k8s.io/v1
  kind: StorageClass
  metadata:
    name: fast
  provisioner: kubernetes.io/aws-ebs
  parameters:
    type: io1
    iopsPerGB: "10"
    fsType: ext4
  ```

  ------

  #### 5. **Volume Mounts（卷挂载）**

  **Volume Mounts** 是在 Pod 定义中将卷挂载到容器文件系统中的过程。通过 `volumeMounts` 指定容器内的挂载路径以及对应的卷。卷可以在多个容器之间共享，且容器可以对挂载点设置不同的权限（读写或只读）。

  ##### Volume Mounts 示例：

  ```
  apiVersion: v1
  kind: Pod
  metadata:
    name: volume-demo
  spec:
    containers:
    - name: container
      image: busybox
      volumeMounts:
      - mountPath: /data
        name: volume
    volumes:
    - name: volume
      persistentVolumeClaim:
        claimName: pvc-nfs
  ```

  ------

  #### 6. **空卷（emptyDir）**

  **emptyDir** 是最简单的卷类型。当 Pod 被调度到节点时，Kubernetes 会在该节点上创建一个空目录供 Pod 使用。`emptyDir` 的数据生命周期与 Pod 绑定，当 Pod 被删除时，数据也会被删除。

  ##### 关键特性：

  - **临时存储**：存储的数据在 Pod 运行期间是持久的，但 Pod 删除时，数据也会丢失。
  - **应用场景**：用于缓存数据、共享容器间数据等。

  ##### emptyDir 示例：

  ```
  apiVersion: v1
  kind: Pod
  metadata:
    name: emptydir-demo
  spec:
    containers:
    - name: container
      image: busybox
      volumeMounts:
      - mountPath: /data
        name: cache
    volumes:
    - name: cache
      emptyDir: {}
  ```

  ------

  #### 7. **HostPath**

  **HostPath** 卷允许容器直接挂载节点上的某个目录或文件。这种卷的主要应用场景是需要直接访问主机文件系统（如日志、配置文件或设备），但这种卷类型存在潜在的安全风险，因为它依赖于主机环境的文件系统。

  ##### HostPath 示例：

  ```
  apiVersion: v1
  kind: Pod
  metadata:
    name: hostpath-demo
  spec:
    containers:
    - name: container
      image: busybox
      volumeMounts:
      - mountPath: /data
        name: host-storage
    volumes:
    - name: host-storage
      hostPath:
        path: /mnt/data
  ```

### 9. 滚动更新和回滚

- **滚动更新**允许对部署的应用程序进行逐步更新，通过逐步替换旧版本的 Pod 实例来确保应用程序的可用性和稳定性。
- **回滚操作**允许将应用程序的部署版本回退到之前的版本，以应对错误或不良的更新。

### 10. 高可用性和自动伸缩

- Kubernetes 提供了高可用性的集群架构，通过在多个节点上调度应用程序的副本来确保应用程序的高可用性。
- 通过控制器和自动伸缩机制，Kubernetes 能够根据需求自动扩展或缩减工作节点上的容器实例数量，以满足应用程序的资源需求。



## 原理

### 1. 容器化基础

Kubernetes 构建在 Linux 容器（如 Docker）的基础上，利用容器技术实现了应用程序的隔离和封装。每个容器都包含了一个或多个应用程序实例及其运行时依赖，但它们共享主机操作系统的内核，因此启动速度快，资源占用少。

### 2. 集群架构

Kubernetes 集群由多个节点组成，主要分为主节点（Master Node）和工作节点（Worker Node）：

- **主节点（Master Node）**：
  - **API Server**：提供了 Kubernetes API 的入口，用于管理集群的各种操作，包括创建、更新和删除资源对象（如 Pod、Service）等。
  - **Scheduler**：负责将新创建的 Pod 调度到合适的工作节点上，考虑到节点的资源利用率和健康状态。
  - **Controller Manager**：运行多个控制器，如 ReplicaSet Controller、Deployment Controller，负责监控集群中的资源对象的状态，并进行必要的调节和修复。
  - **etcd**：分布式键值存储，保存了整个集群的配置信息、状态和元数据，作为主节点的持久化存储。
- **工作节点（Worker Node）**：
  - **Kubelet**：运行在每个节点上的代理服务，负责管理节点上的容器生命周期，与主节点的 API Server 通信，执行主节点下发的任务（如创建、删除 Pod）。
  - **Container Runtime**：负责在节点上运行容器，常见的包括 Docker、containerd 等。
  - **kube-proxy**：负责维护节点上的网络代理和负载均衡，为 Pod 提供网络服务，支持服务发现和负载均衡功能。

### 3. 控制器模式

Kubernetes 采用控制器模式来管理应用程序的部署和运行状态，确保集群中的工作负载符合用户定义的期望状态：

- **ReplicaSet**：确保指定数量的 Pod 副本在集群中运行，处理 Pod 的创建、删除和替换操作。
- **Deployment**：在 ReplicaSet 的基础上实现应用程序的声明式部署和更新，支持滚动更新和回滚操作。
- **StatefulSet**：用于管理有状态应用程序的部署，如数据库服务，保证每个 Pod 有唯一标识和稳定的网络标识。

### 4. 资源调度

Kubernetes 的调度器（Scheduler）负责将新创建的 Pod 分配到合适的工作节点上，考虑以下因素：

- **资源需求和限制**：每个 Pod 可以指定 CPU 和内存的需求与限制，调度器根据节点的可用资源进行匹配。
- **亲和性和反亲和性规则**：可以根据 Pod 之间的关系（如亲和性或反亲和性）来调度它们，确保在同一节点或不同节点上运行。
- **节点健康状态**：考虑节点的负载和健康状况，避免将 Pod 调度到资源紧张或故障的节点上。

### 5. 服务发现和负载均衡

Kubernetes 提供了内置的服务发现机制和负载均衡功能，使得应用程序可以稳定地访问和通信：

- **Service**：定义一组 Pod 的逻辑集合和访问方式，为 Pod 提供稳定的 DNS 名称和统一的入口点。
- **kube-proxy**：在每个节点上运行，维护集群中服务的网络代理和负载均衡规则，实现服务级别的负载均衡和流量转发。

### 6. 滚动更新和回滚

Kubernetes 支持应用程序的滚动更新和回滚操作，确保应用程序的持续可用性和稳定性：

- **滚动更新**：通过逐步替换旧版本的 Pod 实例来实现应用程序的更新，可配置更新策略和健康检查来确保更新过程的安全性。
- **回滚操作**：允许将应用程序的部署版本回退到之前的版本，以应对错误或不良的更新。

### 7. 自动伸缩

Kubernetes 支持基于资源使用情况和应用程序的指标进行自动伸缩，以满足变化的负载需求：

- **水平自动伸缩**（Horizontal Pod Autoscaler，HPA）：根据 CPU 使用率或自定义指标自动调整 Pod 的副本数量。
- **垂直自动伸缩**（Vertical Pod Autoscaler，VPA）：根据单个 Pod 内部资源（如内存、CPU）的使用情况调整 Pod 的资源请求和限制。

### 8. 命名空间和权限控制

Kubernetes 使用命名空间（Namespace）来进行资源的逻辑分组和隔离，同时提供了精细的 RBAC（Role-Based Access Control）机制来管理和控制用户对集群资源的访问权限。

sudo ctr images import kafka/kafka.tar



## 常用命令

### 集群操作命令

1. **连接到集群**

   ```
   kubectl config use-context CONTEXT_NAME
   ```

   - 使用指定的上下文连接到 Kubernetes 集群，其中 `CONTEXT_NAME` 是 `kubectl config get-contexts` 列出的上下文名称。

2. **查看集群信息**

   ```
   kubectl cluster-info
   ```

   - 显示集群的地址信息和状态。

3. **查看节点信息**

   ```
   kubectl get nodes
   ```

   - 列出集群中所有的节点及其状态。

4. **查看集群中的命名空间**

   ```
   kubectl get namespaces
   ```

   - 列出当前集群中所有的命名空间。

### 资源操作命令

1. **查看资源列表**

   ```
   kubectl get RESOURCE_TYPE
   ```

   - 列出指定资源类型的所有实例，如 `pods`、`services`、`deployments` 等。

2. **查看资源详细信息**

   ```
   kubectl describe RESOURCE_TYPE RESOURCE_NAME
   ```

   - 显示指定资源实例的详细信息，如 Pod、Service 等。
   - 查看某个namespace下的某个pod的信息：
     - kubectl describe pod kafk-0 -n dsamp

3. **创建资源**

   ```
   kubectl create -f FILE.yaml
   ```

   - 根据 YAML 或 JSON 文件中的定义创建资源，可以是 Pod、Service、Deployment 等。

4. **删除资源**

   ```
   kubectl delete RESOURCE_TYPE RESOURCE_NAME
   ```

   - 删除指定的资源实例，如 Pod、Service 等。
   - kubectl delete pod my-pod --grace-period=0 --force
   - 检查 PV 是否处于 `Terminating` 状态，这种状态可能被 `Finalizer` 阻止删除。kubectl edit pv <pv-name>找到 `finalizers` 字段并将其删除，保存退出。

5. **修改资源**

   ```
   kubectl apply -f FILE.yaml
   ```

   - 根据 YAML 或 JSON 文件中的定义修改或创建资源。如果资源已存在，则进行更新操作。

### 应用程序和服务管理命令

1. **查看应用程序日志**

   ```
   kubectl logs POD_NAME
   ```

   - 显示指定 Pod 的日志输出。

2. **进入 Pod 内部**

   ```
   kubectl exec -it POD_NAME -n <namaspace> -- /bin/bash
   kubectl exec -it POD_NAME -n <namaspace> -- /bin/sh //也可能是
   ```

   - 在指定 Pod 内部启动一个交互式 Shell。

   - ### 查看某个特定 Pod 的详细信息

     使用以下命令查看某个特定 Pod 的详细信息，包括状态、IP 地址、容器状态等：

   - kubectl describe pod <pod-name> -n <namespace>

   - 也可以内部发送消息

   - kafka-console-producer.sh --broker-list <broker-address>:<port> --topic <topic-name>

3. **暴露服务**

   ```
   kubectl expose RESOURCE_TYPE RESOURCE_NAME --port=PORT --target-port=TARGET_PORT --type=SERVICE_TYPE
   ```

   - 根据 Pod 或 Deployment 暴露服务，定义服务的类型和端口映射。

4. **扩展和缩减副本**

   ```
   kubectl scale --replicas=NUM REPLICATION_CONTROLLER_NAME
   ```

   - 扩展或缩减指定副本控制器（如 Deployment 或 ReplicaSet）的副本数量。

### 网络和存储管理命令

1. **查看集群中的服务**

   ```
   kubectl get services
   ```

   - 列出集群中所有的服务及其相关信息。

2. **查看存储卷**

   ```
   kubectl get pv
   ```

   - 列出集群中所有的持久化存储卷。
   - kubectl get pvc -n <namespace>
   - kubectl get pv -n <namespace>

3. **管理命名空间**

   ```
   kubectl create namespace NAMESPACE_NAME
   kubectl delete namespace NAMESPACE_NAME
   ```

   - 创建或删除命名空间。

### 其他常用命令

1. **查看 API 资源**

   ```
   kubectl api-resources
   ```

   - 列出 Kubernetes API 支持的资源类型。

2. **查看当前上下文**

   ```
   kubectl config current-context
   ```

   - 显示当前使用的 Kubernetes 集群上下文。

3. **设置命名空间**

   ```
   kubectl config set-context --current --namespace=NAMESPACE_NAME
   ```

   - 设置当前上下文的默认命名空间。



## 具体的对象

### 资源范围

查看：

例如：

```
kubectl api-resources
```

输出：

```
NAME              SHORTNAMES   NAMESPACED
pods              po           true
services          svc          true
persistentvolumes pv           false
storageclasses    sc           false
nodes                          false
```

其中：

- `true` → 属于 Namespace
- `false` → 集群级资源



CRD 自定义资源定义

CR 自定义资源实例

**YAML 只是声明“期望状态”；真正干活的是 API Server + 各种 Controller + Scheduler + Kubelet。**

| 分类 | 对象                           | 解决什么问题 |
| ---- | ------------------------------ | ------------ |
| 计算 | Pod / Deployment / StatefulSet | 程序怎么跑   |
| 网络 | Service / Ingress              | 怎么被访问   |
| 存储 | PV / PVC / StorageClass        | 数据放哪     |
| 配置 | ConfigMap / Secret             | 配置怎么给   |
| 管理 | Namespace / Label / Selector   | 怎么管       |

| 分类 | 对象                         | 解决什么问题                  |
| :--- | :--------------------------- | :---------------------------- |
| 计算 | Pod                          | 最小运行单元                  |
| 计算 | Deployment                   | 无状态应用，副本、滚动更新    |
| 计算 | StatefulSet                  | 有状态应用，固定身份和数据    |
| 计算 | DaemonSet                    | 每个节点都跑一个 Pod          |
| 计算 | Job                          | 一次性任务，跑完就结束        |
| 计算 | CronJob                      | 定时任务                      |
| 计算 | ReplicaSet                   | Deployment 底层，保证副本数   |
| 网络 | Service                      | 稳定访问入口                  |
| 网络 | Endpoints / EndpointSlice    | Service 后面实际有哪些 Pod IP |
| 网络 | Ingress                      | 外部 HTTP 入口                |
| 存储 | PV / PVC / StorageClass      | 数据放哪                      |
| 配置 | ConfigMap / Secret           | 配置怎么给                    |
| 管理 | Namespace / Label / Selector | 怎么管                        |
| 伸缩 | HPA                          | 按负载自动扩缩容              |

### 1、存储体系中

#### 1️⃣ 为什么需要它？

现实问题：

- Pod 随时会被删
- 节点会重启
- 容器文件系统不可靠

> **数据不能跟着 Pod 走**

------

#### 2️⃣ PV（PersistentVolume）

> **集群级的“硬盘资源”**

特点：

- 管理员创建
- 与 Pod 无关
- 生命周期独立

```
kind: PersistentVolume
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/mysql
```

工程直觉：

> PV = 后端存储的抽象描述

------

#### 3️⃣ PVC（PersistentVolumeClaim）

> **应用对存储的“申请单”**

```
kind: PersistentVolumeClaim
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

特性：

- 应用创建
- 不关心底层存储
- 只关心容量、访问模式

工程直觉：

> PVC = “我要一个 5G 硬盘”

------

#### 4️⃣ StorageClass（自动化关键）

> **“PVC → 自动生成 PV” 的规则**

```
kind: StorageClass
provisioner: rancher.io/local-path
```

有了它：

```
PVC 创建
↓
自动创建 PV
↓
绑定
```

👉 **平台工程的核心组件之一**

------

#### 5️⃣ 它们的关系

```
Pod → PVC → PV → 实际存储
```

### 2、网络体系中

#### 1️⃣ 为什么需要 Service？

Pod 的 IP：

- 会变
- 不稳定
- 不能依赖

> **访问 Pod 必须有稳定入口**

------

#### 2️⃣ Service 是什么？

> **一组 Pod 的稳定访问入口**

```
kind: Service
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
```

你访问的是：

```
mysql.default.svc.cluster.local
```

而不是 Pod IP。

------

#### 3️⃣ Service 的三种常见类型

| 类型         | 用途                 |
| ------------ | -------------------- |
| ClusterIP    | 集群内部访问（默认） |
| NodePort     | 节点端口暴露         |
| LoadBalancer | 云厂商负载均衡       |

------

#### 4️⃣ 工程直觉

> Service = 负载均衡 + 服务发现

Service 已经能让别人访问一组 Pod，但对外暴露时有局限：

| 方式         | 问题                                                 |
| :----------- | :--------------------------------------------------- |
| ClusterIP    | 只能集群内部访问                                     |
| NodePort     | 每个服务占一个节点端口，端口会很快不够用             |
| LoadBalancer | 每个服务通常要一个云负载均衡器，成本高、域名也不好管 |

如果集群里有很多 Web 服务，你真正想要的是：

一个域名 / 一个入口

  ├── api.example.com     → api-service

  ├── www.example.com     → web-service

  └── example.com/admin   → admin-service

这就是 Ingress。

#### 和 Service 的关系

Ingress 不直接选 Pod，它把请求转到 Service，再由 Service 转到 Pod：

### 3、配置体系

#### 1️⃣ 为什么不能写死在镜像里？

- 每个环境不同
- 频繁修改
- 不想重建镜像

> **配置必须外置**

------

#### 2️⃣ ConfigMap（明文配置）

```
kind: ConfigMap
data:
  my.cnf: |
    max_connections=200
```

用途：

- 配置文件
- 非敏感参数

------

#### 3️⃣ Secret（敏感配置）

```
kind: Secret
data:
  password: base64(...)
```

特点：

- base64（不是加密）
- RBAC 控制访问

------

#### 4️⃣ 它们怎么给 Pod？

##### 方式一：环境变量

```
env:
- name: MYSQL_PASSWORD
  valueFrom:
    secretKeyRef:
      name: mysql
      key: password
```

##### 方式二：挂载成文件（更常用）5

```
volumeMounts:
- mountPath: /etc/mysql
```

------

#### 5️⃣ 工程直觉

> ConfigMap / Secret = “运行时配置输入”





### 4、Deployment vs StatefulSet

| 对比   | Deployment | StatefulSet     |
| ------ | ---------- | --------------- |
| Pod 名 | 随机       | 固定（mysql-0） |
| 存储   | 通常共享   | 每 Pod 一个 PVC |
| 用途   | 无状态     | 有状态          |

MySQL / Redis / Kafka → **StatefulSet**

“有状态 / 无状态”不是指控制器本身，而是指应用是否依赖固定身份和持久数据。

无状态（Stateless）

典型特征：

- Pod 随便删、随便扩、随便换节点，业务都能继续
- 不依赖“我是第几个实例”
- 本地磁盘数据丢了通常没关系
- 多个副本之间基本等价，可互相替代

有状态（Stateful）

典型特征：

- 每个实例有固定身份
- 重启后还要能找回自己的数据
- 实例之间可能不能随意互换
- 启动顺序、主从关系、副本编号常常很重要

### 5、Namespace（资源隔离）

> **逻辑上的“租户 / 环境隔离”**

```
default
dev
test
prod
```

Helm Release 是 **绑定 namespace 的**。

------

### 6、Label / Selector（K8s 的“胶水”）

```
labels:
  app: mysql
```

Service 通过 selector 找 Pod：

```
selector:
  app: mysql
```

> **没有 Label，K8s 就无法工作**



## 总流程

### 1. 先建立总流程

你执行：

kubectl apply -f deployment.yaml

实际发生的是：

kubectl

  ↓ 提交 YAML

API Server

  ↓ 校验并写入

etcd（保存对象）

  ↓ 各 Controller 监听变化

Controller 对比“期望状态 vs 实际状态”

  ↓ 不满足就创建/修改/删除对象

Scheduler / Kubelet / kube-proxy / 外部 Controller 继续处理

这叫 控制循环（Reconcile Loop）：

观察当前状态 → 和期望状态对比 → 如果不一致就修正 → 再观察

所以：不是 Deployment 直接生成 Pod，而是 Deployment Controller 看到 Deployment 后，再去创建 ReplicaSet，再由 ReplicaSet Controller 创建 Pod。

------

### 2. Deployment 是怎么变成 Pod 的

### 对象链路

Deployment

  ↓ Deployment Controller 管理

ReplicaSet

  ↓ ReplicaSet Controller 管理

Pod

  ↓ Scheduler 调度

Node 上的 Kubelet 启动容器

### 各组件职责

| 组件                  | 干什么                              |
| :-------------------- | :---------------------------------- |
| API Server            | 接收 YAML，存到 etcd                |
| Deployment Controller | 看 Deployment，创建/更新 ReplicaSet |
| ReplicaSet Controller | 看 ReplicaSet，保证 Pod 副本数正确  |
| Scheduler             | 给 Pending 状态的 Pod 选节点        |
| Kubelet               | 在节点上拉镜像、挂卷、启容器        |
| Container Runtime     | 真正跑容器，如 containerd           |

### 为什么中间有 ReplicaSet

Deployment 本身不直接管 Pod，它管 ReplicaSet：

Deployment(replicas=3)

  → 创建 ReplicaSet

  → ReplicaSet 创建 3 个 Pod

升级时：

Deployment 更新镜像

  → 新建一个 ReplicaSet

  → 逐步缩旧 RS、扩新 RS

  → 实现滚动更新

所以平时你写 Deployment，ReplicaSet 通常是 K8s 自动生成的，不用手写。

------

### 3. StatefulSet 是谁创建 Pod 的

和 Deployment 类似，但控制器不同：

StatefulSet

  ↓ StatefulSet Controller

Pod（按顺序：mysql-0, mysql-1, mysql-2）

  ↓ 同时可能创建

PVC（每个 Pod 一个）

StatefulSet Controller 负责：

- 按固定顺序创建/删除 Pod
- 保证 Pod 名字稳定
- 为每个 Pod 创建专属 PVC（如果配置了 `volumeClaimTemplates`）

------

### 4. PVC / PV / StorageClass 是谁处理的

存储这条链和 Deployment 不太一样，涉及多个控制器和外部组件。

### 关系

StorageClass（存储规则）

  ↓

PVC（申请存储）

  ↓

PV（真正的一块存储）

  ↓

Pod 挂载 PVC 使用

### 谁负责哪一步

| 组件                     | 负责什么                                 |
| :----------------------- | :--------------------------------------- |
| StorageClass             | 只是规则，本身不创建 Pod                 |
| PVC                      | 应用/workload 引用它申请存储             |
| External Provisioner     | 看到 PVC 后，按 StorageClass 动态创建 PV |
| Attach/Detach Controller | 把云盘挂到节点                           |
| PV Controller            | 处理 PV 和 PVC 的绑定关系                |
| Kubelet                  | Pod 启动时把 PVC 对应的卷挂到容器里      |

### 动态供应流程

StatefulSet Controller 创建 PVC

  ↓

External Provisioner 监听到 PVC

  ↓

按 StorageClass 创建 PV

  ↓

PV Controller 绑定 PV 和 PVC

  ↓

Scheduler 调度 Pod 到合适节点

  ↓

Attach Controller 把盘挂到节点

  ↓

Kubelet 挂载卷并启动容器

### 静态供应流程

如果没有 StorageClass，而是管理员预先创建了 PV：

你创建 PVC

  ↓

PV Controller 找匹配的 PV 并绑定

  ↓

Pod 引用 PVC

  ↓

Kubelet 挂载

------

### 5. Service / Ingress 是谁处理的

### Service

Service YAML

  ↓ API Server 保存

Endpoints/EndpointSlice Controller

  ↓ 根据 Label Selector 找到 Pod IP

生成 Endpoints

  ↓

kube-proxy 读取 Service + Endpoints

  ↓

配置节点 iptables/ipvs 规则

  ↓

流量转发到 Pod

| 组件                 | 作用                                  |
| :------------------- | :------------------------------------ |
| Service              | 定义访问入口和 selector               |
| Endpoints Controller | 维护 Service 后面有哪些 Pod IP        |
| CoreDNS              | 提供 `xxx.svc.cluster.local` 域名解析 |
| kube-proxy           | 节点上做负载均衡转发                  |

### Ingress

Ingress 比较特殊：K8s 内置没有通用 Ingress Controller。

Ingress YAML（规则）

  ↓

Ingress Controller（nginx/traefik/ALB 等，需单独部署）

  ↓ 读取 Ingress 规则

转发到对应 Service

  ↓

Service → Pod

所以：

- Ingress 对象：K8s 原生资源
- Ingress Controller：额外部署的控制器，不是 kube-controller-manager 自带的

------

### 6. ConfigMap / Secret 是谁用的

它们本身只是配置对象，不会“运行”：

ConfigMap / Secret

  ↓

被 Deployment / StatefulSet / Pod 引用

  ↓

API Server 保存对象

  ↓

Kubelet 启动 Pod 时

  \- 注入环境变量

  \- 或挂载成文件

  ↓

容器内应用读取

这里没有专门的“ConfigMap Controller 去跑业务”，Kubelet 负责把配置塞进 Pod。

------

### 7. 一张“谁用谁”的总关系图

你 apply YAML

​    ↓

API Server → etcd

Deployment ──→ Deployment Controller ──→ ReplicaSet

ReplicaSet ──→ ReplicaSet Controller ──→ Pod

StatefulSet ─→ StatefulSet Controller ─→ Pod + PVC

Pod ──→ Scheduler ──→ 选 Node

Pod ──→ Kubelet ──→ 启容器 + 挂 PVC

PVC ──→ External Provisioner ──→ PV

PV/PVC ─→ PV Controller ──→ 绑定

Pod + PVC ─→ Attach Controller ──→ 挂盘到节点

Service ──→ Endpoints Controller ──→ Endpoints

Service + Endpoints ──→ kube-proxy ──→ 转发流量

Ingress ──→ Ingress Controller ──→ 转发到 Service

ConfigMap/Secret ──→ 被 Pod 引用 ──→ Kubelet 注入/挂载





## 问题和解决方案（内部算法）

### 1. **容器调度**

#### **问题描述**：

容器调度可能会导致资源分配不均或某些节点负载过重。

#### **解决措施**：

- **调度器算法**：
  - **调度策略**：Kubernetes 的调度器使用多种调度策略（如 `LeastRequested`, `MostRequested` 和 `BalancedResourceAllocation`）来选择最适合的节点。
  - **亲和性和反亲和性**：通过节点亲和性（`nodeAffinity`）和 Pod 亲和性（`podAffinity` 和 `podAntiAffinity`）来控制 Pods 的调度，优化资源使用和负载均衡。
- **资源请求和限制**：
  - **资源配额**：为 Pods 设置 CPU 和内存的请求和限制，以确保资源合理分配和防止资源竞争。

### 2. **网络问题**

#### **问题描述**：

网络延迟、网络分区或网络配置错误可能导致容器间通信问题或服务不可用。

#### **解决措施**：

- **网络插件**：
  - **CNI 插件**：Kubernetes 支持多种 CNI（Container Network Interface）插件（如 Calico、Flannel 和 Weave），提供高效的网络通信和网络策略管理。
- **网络策略**：
  - **网络隔离**：使用网络策略（`NetworkPolicy`）来控制 Pods 之间的网络流量，增强安全性和网络管理能力。
- **负载均衡**：
  - **内建负载均衡**：Kubernetes 提供内建的负载均衡功能，通过 Service 资源分发流量，确保服务的可用性和性能。

### 3. **存储问题**

#### **问题描述**：

持久化存储可能会遇到数据丢失、存储性能不佳或存储资源不足的问题。

#### **解决措施**：

- **持久卷（Persistent Volumes）**：
  - **动态卷供给**：使用动态卷供给（Dynamic Provisioning）来自动创建和管理持久卷，简化存储管理。
- **存储类（StorageClass）**：
  - **存储策略**：配置存储类来定义存储的性能和特性（如 SSD、HDD 和网络存储），根据需求选择合适的存储方案。
- **数据备份**：
  - **备份策略**：使用备份工具和策略定期备份持久卷的数据，防止数据丢失。

### 4. **服务发现**

#### **问题描述**：

服务发现可能会出现服务无法找到或连接失败的问题。

#### **解决措施**：

- **服务注册和发现**：
  - **内建 DNS**：Kubernetes 内建 DNS 服务（CoreDNS）用于服务发现和负载均衡，自动注册和解析服务名称。
- **头部路由**：
  - **Ingress**：使用 Ingress 资源和 Ingress 控制器管理 HTTP 和 HTTPS 路由，提供灵活的服务路由和负载均衡功能。

### 5. **扩展和弹性**

#### **问题描述**：

容器扩展可能会遇到自动扩展失效、资源不足或扩展速度慢的问题。

#### **解决措施**：

- **水平自动扩展（HPA）**：
  - **自动扩展**：使用 Horizontal Pod Autoscaler（HPA）根据 CPU 和内存利用率自动扩展或缩减 Pod 副本。
- **垂直自动扩展（VPA）**：
  - **资源调整**：使用 Vertical Pod Autoscaler（VPA）根据 Pod 的实际需求自动调整资源请求和限制。
- **集群自动扩展（CA）**：
  - **节点扩展**：使用 Cluster Autoscaler（CA）根据集群中的资源需求自动扩展或缩减节点。

### 6. **安全问题**

#### **问题描述**：

安全问题包括容器漏洞、未授权访问和配置错误。

#### **解决措施**：

- **RBAC（基于角色的访问控制）**：
  - **权限管理**：使用 RBAC 控制用户和服务的权限，确保只有授权的用户可以访问和管理资源。
- **安全上下文**：
  - **权限限制**：设置 Pod 的安全上下文（SecurityContext）来限制容器的权限，避免容器以 root 用户身份运行。
- **网络策略**：
  - **网络隔离**：使用网络策略控制 Pods 之间的网络访问，防止未授权的流量进入或离开 Pod。
- **镜像扫描**：
  - **漏洞扫描**：对容器镜像进行安全扫描，发现并修复镜像中的漏洞。

### 7. **监控和日志**

#### **问题描述**：

监控和日志问题可能导致无法及时发现问题或分析故障原因。

#### **解决措施**：

- **集成监控工具**：
  - **监控系统**：集成监控工具（如 Prometheus 和 Grafana）进行集群和应用的性能监控，及时发现和响应系统问题。
- **日志收集**：
  - **集中日志管理**：使用日志收集工具（如 ELK Stack 或 Fluentd）集中收集和分析日志，帮助排查和解决问题。

### 8. **配置管理**

#### **问题描述**：

配置管理可能会导致应用配置错误、配置不一致或部署失败。

#### **解决措施**：

- **ConfigMap 和 Secret**：
  - **配置管理**：使用 ConfigMap 和 Secret 管理应用配置和敏感数据，实现配置的灵活性和安全性。
- **版本控制**：
  - **配置版本化**：对配置进行版本控制，确保配置变更可追溯，并可以轻松回滚到先前版本。





## 其他：通过goland连接k8s

```
通过From default directory ,就是从默认的 C:\Users\1003716\.kube\config 文件下读取

通过From custom kubeconfigs ，就是从自定义的config文件读取，文件名无后缀

通过Parse kubeconfig content，就是复制远程机器的kubectlconfig文件读取

add时有个可选项是是否加到C:\Users\1003716\.kube\config 这个文件下，这会C:\Users\1003716\.kube\config下一层层多出来很多的配置，建议加上，因为后面的teleprensece连接的时候会默认用这个文件
```

```
1、修改  config  文件中的以下内容 不然goland自动玩会有冲突
clusters:
- name: k3s-212-45

contexts:
- context:
    cluster: k3s-212-45
    user: k3s-212-45-user
  name: k3s-212-45

users:
- name: k3s-212-45-user

current-context: k3s-212-45

2、如果 clusters: 下的 - cluster: 下的 server:  是127.0.0.1:6443 ,还需要改一下这个ip成对应目标ip，只需要改用于连接的本地config即可，不用改远端的config

3、如果不可以进行6443的连接 但是开放了ssh端口 可以将流量代理下走ssh，另外配置文件需要127.0.0.1:可用端口

注意：如果用数字需要加引号防止go build类型错误

4、teleprense使用
怎么配置 windows10有bug
telepresence quit -s
telepresence connect -n test1
telepresence intercept num-consumer --port 8081:8081
telepresence list
telepresence leave num-consumer这是命令行

goland中怎么用
1、teleprencese service类型：deployment，ReplicaSet, StatefulSet ArgoRollout
拦截这四种，调试手段是同一套：流量转到本机，GoLand 里断点、改局部变量。差别不在「能多干一件神奇的事」，而在 集群里真正管这些 Pod 的是哪一种对象，你就拦哪一种。
OSS 常见四种：
拦什么	典型场景 拦截之后实际发生的事
Deployment
无状态 Web / API
进这个 Deployment 对应 Service 的请求转到本机
ReplicaSet
独立 RS，没有上层 Deployment
只劫持 这一个 RS 管的 Pod
StatefulSet
有序号、稳身份的服务（DB、消息队列、主从）
劫持这个 STS 对应 Service 的流量
Argo Rollout
金丝雀 / 蓝绿发布
拦正在灰度的那份工作负载

2、--mount false 参数，不挂盘，只拦流量

3、port的，如果只有一个 Service ，那么可以只写个8080，然后会知道，如果有多个svc 8080:30080/TCP,9090:30090/TCP，就要填8080:XXX,XXX可以是http，gprc这种形式，也可以是30880这种端口形式.
本机 8081（冒号左边）
GoLand 的 Run/Debug 配置里加环境变量：PORT=8081 程序才会 ListenAndServe(":8081")。集群 8082（冒号右边）k8s/hello-tp.yaml 里和“容器在听哪个端口”有关的都要改成 8082：

4、前端拦截 也可以
5、teleprencese还有很多额外参数
```

```
1、kubelet 和 helm的安装，可以在外网用goland下载，然后将可执行文件弄到内网，然后放在path下就可以了
2、teleprence下载下来，解压出可执行文件，然后放到path
另外由于teleprensece连接的时候需要在目标机器安装traffic-manager组件，因此会拉镜像，可以在外网打包根据对应的teleprensece version版本下载对应的包，拷到内网导入，如ghcr.io/telepresenceio/tel2:2.26.2,traffic应该会在这个ambassador namespace下

如果中途失败，本机手动执行telepresence helm install --kubeconfig <45的kubeconfig> --manager-namespace ambassador

3、
```





## 组件

可用 正确 性能 

Kubernetes 自身
│
├── 1. 可用性 Availability
├── 2. 正确性 Correctness
├── 3. 性能 Performance
├── 4. 安全 Security
├── 5. 可维护/可观测/可演进 Operability
└── 6. 成本/资源效率/可扩展性 Efficiency & Scalability

# 原理1

# 1. 调度 Scheduling

## 它到底在解决什么？

一句话：

> **当系统有很多任务、很多资源时，决定“哪个任务应该放到哪个资源上运行”。**

最简单：

```
任务：
Pod A
Pod B
Pod C

机器：
Node 1
Node 2
Node 3
```

调度器要决定：

```
Pod A → Node 2
Pod B → Node 1
Pod C → Node 3
```

这看起来很简单。

但真实系统马上复杂起来。

假设：

```
Node 1:
CPU 16
Memory 64G
GPU 0

Node 2:
CPU 32
Memory 128G
GPU 4

Node 3:
CPU 8
Memory 32G
GPU 1
```

Pod 要求：

```
CPU 8
Memory 16G
GPU 1
```

那至少有：

```
Node 2
Node 3
```

可以运行。

现在还有：

- Pod affinity亲和力
- anti-affinity
- taint/toleration污点/耐受性
- topology拓扑
- priority优先事项
- gang帮派 scheduling
- NUMA
- GPU topology
- bin packing
- resource fragmentation碎片化

问题马上从：

> “放哪里？”

变成：

> **“在一组约束和目标函数下，怎样分配资源，使整个系统最优？”**

这就是调度。

------

## Kubernetes Scheduler 在干什么？

可以粗略理解为：

```
Pod
 ↓
Filter
 ↓
哪些 Node 能放？
 ↓
Score
 ↓
哪个 Node 更合适？
 ↓
Bind
```

例如：

```
Node A ❌ GPU不足
Node B ✅
Node C ✅
Node D ❌ taint
```

然后：

```
B score = 70
C score = 85

Pod → Node C
```

------

## 更深一层

真正的调度问题往往属于：

> **Optimization Problem**

你需要同时考虑：

```
资源利用率
公平性
延迟
成本
可靠性
拓扑
数据本地性
SLA
```

比如 AI Infra：

```
Model A
Model B
Model C

GPU:
A100 × 8
H100 × 8
L40S × 16
```

你不能简单：

> 谁空闲就给谁。

因为：

- 模型大小不同
- 显存需求不同
- GPU 算力不同
- GPU 间互联不同
- 数据位置不同

这时调度就变成了非常大的系统问题。

### 自测

不用看资料，回答：

> ① Kubernetes Scheduler 为什么不能简单地“找到 CPU/Memory 满足要求的第一个 Node 就结束”？
>
> 找第一个满足 CPU/Memory request 的节点，会忽略：
>
> - 约束：节点亲和/反亲和、taint/toleration、拓扑分布、GPU、端口、卷的可挂载性等。
> - 故障域：同一服务副本若都落在一个节点或一个可用区，单点故障会同时带走它们。
> - 资源碎片：随意把小 Pod 放进大节点，可能让后续需要更多 CPU、内存或 GPU 的 Pod 无处可去。
> - 负载与成本：节点虽“还有资源”，但可能已经很忙；也可能应优先填满现有节点以避免额外扩容。
> - 策略与优先级：高优先级 Pod 可能需要抢占，租户配额和公平性也要参与决策。
>
> 所以 Kubernetes Scheduler 通常先 Filter 掉不能运行该 Pod 的节点，再对剩余节点 Score，最后选择综合最合适的节点并 Bind。
>
>  ② Pod anti-affinity 是在解决什么问题？
>
> Pod anti-affinity 用来避免某些 Pod 被调度到同一类位置，降低“同时挂掉”的风险。同一个服务有 3 个副本。若它们都在同一台 Node，上面那台机器宕机时服务会全灭。配置 anti-affinity 后，Scheduler 会尽量或强制把副本分散到不同 Node、可用区或机架。
>
>  ③ GPU 调度和普通 CPU 调度最大的区别是什么？
>
> GPU 更稀缺、更异构，而且往往必须以“整块设备 + 正确拓扑”来分配。
>
> CPU 调度通常主要看可分配核数、内存及负载；少量碎片一般还能接受。GPU 调度还要考虑：
>
> - 型号与能力：A100、H100、L4 的显存、算力、特性不同。
> - 不可随意超卖：普通 GPU 通常按整卡分配；一张空闲一半的卡不一定能给另一个 Pod 使用。
> - 拓扑：多卡训练需要 NVLink、PCIe、NUMA 等连接关系合适，否则通信会成为瓶颈。
> - 碎片影响更大：例如任务需要 4 张互联良好的 GPU，集群虽然总共有 4 张空卡，但分散在不同节点或拓扑不合适，仍无法运行。
> - 设备与驱动约束：还涉及 NVIDIA Device Plugin、驱动/CUDA 版本、MIG 切分等。
>
>  ④ 什么叫 resource fragmentation？
>
> Resource fragmentation（资源碎片化）是指：集群总资源看起来还够，但它们分散在不同节点、不同设备或不合适的组合中，导致某个 Pod 实际上无法被调度。
>
> 例如有两台节点：
>
> - Node A：剩余 6 CPU、2 GB 内存
> - Node B：剩余 2 CPU、10 GB 内存
>
> 新 Pod 需要 `4 CPU + 8 GB 内存`。集群总剩余是 8 CPU、12 GB 内存，数值上足够；但没有一台 Node 同时满足，因此 Pod 仍然 Pending。
>
> GPU 中更明显：总共剩 4 张 GPU，但它们分散在不同节点；一个要求同机 4 卡训练的任务依旧无法运行。

------

# 2. 资源管理 Resource Management

这个和调度非常接近，但不是一回事。

调度回答：

> **谁去哪里。**

资源管理回答：

> **系统到底有多少资源、谁可以用多少、如何分配、如何限制。**

------

## 最简单的模型

机器：

```
Node
CPU = 32
Memory = 128GB
GPU = 4
```

Pod A：

```
requests:
CPU 4
Memory 8GB

limits:
CPU 8
Memory 16GB
```

这里已经出现：

### Request

> 我至少需要这么多。

### Limit

> 最多允许我使用这么多。

------

## Kubernetes 的资源管理

主要包括：

```
CPU
Memory
Ephemeral Storage
GPU
HugePages
```

再往上：

```
ResourceQuota
LimitRange
PriorityClass
QoS
```

例如：

```
Namespace A
quota:
CPU = 100
Memory = 200GB

Namespace B
quota:
CPU = 50
Memory = 100GB
```

这实际上已经是在做：

> **Multi-tenant Resource Management**

------

## 更深的资源管理

真正的大型系统还会出现：

```
Overcommit
Preemption
Fairness
Quota
Isolation
Admission Control
Backpressure
Load Shedding
```

例如：

服务器只有：

```
100 CPU
```

但是所有 Pod 的 requests 加起来：

```
150 CPU
```

怎么办？

这就是：

> **Overcommit / Admission / Scheduling / Runtime behavior**

------

## 自测

> ① request 和 limit 的本质区别是什么？
>  ② Kubernetes 为什么允许 requests 总和和节点真实资源之间出现一些复杂的“超额/实际使用”情况？
>  ③ OOMKilled 是谁决定的？Kubernetes、kubelet、Linux kernel，分别扮演什么角色？
>  ④ ResourceQuota 和 Scheduler 分别解决什么问题？
>
> ① `request` 是调度承诺：Scheduler 依据它判断节点是否“放得下”Pod。
> `limit` 是运行上限：容器最多可用多少资源。
>
> - CPU 超过 `limit`：通常被 cgroup 限流（throttling）。
> - 内存超过 `limit`：通常会触发 cgroup OOM，进程可能被杀。
> - CPU 没设 limit 时可短时多用；内存没有天然“可安全超用”的空间。
>
> ② 因为 Kubernetes 需要在利用率、隔离性和可用性之间权衡。
>
> `requests` 是容量预留，而真实使用量会随时变化；很多工作负载不会持续跑满。允许一定程度的 overcommit，可以提高集群利用率。但这意味着高峰期可能资源竞争、CPU 被限流，或内存压力下发生驱逐／OOM。调度时主要看 request，不是实时使用率。
>
> ③ 最终执行 OOM kill 的是 Linux kernel；Kubernetes 并不直接杀进程。
>
> - Kubernetes：通过 Pod 的 `resources.limits.memory` 和 QoS 等声明期望与策略。
> - kubelet：把这些声明配置进容器运行时/cgroup；在节点内存压力时，也可能主动驱逐 Pod（Eviction）。
> - Linux kernel：检测到 cgroup 达到内存上限，或节点发生全局内存耗尽时，选择并杀掉进程。容器状态因此显示为 `OOMKilled`。
>
> ④ `ResourceQuota` 管“一个 Namespace 最多能申请/使用多少”，属于多租户的准入和治理边界；例如限制团队 A 的 CPU request 总量不超过 100 核。
>
> Scheduler 管“一个具体 Pod 应放到哪台 Node”，根据 request、可用资源、亲和性、污点、拓扑等做筛选与打分。
>
> 简化为：Quota 防止某个租户把集群额度占光；Scheduler 在已获准的资源需求下寻找最合适的落点。



------

# 3. 服务发现 Service Discovery

这是：

> **“我想调用 user-service，但我不知道它现在在哪。”**

因为在云原生环境里：

```
Pod IP
```

非常不稳定。

今天：

```
10.0.1.15
```

明天 Pod 没了：

```
10.0.2.31
```

所以不能让客户端写死：

```
10.0.1.15:8080
```

------

## Kubernetes 的办法

```
Client
  ↓
Service
  ↓
Endpoint / EndpointSlice
  ↓
Pod
```

例如：

```
user-service
```

背后可能有：

```
Pod A
Pod B
Pod C
```

Service 提供一个稳定的虚拟入口。

客户端只需要知道：

```
user-service:8080
```

------

## 再往深

服务发现实际上包括：

```
注册
发现
健康检查
地址更新
负载均衡
故障摘除
DNS
```

例如：

```
Pod A healthy
Pod B healthy
Pod C dead
```

服务发现系统需要让调用方逐渐知道：

```
A ✅
B ✅
C ❌
```

这里已经开始进入：

> **动态成员管理**

------

## 自测

> ① K8s Service 的 IP 为什么可以稳定，而 Pod IP 会变化？
>  ② Service、EndpointSlice、kube-proxy 三者之间大概是什么关系？
>  ③ 如果一个 Pod 已经“看起来健康”，但实际上处理请求很慢，服务发现为什么可能仍然把流量发给它？
>  ④ DNS 服务发现和服务注册中心相比，各自有什么优缺点？
>
> ① Pod 是可替换实例，重建/漂移到新节点后通常获得新 IP；Service 是 API 对象，创建后其 ClusterIP 保持不变，除非删除重建或改为 headless。它提供稳定入口，后端 Pod 可随时变化。
>
> ② 关系是：
>
> ```
> 客户端 → Service（稳定 VIP/DNS）
>              ↓ selector
>       EndpointSlice（当前可用 Pod IP:Port 列表）
>              ↓
>  kube-proxy（将发往 Service IP 的流量转发/NAT 到某个 Endpoint）
> ```
>
> Service 定义“访问谁”；EndpointSlice 记录“当前有哪些后端”；kube-proxy 监听它们的变化，并通过 iptables/IPVS（或在部分集群中由 eBPF 替代）实现实际转发。
>
> ③ “健康”往往只表示 readiness probe 成功，不代表每个请求都快。探针可能只检查 HTTP 200、端口连通或一个很轻的健康接口；它未覆盖慢 SQL、下游依赖卡顿、连接池耗尽、GC 停顿、CPU 抢占等。因此该 Pod 仍留在 EndpointSlice，继续接收流量。要改善，需要把真正的服务能力纳入 readiness、设置合理 timeout，并用限流、熔断和负载保护避免拖垮整体。
>
> ④ DNS 服务发现：
>
> - 优点：简单、通用、客户端零或低侵入；名称稳定，天然适合 K8s Service。
> - 缺点：DNS TTL/客户端缓存可能导致更新不够实时；通常只提供“名称 → 地址”，健康检查、权重、灰度和复杂路由能力较弱。
>
> 注册中心（如 Consul、Nacos、Eureka）：
>
> - 优点：实例注册、健康检查、动态上下线、权重、元数据、治理策略通常更丰富。
> - 缺点：要运行和维护额外控制面；客户端/SDK 耦合更深，故障模式和一致性问题也更多。

------

# 4. 故障恢复 Fault Recovery

这个是最容易被低估的。

正常系统：

```
Request
 ↓
Service
 ↓
Pod
 ↓
DB
```

大家都会设计。

真正难的是：

```
Request
 ↓
Service
 ↓
Pod
 ↓
DB
     ↑
    挂了
```

或者：

```
Pod
 ↓
Node
 ↓
Network
```

突然断。

------

## 故障恢复主要在解决：

> **失败发生以后，系统如何重新回到一个可工作的状态。**

常见机制：

```
Restart
Retry
Timeout
Failover
Replication
Checkpoint
Recovery
Reconciliation
Leader Election
```

K8s 很典型：

```
Desired:
replicas = 3

Actual:
replicas = 2

Controller:
发现不一致

→ 创建 Pod
```

这就是：

> **Reconciliation**

------

## 这里有个特别重要的问题

“自动恢复”并不等于“正确恢复”。

比如：

```
Pod A
 ↓
扣库存
 ↓
请求超时
```

你直接 Retry：

```
Retry
 ↓
再扣一次
```

可能就出事了。

所以：

> **Recovery 必须建立在正确性模型上。**

这就是为什么：

```
Retry
```

永远不能孤立学习。

它必须和：

```
Timeout
Idempotency
Consistency
Transaction
```

一起理解。

### 自测

> ① K8s 为什么能够很好地恢复“无状态服务”，但不能自动恢复业务语义？
>  ② Retry 为什么可能把故障放大？
>  ③ Crash recovery 和 failover 有什么区别？
>  ④ 为什么 reconciliation 是一种非常重要的故障恢复思想？
>
> ① Kubernetes 能恢复的是“运行状态”，不是“业务是否正确完成”。
>
> 无状态服务的副本可被任意重建：Pod 挂了，Deployment 发现副本数少于期望值，就创建新 Pod；请求可由其他副本接手。
>
> 但业务语义涉及外部状态与副作用，例如“扣款是否成功”“订单是否已创建”“消息是否已消费”。K8s 无法仅凭 Pod 重启判断操作执行到哪一步，更不能安全地决定是否重做。这里需要应用自行保证幂等、事务、去重、消息确认和补偿。
>
> ② Retry 会在故障时额外制造请求，形成正反馈：
>
> ```
> 下游变慢 → 超时增多 → 上游重试 → 下游负载更高 → 更慢/更多超时
> ```
>
> 若多层服务都重试，流量可能成倍放大；如果操作非幂等，还会造成重复扣款、重复发货等副作用。应配合超时、有限重试、指数退避、jitter、熔断和幂等键。
>
> ③ Crash recovery 是“同一个实例崩了，如何让它恢复运行”，例如容器退出后 kubelet 重启容器、进程读取 WAL 恢复。
>
> Failover 是“原实例或主节点不可用时，如何切换到另一个可承担职责的实例”，例如数据库 Primary 挂掉后提升 Replica 为新 Primary，并让流量切过去。
>
> ④ Reconciliation 的核心是持续比较“期望状态”和“实际状态”，发现偏差就尝试把实际状态拉回期望状态：
>
> ```
> 期望：3 个健康副本
> 实际：2 个健康副本
> 控制器：创建/修复一个副本
> ```
>
> 它不依赖一次性、脆弱的“出错事件必须刚好被捕获”；控制器可反复执行，短暂失败后继续收敛。这使 K8s 能处理 Pod 消失、节点短暂异常、控制器重启和 API 调用失败等情况。前提是操作尽量幂等，并接受最终收敛而非瞬时完成。

------

# 5. 控制平面 Control Plane

这个和 K8s 非常相关。

你可以先理解成：

> **负责“决定系统应该变成什么样”的那一层。**

而真正干活的：

> **Data Plane**

------

## Kubernetes

大致：

```
Control Plane

API Server
Scheduler
Controller Manager
etcd

        ↓

Data Plane

kubelet
container runtime
Pod
network
storage
```

控制平面：

```
“我要 3 个 Pod”
```

数据平面：

```
“这里真的跑着 3 个 Pod”
```

------

## 为什么要区分？

因为：

> **控制决策和实际执行是两类完全不同的问题。**

比如：

```
Controller：

desired replicas = 3
actual replicas = 2

→ create Pod
```

Controller 不亲自运行 Pod。

它只是：

> 改变系统状态。

------

## 更深一层

很多系统都有：

```
Control Plane
Data Plane
```

例如：

- Kubernetes
- SDN
- Service Mesh
- 云平台
- 网络设备
- AI 集群管理

典型结构：

```
Control Plane
 ↓
Policy / State / Decision
 ↓
Data Plane
 ↓
Actual Execution
```

------

## 自测

> ① kube-scheduler 算 control plane 还是 data plane？为什么？
>  ② kubelet 属于哪一层？
>  ③ etcd 在控制平面中的作用是什么？
>  ④ 为什么控制平面通常不能直接承担全部数据流量？
>
> ① `kube-scheduler` 属于 control plane。它不处理业务请求，而是根据资源、约束与策略决定“新 Pod 应该绑定到哪个 Node”，即修改集群的期望/分配状态。
>
> ② `kubelet` 通常归入 data plane / node plane。它运行在每个 Node 上，接收 PodSpec 后通过容器运行时、CNI、CSI 等真正把 Pod 跑起来并持续汇报状态。它是控制面决策的执行者。
>
> ③ `etcd` 是控制平面的强一致状态存储。API Server 将 Kubernetes 对象的期望状态和大量实际状态写入其中，例如 Pod、Deployment、Service、ConfigMap、Lease。Scheduler、Controller Manager 等通过 API Server 读取和更新这些状态来协作；etcd 不直接承载业务流量。
>
> ④ 控制平面重在决策、协调和全局状态一致性，通常吞吐较低、写入更昂贵，也必须稳定可靠。若它直接承载所有业务数据流量：
>
> - 会被高 QPS、长连接和大流量压垮，影响调度与故障恢复；
> - 数据路径与控制路径耦合，业务流量异常会拖垮集群管理能力；
> - 其一致性与权限边界会成为性能瓶颈和更大的攻击面。
>
> 因此业务流量应由 Service、Ingress/Gateway、kube-proxy/eBPF、CNI 和工作负载等数据面组件处理；控制面只负责声明、调度、配置与协调。

------

# 6. 一致性 Consistency

这个是八个里面非常值得你认真学的。

因为“一致”这个词非常容易被误解。

------

## 最简单

两个节点：

```
A = 100
B = 100
```

A 修改：

```
A = 200
```

现在：

```
A = 200
B = 100
```

那么：

> 系统是不是不一致了？

**不一定。**

这取决于你的 consistency model。

------

## 强一致性

希望：

> 写入成功以后，其他读取者立即看到符合规定的新状态。

典型：

```
Write
 ↓
Commit
 ↓
Read
 ↓
New Value
```

但强一致通常需要付出：

```
Latency
Coordination
Availability
```

------

## 最终一致

允许：

```
A = 200
B = 100
```

暂时存在。

只要最终：

```
A = 200
B = 200
```

就可以。

这就是：

> Eventual Consistency

------

## 再往深

这里会进入：

```
CAP
PACELC
Linearizability
Sequential Consistency
Causal Consistency
Eventual Consistency
Quorum
Consensus
Replication
```

而且一个非常关键的认知：

> **“一致性”不是一个二元开关。**

不是：

```
一致 / 不一致
```

而是一个**一致性模型**。

------

## 自测

> ① “最终一致”为什么不是“系统最终一定正确”？
>  ② CAP 中的 C 到底是什么意思？
>  ③ quorum 是什么？
>  ④ 为什么数据库事务的 ACID 和分布式系统里的 consistency 不是完全一回事？
>
> ① 最终一致（eventual consistency）只承诺：在停止新写入、没有持续故障且副本最终能通信的前提下，副本会收敛到同一状态。它不保证这个状态一定符合业务正确性。
>
> 例如两个重复的“扣款”事件被都正确复制到了所有副本，系统最终一致，但余额仍可能因缺少幂等控制而错误。业务正确性还需要校验、约束、幂等、事务、补偿等机制。
>
> ② CAP 的 C 是 Consistency，通常指线性一致性（linearizability）：一次写入成功后，任何后续读取都不能读到更旧的值；所有客户端观察到的效果像操作在一个全局单一顺序中发生。
>
> 它不是泛泛的“数据没有错”，也不是 ACID 中的 C。
>
> ③ Quorum（法定人数/多数派）是一次决策或提交至少需要多少个节点同意。常见多数派为：
>
> ```
> N 个副本 → quorum = floor(N / 2) + 1
> 3 个节点 → 2
> 5 个节点 → 3
> ```
>
> 多数派之间必然至少有一个共同节点，因此可减少两个分区各自独立作出冲突决定的风险。Raft 常用 quorum 选主、提交日志；但不同系统的 quorum 语义不同，例如 Redis Sentinel 的 quorum 主要用于故障判断。
>
> ④ ACID 的 Consistency 是数据库事务从一个满足数据库约束的状态，转换到另一个满足约束的状态，例如外键、唯一约束、余额不能为负等。
>
> 分布式系统中的 consistency 则是“多个副本/客户端对读写顺序与可见性的保证强度”，如线性一致、顺序一致、因果一致、最终一致。
>
> 所以：
>
> ```
> ACID C：状态是否符合业务/数据库约束
> 分布式 C：不同节点看到的数据与操作顺序是否一致
> ```
>
> 二者可以相关，但互不替代。

------

# 7. 网络 Networking

这是 Infrastructure 最容易形成巨大壁垒的地方。

因为很多人会：

```
ip addr
ping
curl
```

但不知道：

> 请求到底经历了什么。

------

一个 K8s 请求：

```
Client
 ↓
DNS
 ↓
Service
 ↓
kube-proxy / eBPF
 ↓
CNI
 ↓
Node Network
 ↓
Pod Network Namespace
 ↓
Application
```

如果跨 Node：

```
Pod A
 ↓
veth
 ↓
CNI
 ↓
Node 1
 ↓
Overlay / Routing
 ↓
Node 2
 ↓
CNI
 ↓
Pod B
```

每一层都可能有问题。

------

## 你至少应该理解：

```
L2
L3
L4
L7
```

以及：

```
TCP
UDP
HTTP
DNS
TLS
NAT
Routing
Load Balancing
```

在 K8s：

```
Service
Ingress
CNI
NetworkPolicy
DNS
```

------

## 更深一层

开始研究：

```
Connection Tracking
NAT
Packet Loss
MTU
TCP Congestion
Latency
Head-of-Line Blocking
Connection Pool
Keepalive
```

这时候你遇到性能问题，才有能力定位：

> 到底是应用慢、网络慢，还是连接模型有问题。

------

## 自测

> ① Pod 到 Pod 的网络，为什么一般不应该简单理解为“两个 Pod 之间直接插一根网线”？
>  ② Service IP 是不是一个真正绑定在某个网卡上的普通 IP？
>  ③ TCP 三次握手和 HTTP 请求之间是什么关系？
>  ④ MTU 为什么可能导致一些“只有大包才出问题”的诡异故障？
>
> ① Pod 通信通常会经过多个网络层：Pod 的 network namespace、veth、CNI 网桥或路由、Node 网络、跨节点 overlay/路由、目标 Node 的 CNI，最后才到目标 Pod。中间还可能有 NAT、NetworkPolicy、iptables/IPVS/eBPF、conntrack。因此它不是一条固定的“直连网线”，路径、封装和策略都可能影响通信。
>
> ② 一般不是。`ClusterIP` 是逻辑虚拟 IP，通常不配置在某块网卡上。`kube-proxy` 通过 iptables 或 IPVS（有些集群是 eBPF 数据面）把访问该 IP:Port 的流量重写/转发到真实 Pod Endpoint。它更像一个稳定的流量匹配规则，而非某台机器真实持有的地址。
>
> ③ TCP 是 HTTP 常见的传输基础。通常流程是：
>
> ```
> TCP 三次握手 → 建立可靠连接 → 在该连接中发送 HTTP 请求/响应
> ```
>
> HTTP/1.1 可在同一 TCP 连接上复用多个顺序请求；HTTP/2 可在同一连接上多路复用多个流。若已有 keep-alive 连接，新的 HTTP 请求不需要每次重新三次握手。HTTP/3 则使用基于 UDP 的 QUIC，不走 TCP 三次握手。
>
> ④ MTU 是单个链路可承载的最大 IP 包大小。跨节点 overlay 会增加封装头，实际可用 MTU 变小。
>
> 小包能通过，但大包若超过路径 MTU：
>
> - 若允许分片或 PMTUD 正常，可能被分片或发送方自动缩小；
> - 若中间设备丢弃分片，或 ICMP “Fragmentation Needed / Packet Too Big” 被防火墙拦截，发送方不知道该缩小包；
> - 大包就会反复重传、超时或卡住。
>
> 于是健康检查、小请求、`ping` 可能正常，上传文件、TLS 握手中的较大证书、较大响应或数据库复制却失败。这就是典型的 MTU black hole。

------

# 8. 存储 Storage

K8s 用户经常把存储理解成：

```
PV
PVC
StorageClass
```

这是：

> **存储编排层。**

但 Storage 本身远远更大。

------

## 先区分几个东西

```
Application
 ↓
Filesystem
 ↓
Block Device
 ↓
Storage System
 ↓
Physical Disk
```

例如：

```
Pod
 ↓
PVC
 ↓
CSI
 ↓
Ceph RBD
 ↓
OSD
 ↓
Disk
```

或者：

```
Pod
 ↓
PVC
 ↓
CSI
 ↓
EBS
```

------

## 存储真正关心：

```
Durability
Availability
Consistency
Replication
IOPS
Throughput
Latency
Failure Recovery
Snapshot
Backup
Locality
```

例如：

```
1TB 数据

Disk A
Disk B
Disk C
```

复制三份。

那么：

> 一块盘挂掉怎么办？

这只是最基本的问题。

------

## 再复杂一点

Ceph：

```
Client
 ↓
RBD
 ↓
Ceph Cluster
 ↓
OSD
 ↓
Disk
```

它内部自己就是一个复杂的分布式系统。

你之前碰过 CephFS / MinIO / HDFS，这实际上非常有价值。

------

## 自测

> ① PV/PVC 解决的是“存储本身”，还是“如何向工作负载提供存储”？
>  ② Block Storage 和 Object Storage 最核心的抽象差别是什么？
>  ③ 为什么数据库通常不能简单地直接放在任意共享文件系统上？
>  ④ Ceph 为什么本身也属于 Distributed Systems 的典型案例？
>
> ① PV/PVC 主要解决“如何把底层存储以可声明、可绑定、可挂载的方式提供给工作负载”，不是实现存储介质本身。
>
> - PV：集群中一份可供使用的存储资源抽象。
> - PVC：工作负载提出的容量、访问模式、性能等需求。
> - StorageClass/CSI：负责按需创建、挂载和回收实际后端存储。
>
> ② Block Storage 抽象为“可随机读写的裸磁盘块设备”：应用通常要自己创建文件系统、管理目录和文件，例如云盘、Ceph RBD。
>
> Object Storage 抽象为“通过 API 操作的对象”：每个对象有 key、内容和元数据，通常通过 HTTP/S3 API 读写；没有传统 POSIX 目录、块地址与原地随机写语义，例如 S3、Ceph RGW。
>
> ③ 数据库依赖很严格的持久化与并发语义：`fsync`/flush 是否真落盘、原子 rename、文件锁、缓存一致性、延迟抖动、断连后的行为等。任意共享文件系统未必正确实现或稳定提供这些语义；网络抖动、锁异常、缓存陈旧、性能长尾都可能导致数据库损坏或严重退化。
>
> 并非“共享文件系统一定不能用”。关键是数据库官方是否支持该文件系统，以及它是否能可靠满足所需 POSIX、持久化和性能保证。
>
> ④ Ceph 是典型分布式系统，因为它把数据和控制分布在多台节点上，并必须处理：
>
> - 数据分片与多副本/纠删码；
> - 节点、磁盘、网络故障后的检测、重平衡与恢复；
> - 一致性、写入确认与副本协调；
> - CRUSH 数据放置、容量变化与故障域；
> - 监控节点的 quorum/选主，以及客户端并发访问。
>
> 也就是说，Ceph 不只是“很多磁盘拼起来”，而是在不可靠节点与网络上持续提供一个统一、可恢复的存储系统。

                    Control Plane
                         │
            ┌────────────┼────────────┐
            ↓            ↓            ↓
        Scheduling    Resource      Service
                     Management     Discovery
            │            │            │
            └────────────┼────────────┘
                         ↓
                    Workloads
                         │
             ┌───────────┼───────────┐
             ↓           ↓           ↓
          Network      Storage      Compute
             │           │           │
             └───────────┼───────────┘
                         ↓
                    Fault Recovery
                         │
                         ↓
                    Consistency

