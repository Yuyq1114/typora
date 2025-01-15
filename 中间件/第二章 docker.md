## 基本概念介绍

### 1. 容器（Container）

- **容器**是 Docker 打包应用程序和其依赖项的标准单元。每个容器都是一个运行时环境，包括应用程序代码、运行时环境、系统工具、系统库等。
- 容器与虚拟机不同，它们不需要单独的操作系统实例，而是共享主机操作系统的内核，因此更加轻量级和快速。

### 2. 镜像（Image）

- **镜像**是容器的只读模板，用于创建容器实例。它包含应用程序运行所需的所有文件系统内容、配置和环境变量等。
- 镜像可以通过 Dockerfile 定义，其中包含构建容器的指令和配置。

### 3. Docker 引擎（Engine）

- **Docker 引擎**是 Docker 的核心组件，负责管理容器的生命周期和构建、运行、分发容器应用程序。
- Docker 引擎包括 Docker 客户端和 Docker 服务器，可以在单个主机上运行多个容器。

### 4. 仓库（Registry）

- **仓库**是用于存储和管理 Docker 镜像的地方。最常见的 Docker 仓库是 Docker Hub，它是一个公共的、云端的仓库，包含了大量的镜像供用户使用。
- 用户可以创建私有仓库或者使用其他公共仓库存储自己的镜像。

### 5. Dockerfile

- **Dockerfile** 是一个文本文件，包含了用于构建 Docker 镜像的指令和配置。通过 Dockerfile 可以定义镜像的内容、环境变量、运行命令等。
- 使用 Dockerfile 可以实现自动化地构建镜像，确保在不同环境中应用程序的一致性。

### 6. 容器编排（Orchestration）

- **容器编排**是管理和协调多个容器运行的过程。Docker 提供了多种容器编排工具，如 Docker Compose、Kubernetes 等，用于在集群中自动化部署、扩展和管理容器化应用程序。

### 7. Docker Compose

- **Docker Compose** 是一个用于定义和运行多容器 Docker 应用程序的工具。通过一个单独的 YAML 文件定义一组服务、网络和卷等配置，然后使用 `docker-compose` 命令启动、停止和管理这些服务。

### 8. 虚拟化与容器化的比较

- **虚拟化**（Virtualization）是通过在物理硬件上运行多个虚拟机实例来实现应用隔离和资源分配。每个虚拟机包含完整的操作系统和应用程序，需要更多的资源。
- **容器化**（Containerization）是利用容器技术实现应用程序和依赖项的隔离，容器共享主机操作系统内核，因此更轻量级、启动更快，并且更加适合于微服务架构和持续集成、持续部署（CI/CD）的场景。

### 示例操作流程

1. **编写 Dockerfile**：

   - 创建一个 Dockerfile 定义应用程序的环境和运行时依赖。

   ```
   # 使用官方 Python 镜像作为基础镜像
   FROM python:3.9-alpine
   
   # 设置工作目录
   WORKDIR /app
   
   # 复制当前目录下的所有文件到工作目录
   COPY . .
   
   # 安装应用程序依赖
   RUN pip install -r requirements.txt
   
   # 定义容器启动时执行的命令
   CMD ["python", "app.py"]
   ```

2. **构建镜像**：

   - 使用 `docker build` 命令根据 Dockerfile 构建镜像。

   ```
   docker build -t my-python-app .
   ```

3. **运行容器**：

   - 使用 `docker run` 命令运行镜像创建并启动容器实例。

   ```
   docker run -d --name my-app-container -p 8000:8000 my-python-app
   ```

4. **使用 Docker Compose 管理多容器应用**：

   - 编写 `docker-compose.yml` 文件定义多个服务，并使用 `docker-compose` 命令管理这些服务的生命周期。

   ```
   version: '3'
   services:
     web:
       build: .
       ports:
         - "8000:8000"
   ```

## Docker原理

### 1. 容器化技术基础

Docker 使用了 Linux 内核中的一些核心特性来实现容器化，主要包括：

- **命名空间（Namespaces）**：用于隔离应用程序的视图，包括进程、网络接口、文件系统等。Docker 使用多个命名空间来实现容器内部的隔离，如 PID 命名空间、网络命名空间、挂载命名空间等。
- **控制组（Cgroups）**：用于限制、账户和隔离资源（如 CPU、内存、磁盘 I/O）的使用。Docker 使用 Cgroups 来限制容器对主机资源的访问和使用，确保容器间的资源隔离和公平共享。

