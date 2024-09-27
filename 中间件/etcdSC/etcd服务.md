# RAFT算法

etcd-main\server\etcdserver\raft.go

### 详细过程分析

#### 1. **Leader 选举过程**

- 集群初始化时，所有节点都是 `Follower`。
- 如果某个 `Follower` 超过一定时间（称为选举超时）没有收到 `Leader` 的心跳，它就会转变为 `Candidate`，进入选举状态。
- `Candidate` 向其他节点发送 `RequestVote` 请求，并增加自己的任期编号。
- 其他节点响应该请求，如果该 `Candidate` 的 `term` 比自己的 `term` 大，并且该节点还没有投票给其他 `Candidate`，就会将选票投给这个 `Candidate`。
- 如果某个 `Candidate` 获得超过半数节点的选票，它就会成为 `Leader`，lead就是一个原子uint64，`atomic.LoadUint64(&s.lead)`，并开始向其他节点发送 `AppendEntries`（心跳）信号来保持领导地位。
- 如果没有 `Candidate` 在超时时间内获得足够的选票，新的选举将开始，所有节点会递增 `term` 并重新投票。

#### 2. **日志复制过程**

- `Leader` 收到客户端请求后，将请求转换成日志条目，并将该条目追加到自己的日志中。
- `Leader` 通过 `AppendEntries` RPC 将日志条目复制到 `Follower`。
- 当超过半数的节点接收并保存了该条目，`Leader` 就提交该条目，并通知 `Follower` 也提交该条目。
- 提交后的日志条目会被应用到状态机，保证系统的最终一致性。

#### 3. **日志修复**

- 如果某个 `Follower` 的日志与 `Leader` 的日志不一致（如日志条目丢失或有不同的条目），`Leader` 会通过日志回溯的方式修复该 `Follower`。
- `Leader` 会根据 `Follower` 返回的日志匹配信息，找到最后一个共同的日志条目位置，然后重新从这个位置开始发送后续的日志条目，直到 `Follower` 的日志与 `Leader` 一致。

## 源码分析

raftnode结构如下：

```
type raftNode struct {
    lg *zap.Logger//指向日志记录器的指针，用于记录节点的运行状态和事件。zap 是一个高性能的日志库，适用于分布式系统的日志记录。

    tickMu *sync.RWMutex//保护 latestTickTs 字段的并发访问。它允许多个 goroutine 读取时间戳而不会引发冲突，同时在写入时阻止其他读写操作。

    latestTickTs time.Time//存储最新的 tick 时间戳。tick 是 Raft 的心跳机制的一部分，通常用于检测节点是否处于活跃状态。这个时间戳可以帮助判断节点是否需要进行选举。
    raftNodeConfig//包含节点配置的结构体

    msgSnapC chan raftpb.Message//通道，用于发送和接收快照消息。这对于状态机的恢复和快照同步非常重要，确保节点可以处理来自其他节点的快照请求。

    applyc chan toApply//发送需要应用到状态机的日志条目。toApply 结构体通常包含要应用的日志条目的信息，帮助将日志条目转换为状态机的实际操作。

    readStateC chan raft.ReadState//用于发送读取状态请求的结果。

    ticker *time.Ticker//定时器，用于定期触发事件（如心跳或选举）。这个定时器通常用来控制节点的心跳发送频率，以保持与其他节点的连接。

    td *contention.TimeoutDetector//指向竞争超时检测器的指针。这个检测器用于监测 Raft 心跳消息的发送情况，以防止心跳消息过于频繁导致的资源竞争。

    stopped chan struct{}//指示节点是否已停止
    done    chan struct{}//，指示节点的运行已经完成。
}
```

具体的raftNodeConfig节点如下：

```
type raftNodeConfig struct {
    lg *zap.Logger//指向日志记录器的指针，用于记录节点的运行状态和重要事件。

    isIDRemoved func(id uint64) bool//函数类型的字段，用于检查某个节点的 ID 是否已被从集群中移除。通过提供这个函数，Raft 节点可以动态地确认集群状态，确保在处理消息时只与活跃的节点交互。
    raft.Node//嵌入的 raft.Node 结构体，提供了 Raft 节点的基本功能和状态。
    raftStorage *raft.MemoryStorage//指向内存存储的指针，用于存储 Raft 的日志和状态信息。MemoryStorage 提供了高效的内存持久化，适合于快速访问和临时存储，但在节点重启后不会保留数据。
    storage     serverstorage.Storage//持久化存储的接口，通常实现为磁盘存储，确保日志和状态信息在节点重启后能够持久保留。
    heartbeat   time.Duration // 定义了心跳的时间间隔，通常用于日志记录。

    transport rafthttp.Transporter//用于发送和接收消息的传输接口。它负责节点之间的消息传递，要求发送消息时不能阻塞。如果没有传输实现，节点会导致 panic。
}
```

## 新建raft节点

```
func newRaftNode(cfg raftNodeConfig) *raftNode {//传入raftnodeconfig值，返回结点地址
lg = NewRaftLoggerZap(cfg.lg)//创建日志
raft.SetLogger(lg)//设置日志
r := &raftNode{//填充字段
//返回r
```