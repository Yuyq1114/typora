![](D:\ProgramFile\typora\中间件\image\k8s_architecture.png)

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
  - 