### 2. Docker 架构

Docker 的架构主要包括三个核心组件：

- **Docker Daemon**：也称为 Docker 服务端，负责管理 Docker 容器、镜像、网络和存储卷等。它监听 Docker API 请求，并管理容器的生命周期。
- **Docker Client**：也称为 Docker 客户端，通过 Docker 命令行接口（CLI）或者 API 向 Docker Daemon 发送命令和请求。用户通过 Docker 客户端与 Docker Daemon 交互，管理和操作容器和镜像等。
- **Docker Registry**：用于存储 Docker 镜像的仓库。最常见的是 Docker Hub，一个公共的 Docker Registry，也可以搭建私有的 Docker Registry 用于内部使用或安全性考虑。

### 3. 镜像（Image）和容器（Container）

- **镜像**：镜像是容器的只读模板，包含应用程序运行所需的所有文件系统内容、配置和环境变量等。镜像是通过 Dockerfile 定义的，Dockerfile 中包含了构建镜像的指令和配置。
- **容器**：容器是镜像的运行实例，每个容器都是一个独立的、轻量级的虚拟化环境。容器与虚拟机不同，它们共享宿主机的操作系统内核，并且启动更快、占用资源更少。

### 4. 容器生命周期管理

Docker 根据容器的生命周期提供了一些基本操作：

- **创建容器**：通过 `docker run` 命令创建并启动一个新的容器实例。
- **启动和停止容器**：使用 `docker start` 和 `docker stop` 命令分别启动和停止容器。
- **删除容器**：使用 `docker rm` 命令删除已停止的容器。
- **查看容器状态和日志**：通过 `docker ps` 和 `docker logs` 命令查看容器的状态和日志输出。

### 5. Docker 的优势和应用场景

- **环境一致性**：通过镜像，可以确保在不同环境中应用程序的一致性，避免了 "在我的机器上可以工作" 的问题。
- **快速部署和扩展**：Docker 容器可以快速部署和启动，支持水平扩展和自动化部署，适合于微服务架构和持续集成、持续部署（CI/CD）的实践。
- **资源利用效率**：与传统虚拟化技术相比，Docker 容器启动更快、占用资源更少，更加轻量级。
- **隔离性和安全性**：Docker 使用内核级别的隔离技术，提供了良好的安全性和资源隔离，使得多个容器可以安全地共存于同一主机上。

## 常用命令

### 容器管理命令

1. **创建并启动容器**

   ```
   docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
   ```

   - `docker run` 命令用于创建并启动一个新的容器实例。

   - 示例：

     ```
     docker run -d --name my-container nginx
     ```

2. **列出正在运行的容器**

   ```
   docker ps [OPTIONS]
   ```

   - `docker ps` 命令用于列出正在运行的容器。

   - 示例：

     ```
     docker ps
     ```

3. **列出所有容器（包括停止的）**

   ```
   docker ps -a
   ```

   - 使用 `-a` 参数可以列出所有的容器，包括正在运行和已停止的。

4. **停止容器**

   ```
   docker stop CONTAINER [CONTAINER...]
   ```

   - `docker stop` 命令用于停止一个或多个正在运行的容器。

   - 示例：

     ```
     docker stop my-container
     ```

5. **启动已停止的容器**

   ```
   docker start CONTAINER [CONTAINER...]
   ```

   - `docker start` 命令用于启动一个或多个已停止的容器。

   - 示例：

     ```
     docker start my-container
     ```

6. **删除容器**

   ```
   docker rm CONTAINER [CONTAINER...]
   ```

   - `docker rm` 命令用于删除一个或多个已停止的容器。

   - 示例：

     ```
     docker rm my-container
     ```

7. **查看容器日志**

   ```
   docker logs CONTAINER
   ```

   - `docker logs` 命令用于查看容器的日志输出。

   - 示例：

     ```
     docker logs my-container
     ```

8. **进入运行中的容器**

   ```
   docker exec [OPTIONS] CONTAINER COMMAND [ARG...]
   ```

   - `docker exec` 命令用于在运行中的容器内执行命令。

   - 示例：

     ```
     docker exec -it my-container bash
     ```

     这里的 

     ```
     -it
     ```

      参数用于交互式地进入容器内的 Bash Shell。

### 镜像管理命令

