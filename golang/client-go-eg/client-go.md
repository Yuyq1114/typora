配置
├── tools/clientcmd   用途：读取 kubeconfig最后生成*rest.Config
└── rest 这是 client-go 的底层连接层。读取 Pod 中自动挂载的：
ServiceAccount Token
CA 证书
KUBERNETES_SERVICE_HOST
KUBERNETES_SERVICE_PORT  创建底层 REST 客户端

客户端
├── kubernetes   提供操作 Kubernetes 内置资源的 Typed Clientset。
├── dynamic   通过 GVR 操作任意 Kubernetes 资源，尤其是 CRD。
├── metadata  只读取资源的元数据，不获取完整 `spec/status`。
└── discovery   询问 API Server

Controller 基础设施
├── informers 为 Kubernetes 内置资源创建 SharedInformer。
├── listers       Informer 的 Indexer 本地缓存
├── tools/cache     Controller 基础设施的核心包。平常使用生成好的 `informers` 即可；需要监听特殊 API 时才会直接用 `tools/cache` 构造。
└── util/workqueue   保存待 Reconcile 的 key，并处理重试、限速和去重。

Controller 增强能力
├── tools/leaderelection   多个 Controller 副本中只允许一个执行写操作。
├── tools/record  创建通过 `kubectl describe` 能看到的 Kubernetes Event。
├── util/retry   处理 Kubernetes 乐观锁冲突。
└── applyconfigurations   构建 Server-Side Apply 对象。

特殊操作
├── tools/remotecommand
├── tools/portforward
└── tools/watch    等待某个 Watch 条件满足。

测试
├── kubernetes/fake
├── dynamic/fake
└── testing



下面这些入口就够用了。

```
配置层
  clientcmd / rest
       ↓
客户端层
  kubernetes / dynamic / discovery
       ↓
监听缓存层
  informers / tools/cache / listers
       ↓
Controller 层
  workqueue / leaderelection / record
       ↓
高级操作
  applyconfigurations / remotecommand / portforward
```

