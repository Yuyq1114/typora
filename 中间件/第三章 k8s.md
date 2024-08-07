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
   kubectl exec -it POD_NAME -- /bin/bash
   ```

   - 在指定 Pod 内部启动一个交互式 Shell。

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