1. **列出本地镜像**

   ```
   docker images [OPTIONS]
   ```

   - `docker images` 命令用于列出本地所有的镜像。

   - 示例：

     ```
     docker images
     ```

2. **搜索镜像**

   ```
   docker search IMAGE_NAME
   ```

   - `docker search` 命令用于在 Docker Hub 上搜索镜像。

   - 示例：

     ```
     docker search nginx
     ```

3. **拉取镜像**

   ```
   docker pull IMAGE_NAME[:TAG]
   ```

   - `docker pull` 命令用于从远程仓库拉取镜像到本地。

   - 示例：

     ```
     docker pull nginx:latest
     ```

4. **删除镜像**

   ```
   docker rmi IMAGE_NAME[:TAG]
   ```

   - `docker rmi` 命令用于删除本地的一个或多个镜像。

   - 示例：

     ```
     docker rmi nginx:latest
     ```

### 其他常用命令

1. **查看 Docker 版本信息**

   ```
   docker version
   ```

   - `docker version` 命令用于查看 Docker 客户端和服务端的版本信息。

2. **查看 Docker 系统信息**

   ```
   docker info
   ```

   - `docker info` 命令用于查看 Docker 系统的详细信息，包括容器、镜像、存储等。

3. **管理 Docker 网络**

   ```
   docker network [SUBCOMMAND]
   ```

   - `docker network` 命令用于管理 Docker 网络，包括创建、连接、断开、移除网络等操作。

4. **管理 Docker 数据卷**

   ```
   docker volume [SUBCOMMAND]
   ```

   - `docker volume` 命令用于管理 Docker 数据卷，包括创建、删除、列出数据卷等操作。

##   实际命令

docker start $(docker ps -a | awk '{ print $1}' | tail -n +2)//开启

docker stop $(docker ps -a | awk '{ print $1}' | tail -n +2)//关闭所有容器

docker rm $(docker ps -a | awk '{ print $1}' | tail -n +2)//删除所有容器

docker-compose down //停止且删除

docker exec -it **podname** bash// 进入bash

redis-cli -c -h ip -p 7001 -a mypassword cluster nodes//查看集群状态

redis-cli -c -h 192.168.1.2 -p 6379 -a mypassword cluster nodes //获取集群结点

redis-cli -h 118.178.127.89 -p 7001 -a mypassword PING //ping结点

redis-cli -c -h 118.178.127.89 -p 7001 -a mypassword cluster nodes

118.178.127.89

## 安装docker

yum update -y

yum remove docker \
                docker-client \
                docker-client-latest \
                docker-common \
                docker-latest \
                docker-latest-logrotate \
                docker-logrotate \
                docker-engine

yum install -y yum-utils device-mapper-persistent-data lvm2

yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

yum install -y docker-ce docker-ce-cli containerd.io

systemctl start docker

systemctl enable docker

docker --version

//换源

vim /etc/docker/daemon.json

{
  "registry-mirrors": [
    "https://docker.registry.cyou",
    "https://docker-cf.registry.cyou",
    "https://dockercf.jsdelivr.fyi",
    "https://docker.jsdelivr.fyi",
    "https://dockertest.jsdelivr.fyi",
    "https://mirror.aliyuncs.com",
    "https://dockerproxy.com",
    "https://mirror.baidubce.com",
    "https://docker.m.daocloud.io",
    "https://docker.nju.edu.cn",
    "https://docker.mirrors.sjtug.sjtu.edu.cn",
    "https://docker.mirrors.ustc.edu.cn",
    "https://mirror.iscas.ac.cn",
    "https://docker.rainbond.cc"
  ]
}

systemctl restart docker

docker run hello-world



## 问题和解决方案（内部算法）

### 1. **容器隔离和安全性**

#### **问题描述**：

容器之间和容器与主机系统之间的隔离可能会受到攻击，影响系统的安全性。

#### **解决措施**：

- **命名空间（Namespaces）**：
  - Docker 使用 Linux 命名空间来提供进程隔离，包括进程 ID、网络、挂载和用户命名空间。这样可以确保容器之间的隔离。
- **控制组（Cgroups）**：
  - 控制组用于限制容器的资源使用（如 CPU、内存、I/O），防止容器过度消耗系统资源。
- **Seccomp 和 AppArmor**：
  - Docker 使用 Seccomp 配置文件来限制容器中可以使用的系统调用。
  - AppArmor（在支持的系统上）为容器提供额外的安全策略，控制容器的行为。