官方仓库中还有大量由代码生成的 API Group 子包，不需要逐个背。[client-go 官方仓库](https://github.com/kubernetes/client-go)

# 一、配置与连接

## 1. `tools/clientcmd`

路径：

```
k8s.io/client-go/tools/clientcmd
```

作用：从 kubeconfig 加载集群地址、证书、用户和 context。

最常用入口：

```
config, err := clientcmd.BuildConfigFromFlags(
    "",
    "C:\\Users\\me\\.kube\\config",
)
```

返回：

```
*rest.Config
```

它适合在集群外运行的程序。

如果需要遵循默认 kubeconfig 查找和合并规则：

```
config, err :=
    clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
        clientcmd.NewDefaultClientConfigLoadingRules(),
        &clientcmd.ConfigOverrides{},
    ).ClientConfig()
```

### 应该记住

```
集群外运行
    ↓
clientcmd
    ↓
*rest.Config
```

------

## 2. `rest`

路径：

```
k8s.io/client-go/rest
```

作用：保存连接配置，并提供底层 REST 客户端。

最常用入口是集群内配置：

```
config, err := rest.InClusterConfig()
```

它会使用 Pod 中自动挂载的 ServiceAccount 信息。

常见入口：

```
rest.InClusterConfig()
rest.CopyConfig(config)
rest.HTTPClientFor(config)
rest.RESTClientFor(config)
```

底层 REST 请求可以这样写：

```
restClient.Get().
    Namespace("default").
    Resource("pods").
    Name("nginx").
    Do(ctx)
```

但普通业务代码一般不直接使用 `RESTClient`，而是使用上层的 Typed Client 或 Dynamic Client。[rest 包文档](https://pkg.go.dev/k8s.io/client-go/rest)

### 应该记住

```
rest.Config 是所有 Kubernetes 客户端的共同配置
```

常用字段包括：

```
config.Host
config.QPS
config.Burst
config.Timeout
config.UserAgent
config.Proxy
```

------

## 3. `transport`

路径：

```
k8s.io/client-go/transport
```

作用：处理认证、TLS、HTTP RoundTripper 等底层网络能力。

常见入口：

```
transport.New(config)
```

一般不会直接使用它，因为：

```
kubernetes.NewForConfig(config)
dynamic.NewForConfig(config)
```

内部已经处理了 transport。

------

# 二、Kubernetes API 客户端

## 4. `kubernetes`

路径：

```
k8s.io/client-go/kubernetes
```

这是操作 Kubernetes 内置资源最重要的包。

入口：

```
clientset, err := kubernetes.NewForConfig(config)
```

返回：

```
*kubernetes.Clientset
```

然后按照 API Group 进入：

```
clientset.CoreV1()
clientset.AppsV1()
clientset.BatchV1()
clientset.RbacV1()
clientset.NetworkingV1()
clientset.StorageV1()
clientset.AutoscalingV2()
```

接着选择资源：

```
clientset.CoreV1().Pods(namespace)
clientset.CoreV1().Services(namespace)
clientset.CoreV1().ConfigMaps(namespace)

clientset.AppsV1().Deployments(namespace)
clientset.AppsV1().StatefulSets(namespace)
clientset.AppsV1().DaemonSets(namespace)

clientset.BatchV1().Jobs(namespace)
clientset.BatchV1().CronJobs(namespace)
```

最后执行 CRUD：

```
Get(ctx, name, options)
List(ctx, options)
Watch(ctx, options)
Create(ctx, object, options)
Update(ctx, object, options)
Patch(ctx, name, patchType, data, options)
Delete(ctx, name, options)
```

完整调用：

```
deployment, err := clientset.
    AppsV1().
    Deployments("default").
    Get(ctx, "nginx", metav1.GetOptions{})
```

### 入口记忆方式

```
kubernetes.NewForConfig
    ↓
Clientset
    ↓
API Group
    ↓
Resource
    ↓
CRUD
```

------

## 5. `kubernetes/typed/...`

例如：

```
k8s.io/client-go/kubernetes/typed/apps/v1
k8s.io/client-go/kubernetes/typed/core/v1
k8s.io/client-go/kubernetes/typed/batch/v1
```

这些是生成出来的具体 API 客户端。

一般不会自己写：

```
appsv1client.NewForConfig(config)
```

通常从 Clientset 进入：

```
clientset.AppsV1()
```

但在需要缩小依赖或只使用某个 API Group 时，可以直接创建：

```
appsClient, err :=
    appsv1client.NewForConfig(config)
```

然后：

```
appsClient.Deployments("default").List(...)
```

------

## 6. `dynamic`

路径：

```
k8s.io/client-go/dynamic
```

作用：操作任意 Kubernetes 资源，尤其是 CRD。

入口：

```
client, err := dynamic.NewForConfig(config)
```

定义 GVR：

```
gvr := schema.GroupVersionResource{
    Group:    "demo.example.com",
    Version:  "v1",
    Resource: "webapps",
}
```

选择资源：

```
webApps := client.
    Resource(gvr).
    Namespace("default")
```

执行操作：

```
obj, err := webApps.Get(
    ctx,
    "demo",
    metav1.GetOptions{},
)
```

返回：

```
*unstructured.Unstructured
```

Dynamic Client 同样支持：

```
Get
List
Watch
Create
Update
Patch
Delete
Apply
```

### 应该记住

```
dynamic.NewForConfig
    ↓
Resource(GVR)
    ↓
Namespace
    ↓
CRUD
```

它没有：

```
client.WebAppsV1()
```

因为 Dynamic Client 在编译时不需要知道 WebApp 类型。

------

## 7. `metadata`

路径：

```
k8s.io/client-go/metadata
```

作用：只读取对象的元数据，不下载完整 spec/status。

入口：

```
client, err := metadata.NewForConfig(config)
```

使用：

```
obj, err := client.
    Resource(gvr).
    Namespace("default").
    Get(ctx, "demo", metav1.GetOptions{})
```

返回的信息主要是：

```
name
namespace
labels
annotations
UID
resourceVersion
ownerReferences
```

适合只关心元数据的大规模控制器或通用工具。

------

# 三、API 发现

## 8. `discovery`

路径：

```
k8s.io/client-go/discovery
```

作用：询问 API Server 当前支持哪些 Group、Version 和 Resource。

入口：

```
client, err :=
    discovery.NewDiscoveryClientForConfig(config)
```

常见方法：

```
client.ServerVersion()
client.ServerGroups()
client.ServerResourcesForGroupVersion("apps/v1")
client.ServerGroupsAndResources()
```

例如判断集群是否支持某个 CRD：

```
resources, err :=
    client.ServerResourcesForGroupVersion(
        "demo.example.com/v1",
    )
```

`kubectl api-resources` 底层思想就和 Discovery 有关。

------

## 9. `discovery/cached/memory`

路径：

```
k8s.io/client-go/discovery/cached/memory
```

Discovery 数据不会频繁变化，可以缓存。

入口：

```
cachedClient :=
    memory.NewMemCacheClient(discoveryClient)
```

还有磁盘缓存：

```
k8s.io/client-go/discovery/cached/disk
```

------

## 10. `restmapper`

路径：

```
k8s.io/client-go/restmapper
```

作用：在 GVK 和 GVR 之间转换。

例如：

```
GVK
Group: apps
Version: v1
Kind: Deployment
```

转换成：

```
GVR
Group: apps
Version: v1
Resource: deployments
```

典型入口：

```
mapper := restmapper.NewDeferredDiscoveryRESTMapper(
    cachedDiscoveryClient,
)
```

然后：

```
mapping, err := mapper.RESTMapping(
    schema.GroupKind{
        Group: "apps",
        Kind:  "Deployment",
    },
    "v1",
)
```

得到：

```
mapping.Resource
mapping.Scope
```

Dynamic Client 需要 GVR；当用户只给了 `kind: Deployment` 时，通常借助 Discovery 和 RESTMapper 找到 GVR。

------

# 四、Informer 和本地缓存

## 11. `informers`

路径：

```
k8s.io/client-go/informers
```

作用：为 Kubernetes 内置资源创建 Shared Informer。

主要入口：

```
factory :=
    informers.NewSharedInformerFactory(
        clientset,
        0,
    )
```

只监听指定 namespace：

```
factory :=
    informers.NewSharedInformerFactoryWithOptions(
        clientset,
        0,
        informers.WithNamespace("test1"),
    )
```

按 label 过滤：

```
factory :=
    informers.NewSharedInformerFactoryWithOptions(
        clientset,
        0,
        informers.WithTweakListOptions(
            func(options *metav1.ListOptions) {
                options.LabelSelector = "app=nginx"
            },
        ),
    )
```

取得具体 Informer：

```
podInformer :=
    factory.Core().V1().Pods()

deploymentInformer :=
    factory.Apps().V1().Deployments()
```

核心方法：

```
podInformer.Informer()
podInformer.Lister()

factory.Start(ctx.Done())
factory.WaitForCacheSync(ctx.Done())
factory.Shutdown()
```

[SharedInformerFactory 源码](https://github.com/kubernetes/client-go/blob/master/informers/factory.go)

------

## 12. `informers/<group>/<version>`

例如：

```
k8s.io/client-go/informers/core/v1
k8s.io/client-go/informers/apps/v1
k8s.io/client-go/informers/batch/v1
```

这些是生成的具体 Informer。

一般从 Factory 进入：

```
factory.Core().V1().Pods()
factory.Apps().V1().Deployments()
```

也可以直接创建：

```
podInformer := coreinformers.NewPodInformer(
    clientset,
    namespace,
    resyncPeriod,
    indexers,
)
```

但使用 SharedInformerFactory 更常见，因为多个模块可以共享同一套 Watch 和缓存。

------

## 13. `dynamic/dynamicinformer`

路径：

```
k8s.io/client-go/dynamic/dynamicinformer
```

作用：为 CRD 或任意 GVR 创建 Dynamic Informer。

入口：

```
factory :=
    dynamicinformer.NewDynamicSharedInformerFactory(
        dynamicClient,
        0,
    )
```

取得某个 GVR 的 Informer：

```
genericInformer := factory.ForResource(gvr)
```

然后：

```
genericInformer.Informer()
genericInformer.Lister()

factory.Start(ctx.Done())
```

对象类型是：

```
*unstructured.Unstructured
```

它内部就是通过 Dynamic Client 的 `List()` 和 `Watch()` 建立 Informer。[Dynamic Informer 源码](https://github.com/kubernetes/client-go/blob/master/dynamic/dynamicinformer/informer.go)

------

## 14. `tools/cache`

路径：

```
k8s.io/client-go/tools/cache
```

这是 Informer 机制的核心基础包。

最常见的类型和入口：

### 注册事件

```
cache.ResourceEventHandlerFuncs{
    AddFunc:    ...,
    UpdateFunc: ...,
    DeleteFunc: ...,
}
```

### 生成资源 key

```
key, err :=
    cache.MetaNamespaceKeyFunc(obj)
```

得到：

```
test1/nginx
```

支持删除 tombstone：

```
key, err :=
    cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
```

拆分 key：

```
namespace, name, err :=
    cache.SplitMetaNamespaceKey(key)
```

### 等待缓存同步

```
cache.WaitForCacheSync(
    ctx.Done(),
    informer.HasSynced,
)
```

### 手动创建 Informer

```
cache.NewSharedIndexInformer(...)
cache.NewSharedIndexInformerWithOptions(...)
```

### 底层组件

这个包中还能看到：

```
Reflector
DeltaFIFO
Indexer
Store
ListWatch
SharedIndexInformer
```

普通业务不会直接创建 Reflector 和 DeltaFIFO，但理解 Informer 源码时会遇到它们。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

------

# 五、Lister

## 15. `listers`

路径：

```
k8s.io/client-go/listers
```

下面也是按 API Group 生成的：

```
k8s.io/client-go/listers/core/v1
k8s.io/client-go/listers/apps/v1
k8s.io/client-go/listers/batch/v1
```

通常不手动创建，直接从 Informer 取得：

```
podLister :=
    factory.Core().V1().Pods().Lister()

deploymentLister :=
    factory.Apps().V1().Deployments().Lister()
```

使用：

```
pod, err := podLister.
    Pods("default").
    Get("nginx")

deployments, err := deploymentLister.
    Deployments("default").
    List(labels.Everything())
```

Lister 读取的是 Informer 的本地 Indexer，而不是 API Server。

如果手动创建，可以看到类似入口：

```
corelisters.NewPodLister(indexer)
appslisters.NewDeploymentLister(indexer)
```

日常代码应优先：

```
informer.Lister()
```

------

## 16. `dynamic/dynamiclister`

路径：

```
k8s.io/client-go/dynamic/dynamiclister
```

作用：从 Dynamic Informer 缓存读取 CRD。

通常也是通过：

```
genericInformer.Lister()
```

取得。

使用：

```
obj, err := lister.
    ByNamespace("test1").
    Get("webapp-demo")
```

返回的是通用对象，而不是具体的 WebApp struct。

------

# 六、工作队列

## 17. `util/workqueue`

路径：

```
k8s.io/client-go/util/workqueue
```

作用：将 Informer 事件和 Reconcile 处理分开。

现在推荐 typed WorkQueue：

```
queue :=
    workqueue.NewTypedRateLimitingQueue(
        workqueue.
            DefaultTypedControllerRateLimiter[string](),
    )
```

常用方法：

```
queue.Add(key)
queue.AddAfter(key, 5*time.Second)
queue.AddRateLimited(key)

key, shutdown := queue.Get()

queue.Done(key)
queue.Forget(key)
queue.NumRequeues(key)
queue.ShutDown()
```

典型处理：

```
key, shutdown := queue.Get()
if shutdown {
    return
}
defer queue.Done(key)

if err := reconcile(key); err != nil {
    queue.AddRateLimited(key)
    return
}

queue.Forget(key)
```

应记住三个主要队列层次：

```
TypedInterface
TypedDelayingInterface
TypedRateLimitingInterface
```

Controller 通常使用 `TypedRateLimitingInterface`。

------

# 七、Server-Side Apply

## 18. `applyconfigurations`

路径：

```
k8s.io/client-go/applyconfigurations
```

下面按 API Group 生成：

```
k8s.io/client-go/applyconfigurations/apps/v1
k8s.io/client-go/applyconfigurations/core/v1
```

入口一般是资源名构造函数：

```
deployment :=
    applyappsv1.Deployment(
        "nginx",
        "default",
    ).
        WithSpec(...)
```

然后通过 typed client 的 `Apply`：

```
result, err := clientset.
    AppsV1().
    Deployments("default").
    Apply(
        ctx,
        deployment,
        metav1.ApplyOptions{
            FieldManager: "my-controller",
        },
    )
```

Server-Side Apply 会记录字段所有者，适合声明式维护资源。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

------

# 八、Controller 的辅助包

## 19. `tools/leaderelection`

路径：

```
k8s.io/client-go/tools/leaderelection
```

作用：当 Controller 部署多个副本时，只让一个副本执行主要工作。

核心入口：

```
leaderelection.RunOrDie(
    ctx,
    leaderelection.LeaderElectionConfig{
        Lock: lock,

        Callbacks: leaderelection.LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                runController(ctx)
            },

            OnStoppedLeading: func() {
                os.Exit(1)
            },
        },
    },
)
```

锁通常来自：

```
k8s.io/client-go/tools/leaderelection/resourcelock
```

核心入口：

```
resourcelock.New(...)
```

锁资源一般使用 Kubernetes `Lease`。

------

## 20. `tools/record`

路径：

```
k8s.io/client-go/tools/record
```

作用：让 Controller 向 Kubernetes 写 Event。

入口：

```
broadcaster := record.NewBroadcaster()
```

创建 Recorder：

```
recorder := broadcaster.NewRecorder(
    scheme,
    corev1.EventSource{
        Component: "webapp-controller",
    },
)
```

记录事件：

```
recorder.Event(
    webApp,
    corev1.EventTypeNormal,
    "Created",
    "Created the Deployment",
)
```

然后可以看到：

```
kubectl describe webapp webapp-demo
kubectl get events
```

------

## 21. `util/retry`

路径：

```
k8s.io/client-go/util/retry
```

作用：处理 `resourceVersion` 冲突。

入口：

```
err := retry.RetryOnConflict(
    retry.DefaultRetry,
    func() error {
        deployment, err := client.Get(...)
        if err != nil {
            return err
        }

        copy := deployment.DeepCopy()
        copy.Spec.Replicas = ptr.To(int32(3))

        _, err = client.Update(
            ctx,
            copy,
            metav1.UpdateOptions{},
        )
        return err
    },
)
```

这个包非常适合传统的 `Get → 修改 → Update` 流程。

------

# 九、执行、日志和端口转发

## 22. `tools/remotecommand`

路径：

```
k8s.io/client-go/tools/remotecommand
```

作用：实现类似：

```
kubectl exec
```

入口：

```
executor, err :=
    remotecommand.NewSPDYExecutor(
        config,
        http.MethodPost,
        request.URL(),
    )
```

运行：

```
err = executor.StreamWithContext(
    ctx,
    remotecommand.StreamOptions{
        Stdin:  os.Stdin,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
        Tty:    true,
    },
)
```

------

## 23. `tools/portforward`

路径：

```
k8s.io/client-go/tools/portforward
```

作用：实现类似：

```
kubectl port-forward
```

入口：

```
forwarder, err :=
    portforward.New(...)
```

然后：

```
forwarder.ForwardPorts()
```

------

## 24. `tools/clientcmd/api`

路径：

```
k8s.io/client-go/tools/clientcmd/api
```

作用：表示 kubeconfig 的 Go 数据结构。

核心类型：

```
api.Config
api.Cluster
api.AuthInfo
api.Context
```

例如代码生成或修改 kubeconfig 时会使用：

```
config := api.NewConfig()
```

普通 Controller 只读取 kubeconfig，通常不直接用这个包。

------

# 十、认证插件

## 25. `plugin/pkg/client/auth`

路径：

```
k8s.io/client-go/plugin/pkg/client/auth
```

作用：注册额外认证插件。

通常通过空导入启用：

```
import (
    _ "k8s.io/client-go/plugin/pkg/client/auth"
)
```

或者只启用某个插件：

```
import (
    _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)
```

代码中不直接调用函数；空导入会执行包的注册逻辑。

------

# 十一、测试包

## 26. `kubernetes/fake`

路径：

```
k8s.io/client-go/kubernetes/fake
```

入口：

```
client :=
    fake.NewSimpleClientset(
        existingObjects...,
    )
```

它实现：

```
kubernetes.Interface
```

因此业务代码如果依赖接口：

```
type Controller struct {
    client kubernetes.Interface
}
```

测试时可以传 fake client。

不过 fake client 并不完全等价于真实 API Server，例如：

- 默认没有完整准入校验
- 部分 resourceVersion 行为不同
- 不会真的运行 Deployment controller
- 创建 Deployment 不会自动创建 Pod

------

## 27. `dynamic/fake`

路径：

```
k8s.io/client-go/dynamic/fake
```

入口：

```
client :=
    fake.NewSimpleDynamicClient(
        scheme,
        objects...,
    )
```

适合测试 Dynamic Client 和 CRD 代码。

------

# 十二、几个容易混淆但不属于 client-go 的包

下面经常与 client-go 一起使用，但属于其他 Go module。

## `k8s.io/api`

资源类型定义：

```
corev1.Pod
corev1.ConfigMap
appsv1.Deployment
batchv1.Job
```

## `k8s.io/apimachinery`

Kubernetes 通用基础类型：

```
metav1.ObjectMeta
metav1.GetOptions
schema.GroupVersionResource
unstructured.Unstructured
labels.Selector
runtime.Scheme
types.UID
wait.Until
```

以及 API 错误：

```
apierrors.IsNotFound(err)
apierrors.IsConflict(err)
```

## `k8s.io/code-generator`

为自定义 API 生成：

```
typed clientset
typed informer
typed lister
deepcopy
apply configurations
```

它是另一个模块，不是 client-go 的运行时包。

------

# 最应该先记住的入口

初学阶段，只需要先记住这十个：

| 需求              | 包                | 入口                                            |
| ----------------- | ----------------- | ----------------------------------------------- |
| 读取 kubeconfig   | `tools/clientcmd` | `BuildConfigFromFlags`                          |
| Pod 内连接        | `rest`            | `InClusterConfig`                               |
| 操作内置资源      | `kubernetes`      | `NewForConfig`                                  |
| 操作 CRD          | `dynamic`         | `NewForConfig`                                  |
| 查询 API 能力     | `discovery`       | `NewDiscoveryClientForConfig`                   |
| 内置资源 Informer | `informers`       | `NewSharedInformerFactory`                      |
| CRD Informer      | `dynamicinformer` | `NewDynamicSharedInformerFactory`               |
| Informer 基础能力 | `tools/cache`     | `ResourceEventHandlerFuncs`、`WaitForCacheSync` |
| 本地缓存读取      | `listers`         | 通常从 `Informer.Lister()` 取得                 |
| 工作队列          | `util/workqueue`  | `NewTypedRateLimitingQueue`                     |

把它们组合起来就是：

```
config, _ :=
    clientcmd.BuildConfigFromFlags("", kubeconfig)

clientset, _ :=
    kubernetes.NewForConfig(config)

factory :=
    informers.NewSharedInformerFactory(
        clientset,
        0,
    )

podInformer :=
    factory.Core().V1().Pods()

podInformer.Informer().
    AddEventHandler(...)

podLister :=
    podInformer.Lister()

queue :=
    workqueue.NewTypedRateLimitingQueue(
        workqueue.
            DefaultTypedControllerRateLimiter[string](),
    )

factory.Start(ctx.Done())

cache.WaitForCacheSync(
    ctx.Done(),
    podInformer.Informer().HasSynced,
)
```

这就是 client-go 最核心的包结构：



# 基础流程



`client-go` 是 Kubernetes 官方的 Go 客户端库。它提供两类能力：一类是向 API Server 发请求，另一类是支撑长期运行的 Controller，例如 Informer、缓存和工作队列。它不会替你决定“看到什么资源后该创建什么”；这部分规则要由你的代码编写。[client-go 官方说明](https://github.com/kubernetes/client-go)

## 1. 先建立连接：`rest.Config`

几乎所有 client-go 程序都从 `*rest.Config` 开始。它保存 API Server 地址、认证信息、TLS 设置等连接参数。

程序在集群**外**运行，通常读取 kubeconfig：

```
config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
if err != nil {
    return err
}
```

程序作为 Pod 在集群**内**运行，通常使用 Pod 的 ServiceAccount：

```
config, err := rest.InClusterConfig()
```

有了 `config`，才能创建下面的各种客户端。`clientcmd` 负责加载配置，`rest.Config` 是加载后的连接配置；两者还没有读取任何 Pod。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

## 2. Typed Client：操作已知类型

最常用的是 `kubernetes.Clientset`：

```
clientset, err := kubernetes.NewForConfig(config)
```

之后按 **API 组 → 资源 → 命名空间 → 操作** 调用：

```
pod, err := clientset.CoreV1().
    Pods("default").
    Get(ctx, "nginx-pod", metav1.GetOptions{})

deployment, err := clientset.AppsV1().
    Deployments("default").
    Get(ctx, "nginx", metav1.GetOptions{})
```

它叫 *typed*，因为返回值分别是 `*corev1.Pod`、`*appsv1.Deployment` 等明确的 Go 类型。你可以写 `deployment.Spec.Replicas`，编译器会检查字段名和类型。

常见方法如下：

| 方法     | 含义                                   |
| -------- | -------------------------------------- |
| `Get`    | 按名字读取一个资源                     |
| `List`   | 读取一批资源，可用 label selector 过滤 |
| `Create` | 创建资源                               |
| `Update` | 提交修改后的完整资源                   |
| `Patch`  | 修改指定字段                           |
| `Delete` | 删除资源                               |
| `Watch`  | 持续接收资源变化                       |

这些调用对应 Kubernetes API 请求。`Create` 成功意味着 API Server 接受了资源，**不意味着相关 Pod 已经 Ready**；后续还有控制器调谐、调度和容器启动。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

### 修改资源时为什么常常先 `Get`？

`Update` 通常采用“读出对象 → 修改副本 → 提交”的方式：

```
deployment, err := clientset.AppsV1().
    Deployments("default").
    Get(ctx, "nginx", metav1.GetOptions{})
if err != nil {
    return err
}

copy := deployment.DeepCopy()
replicas := int32(3)
copy.Spec.Replicas = &replicas

_, err = clientset.AppsV1().
    Deployments("default").
    Update(ctx, copy, metav1.UpdateOptions{})
```

对象带有 `resourceVersion`。如果在你读取后，别人先修改了它，你的 `Update` 可能得到 **Conflict**。正确做法是重新读取、重新计算修改，而不是盲目覆盖。读取错误也应按类型处理，例如 `apierrors.IsNotFound(err)`。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

## 3. Dynamic Client：操作运行时才知道的类型

内置的 `Clientset` 有 `CoreV1().Pods()`，却不会自动生成 `MyAppV1().WebApps()`。操作任意 CRD 时可以用 Dynamic Client：

```
dynamicClient, err := dynamic.NewForConfig(config)

gvr := schema.GroupVersionResource{
    Group:    "demo.example.com",
    Version:  "v1",
    Resource: "webapps", // 资源的复数名
}

webApp, err := dynamicClient.
    Resource(gvr).
    Namespace("default").
    Get(ctx, "demo", metav1.GetOptions{})
```

这里的关键是 **GVR**：Group、Version、Resource。返回值是 `*unstructured.Unstructured`，底层字段类似一个嵌套的 `map[string]interface{}`：

```
image, found, err := unstructured.NestedString(
    webApp.Object, "spec", "image",
)
```

Dynamic Client 灵活，适合 CRD 和通用工具；代价是没有 typed client 那样的编译期字段检查。若一个 CRD 项目需要大量强类型代码，也可以为它生成专用 clientset、lister 和 informer。[Dynamic Client 文档](https://pkg.go.dev/k8s.io/client-go/dynamic)、[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

## 4. 一次性查询与长期监听是两种写法

如果程序只需打印一次 Pod 列表，直接 `List` 即可：

```
pods, err := clientset.CoreV1().
    Pods("default").
    List(ctx, metav1.ListOptions{})
```

如果程序要长期响应变化，可以直接 `Watch`，但你得自己处理断线、重新列举等细节。写 Controller 时通常使用 **Informer**：它先 `List` 取得初始数据，再 `Watch` 后续变化，并维护一份本地缓存。其内部大致是：

```
API Server
  → Reflector（List / Watch）
  → DeltaFIFO
  → Informer
  → 本地缓存 + 事件处理函数
```

`Lister` 从这份**本地缓存**读数据，不是每次都向 API Server 发 GET。缓存可能略有延迟，所以 Controller 要容忍重复事件、暂时读不到刚创建的对象等情况。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

## 5. WorkQueue：把“发现变化”和“处理变化”分开

典型 Controller 的事件处理函数只做一件事：把资源的 `namespace/name` 放进队列。

```
AddFunc: func(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err == nil {
        queue.Add(key)
    }
}
```

Worker 再取出 key：

```
key, shutdown := queue.Get()
if shutdown {
    return
}
defer queue.Done(key)

if err := reconcile(key); err != nil {
    queue.AddRateLimited(key) // 失败后延迟重试
} else {
    queue.Forget(key)         // 成功，清除重试记录
}
```

这里最容易漏的是：**`Done` 表示本次处理结束；`Forget` 表示不要再保留这个 key 的失败重试状态。** WorkQueue 让事件处理不被慢操作阻塞，也让失败重试有固定位置。[client-go 架构文档](https://github.com/kubernetes/client-go/blob/master/ARCHITECTURE.md)

## 6. Reconcile：client-go 不替你写的部分

Informer 和队列解决“什么时候处理”。`reconcile(key)` 解决“应该做什么”：

```
从 Lister 读取期望状态
→ 从 Lister 读取实际状态
→ 比较
→ 如有差异，用 Client 创建或更新资源
```

例如期望副本数是 3、实际是 2，就更新 Deployment；已经是 3 就直接返回。**Reconcile 应允许同一个 key 被重复处理**，因为事件可能重复，失败也会重试。Kubernetes Controller 的目标是持续让实际状态接近期望状态。[Kubernetes Controller 文档](https://kubernetes.io/docs/concepts/architecture/controller/)

## 7. 这些包怎样分工

| 包                                 | 你通常用它做什么                |
| ---------------------------------- | ------------------------------- |
| `k8s.io/client-go/tools/clientcmd` | 读取 kubeconfig                 |
| `k8s.io/client-go/rest`            | 连接配置及较底层的 REST 能力    |
| `k8s.io/client-go/kubernetes`      | 内置资源的 typed clientset      |
| `k8s.io/client-go/dynamic`         | 通过 GVR 操作 CRD 等任意资源    |
| `k8s.io/client-go/informers`       | 内置资源的共享 Informer         |
| `k8s.io/client-go/tools/cache`     | Informer、Indexer、缓存相关工具 |
| `k8s.io/client-go/listers`         | 从 Informer 缓存读取内置资源    |
| `k8s.io/client-go/util/workqueue`  | Controller 的待处理队列         |

旁边常一起出现的 `k8s.io/api` 放 Pod、Deployment 等类型定义；`k8s.io/apimachinery` 放 `metav1`、GVR、`Unstructured`、错误判断等通用类型。它们与 client-go 配合使用，依赖版本通常应保持同一 Kubernetes 次版本线。[client-go 项目说明](https://github.com/kubernetes/client-go)、[版本对应说明](https://github.com/kubernetes/client-go#versioning)

**最实用的记法**：一次性工具通常是 `Config → Client → Get/List/Create`；长期运行的 Controller 则是 `Config → Client → Informer → Cache/Lister → WorkQueue → Reconcile → Client 写回`。