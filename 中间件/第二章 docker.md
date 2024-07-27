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

## docker配置各个中间件的参数

### mysql

**端口映射**

**环境变量**MYSQL_ROOT_PASSWORD