- **Docker Content Trust**：
  - Docker Content Trust（DCT）通过签名和验证镜像来确保镜像的来源和完整性，防止恶意镜像的使用。

### 2. **网络配置和连接**

#### **问题描述**：

容器网络可能遇到连接问题、网络冲突或配置复杂性。

#### **解决措施**：

- **网络模式**：
  - Docker 提供多种网络模式（如 bridge、host、overlay、macvlan）以满足不同的需求。合理选择网络模式可以解决大多数网络配置问题。
- **服务发现和负载均衡**：
  - 使用 Docker 内置的服务发现机制和负载均衡工具（如 Docker Swarm 的内置 DNS 和负载均衡器）来管理容器之间的通信。
- **自定义网络**：
  - 创建自定义网络并将容器连接到这些网络，使用网络驱动程序（如 bridge 或 overlay）来解决网络冲突和隔离问题。

### 3. **存储管理**

#### **问题描述**：

容器的持久化数据可能会面临存储不足、数据丢失或数据迁移问题。

#### **解决措施**：

- **数据卷（Volumes）**：
  - 使用 Docker 数据卷来管理容器数据，数据卷可以在容器之间共享，并持久化存储数据，即使容器删除也不会丢失。
- **绑定挂载（Bind Mounts）**：
  - 通过绑定挂载将主机文件系统的目录挂载到容器内，确保数据持久化和主机数据访问。
- **存储驱动**：
  - Docker 支持多种存储驱动（如 overlay2、aufs、btrfs），选择合适的存储驱动可以优化性能和兼容性。
- **数据备份和恢复**：
  - 定期备份数据卷，并在需要时恢复数据，以防止数据丢失。

### 4. **性能问题**

#### **问题描述**：

容器可能会遇到性能瓶颈，包括启动时间、资源使用率和 I/O 性能问题。

#### **解决措施**：

- **性能优化**：
  - 使用适当的资源限制（如 CPU 限制、内存限制）和优化容器配置来提高性能。
  - 使用 Docker 的 `--memory`、`--cpu` 和 `--blkio` 选项来限制容器的资源使用。
- **多阶段构建**：
  - 在 Dockerfile 中使用多阶段构建来减少镜像大小，提高构建效率。
- **镜像优化**：
  - 选择轻量级基础镜像，并清理不必要的文件和层，减少镜像大小，提高容器启动速度。

### 5. **版本兼容性和依赖管理**

#### **问题描述**：

不同版本的 Docker 或不同的容器镜像可能存在兼容性问题，导致运行失败或行为异常。

#### **解决措施**：

- **版本控制**：
  - 使用 Docker 的版本管理工具，确保 Docker 引擎和容器镜像的版本兼容性。
  - 在 Dockerfile 中指定基础镜像的版本，确保构建的一致性。
- **依赖管理**：
  - 在构建容器镜像时明确指定依赖项版本，确保应用在容器内的一致性和可重复性。
- **测试和验证**：
  - 在生产环境部署之前，充分测试容器镜像和应用，确保版本兼容性和稳定性。

### 6. **日志管理**

#### **问题描述**：

容器生成的日志可能难以集中管理、分析和存储。

#### **解决措施**：

- **集中日志管理**：
  - 配置 Docker 将容器日志输出到集中日志系统（如 ELK Stack、Fluentd、Graylog），便于日志管理和分析。
  - 使用 Docker 的日志驱动程序（如 json-file、syslog、journald）来配置日志记录和传输。
- **日志轮转和存储**：
  - 配置日志轮转和存储策略，防止日志文件过大，影响系统性能和存储空间。

### 7. **故障排除和调试**

#### **问题描述**：

容器故障可能难以排查和调试，特别是在复杂的应用环境中。

#### **解决措施**：

- **调试工具**：
  - 使用 Docker 提供的工具（如 `docker logs`、`docker exec`、`docker inspect`）来获取容器状态和运行时信息。
  - 利用 Docker Compose 的服务输出，调试和分析服务间的交互。
- **性能监控**：
  - 使用容器监控工具（如 Prometheus、Grafana、cAdvisor）来监控容器的性能指标，及时发现和解决性能问题。





## docker配置各个中间件的参数

### 换源

vim /etc/docker/daemon.json

```json
{
    "registry-mirrors" : [
    	"https://registry.docker-cn.com",
    	"http://hub-mirror.c.163.com",
    	"https://docker.mirrors.ustc.edu.cn",
    	"https://cr.console.aliyun.com",
    	"https://mirror.ccs.tencentyun.com"
  ]
}
```

重启

### 离线安装镜像

docker save -o dorisfe.tar apache/doris:doris-fe-2.1.7

docker save -o dorisbe.tar apache/doris:doris-be-2.1.7

echo %cd%

docker load -i dorisfe.tar





### mysql

**端口映射**3306

**环境变量**MYSQL_ROOT_PASSWORD

docker run -d --name mysql-container \
  -e MYSQL_ROOT_PASSWORD=mypassword \
  -e MYSQL_DATABASE=mydatabase \
  -e MYSQL_USER=myuser \
  -e MYSQL_PASSWORD=mypassword \
  -p 3306:3306 \
  mysql:8.0



### kafka

**使用kraft和apach的镜像**

docker run -d    --name broker -p 9092:9092   -e KAFKA_NODE_ID=1   -e KAFKA_PROCESS_ROLES=broker,controller   -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://:9093   -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://**外网ip**:9092   -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER   -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT   -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093   -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1   -e KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1   -e KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1   -e KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0   -e KAFKA_NUM_PARTITIONS=3   apache/kafka:latest



**使用zookeeper可能有问题**

docker run -d \
  --name zookeeper \
  -p 2181:2181 \
  zookeeper:latest



docker run -d \
  --name kafka \
  -p 9092:9092 \
  -e KAFKA_BROKER_ID=1 \
  -e KAFKA_ZOOKEEPER_CONNECT=120.46.80.186:2181 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://120.46.80.186:9092 \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
  wurstmeister/kafka:latest





**拉取apache/kafka:latest,不要拉取bitnami**

**端口映射**

Start a Kafka broker:

```console
docker run -d --name broker apache/kafka:latest
```

Open a shell in the broker container:

```console
docker exec --workdir /opt/kafka/bin/ -it broker sh
```

A *topic* is a logical grouping of events in Kafka. From inside the container, create a topic called `test-topic`:

```console
./kafka-topics.sh --bootstrap-server localhost:9092 --create --topic test-topic
```

Write two string events into the `test-topic` topic using the console producer that ships with Kafka:

```console
./kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test-topic
```

This command will wait for input at a `>` prompt. Enter `hello`, press `Enter`, then `world`, and press `Enter` again. Enter `Ctrl+C` to exit the console producer.

Now read the events in the `test-topic` topic from the beginning of the log:

```console
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test-topic --from-beginning
```

You will see the two strings that you previously produced:

```
hello
world
```

The consumer will continue to run until you exit out of it by entering `Ctrl+C`.

When you are finished, stop and remove the container by running the following command on your host machine:

```console
docker rm -f broker
```



### redis

goland连接时只需要给出ip和一个主节点的端口就可以，模式选cluster

**端口映射**6379

进入后redis-cli进入客户端

测试发送ping 收到pong



docker run -d \
  --name redis-node1 \
  -p 6379:6379 redis  \--requirepass mypassword

**集群test**

docker run -d --name redis-node-1 --net host --privileged=true  -v /data/redis/share/redis-node-1:/data redis:latest  --cluster-enabled yes  --appendonly yes --port 7001  --requirepass mypassword

docker run -d --name redis-node-2 --net host --privileged=true  -v /data/redis/share/redis-node-2:/data redis:latest  --cluster-enabled yes  --appendonly yes --port 7002  --requirepass mypassword

docker run -d --name redis-node-3 --net host --privileged=true  -v /data/redis/share/redis-node-3:/data redis:latest  --cluster-enabled yes  --appendonly yes --port 7003  --requirepass mypassword

docker run -d --name redis-node-4 --net host --privileged=true  -v /data/redis/share/redis-node-4:/data redis:latest  --cluster-enabled yes  --appendonly yes --port 7004  --requirepass mypassword

docker run -d --name redis-node-5 --net host --privileged=true  -v /data/redis/share/redis-node-5:/data redis:latest  --cluster-enabled yes  --appendonly yes --port 7005  --requirepass mypassword

docker run -d --name redis-node-6 --net host --privileged=true  -v /data/redis/share/redis-node-6:/data redis:latest  --cluster-enabled yes  --appendonly yes --port 7006  --requirepass mypassword

docker exec -it redis-node-1 /bin/bash

redis-cli --cluster create 117.50.85.130:7001 117.50.85.130:7002 117.50.85.130:7003 117.50.85.130:7004 117.50.85.130:7005 117.50.85.130:7006 --cluster-replicas 1 -a mypassword

redis-cli -h 127.0.0.1 -p 7001 -a mypassword



**集群支持**

docker network create --driver bridge redis-cluster-network

```
version: "3.3"
services:
  redis-node-1:
    image: redis:latest
    container_name: redis-node-1
    command: ["redis-server", "--requirepass", "mypassword", "--cluster-enabled", "yes", "--cluster-config-file", "nodes.conf", "--cluster-node-timeout", "5000", "--appendonly", "yes"]
    ports:
      - "7001:6379"
    volumes:
      - redis-node-1-data:/data
    networks:
      - redis-cluster-network

  redis-node-2:
    image: redis:latest
    container_name: redis-node-2
    command: ["redis-server", "--requirepass", "mypassword", "--cluster-enabled", "yes", "--cluster-config-file", "nodes.conf", "--cluster-node-timeout", "5000", "--appendonly", "yes"]
    ports:
      - "7002:6379"
    volumes:
      - redis-node-2-data:/data
    networks:
      - redis-cluster-network

  redis-node-3:
    image: redis:latest
    container_name: redis-node-3
    command: ["redis-server", "--requirepass", "mypassword", "--cluster-enabled", "yes", "--cluster-config-file", "nodes.conf", "--cluster-node-timeout", "5000", "--appendonly", "yes"]
    ports:
      - "7003:6379"
    volumes:
      - redis-node-3-data:/data
    networks:
      - redis-cluster-network

volumes:
  redis-node-1-data:
  redis-node-2-data:
  redis-node-3-data:

networks:
  redis-cluster-network:
    driver: bridge

```

docker-compose up -d

docker exec -it redis-node-1 redis-cli -a mypassword -h redis-node-2 -p 6379 ping



### postgreSQL

**端口映射**5432

**环境变量**POSTGRES_PASSWORD

docker run -d \
  --name my-postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydatabase \
  -p 5432:5432 \
  postgres



### doris

**//单结点模式**

https://github.com/apache/doris/blob/master/docker/runtime/docker-compose-demo/build-cluster/rum-command/3fe_3be.sh

sysctl -w vm.max_map_count=2000000

docker network create --driver bridge --subnet=172.20.80.0/24 doris-network
docker run -itd \
    --name=fe-01 \
    --env FE_SERVERS="fe1:172.20.80.2:9010,fe2:172.20.80.3:9010,fe3:172.20.80.4:9010" \
    --env FE_ID=1 \
    -p 8031:8030 \
    -p 9031:9030 \
    -v /data/fe-01/doris-meta:/opt/apache-doris/fe/doris-meta \
    -v /data/fe-01/log:/opt/apache-doris/fe/log \
    --network=doris-network \
    --ip=172.20.80.2 \
    apache/doris:doris-fe-2.1.7



docker run -itd \
    --name=be-01 \
    --env FE_SERVERS="fe1:172.20.80.2:9010,fe2:172.20.80.3:9010,fe3:172.20.80.4:9010" \
    --env BE_ADDR="172.20.80.5:9050" \
    -p 8041:8040 \
    -v /data/be-01/storage:/opt/apache-doris/be/storage \
    -v /data/be-01/log:/opt/apache-doris/be/log \
    --network=doris-network \
    --ip=172.20.80.5 \
    apache/doris:doris-be-2.1.7





**修改密码**

docker exec -it fe-01 bash

mysql -h 127.0.0.1 -P 9030 -u root -p

SET PASSWORD FOR 'root' = PASSWORD('mypassword');

SET PASSWORD FOR 'admin' = PASSWORD('mypassword');

**访问**

http://ip:8030/login



### MongoDB

**端口映射**27017



### rabbitMQ

**下载镜像选择-management**

**端口映射**    **5672 15672**

guest**只能本地**

访问地址
localhost:15672，这里的用户名和密码默认都是guest



### nacos

docker run -d \
  --name nacos \
  -e MODE=standalone \
  -p 8848:8848 \
  nacos/nacos-server:latest





http://120.26.84.229:8848/nacos/#/login
