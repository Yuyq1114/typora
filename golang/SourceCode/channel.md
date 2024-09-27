# 基本结构

src\runtime\chan.go

```
type hchan struct {
    qcount   uint           // total data in the queue表示 channel 中当前存储的元素个数。
    dataqsiz uint           // size of the circular queue表示环形缓冲区的大小。
    buf      unsafe.Pointer // points to an array of dataqsiz elements指向缓冲区数组的指针。
    elemsize uint16// 表示每个元素的大小（以字节为单位）。
    closed   uint32// 表示 channel 是否已关闭。
    timer    *timer // timer feeding this chan指向与此 channel 相关联的定时器。
    elemtype *_type // element type指向 channel 中元素类型的元信息。
    sendx    uint   // send index表示环形缓冲区中的发送下标。
    recvx    uint   // receive index表示环形缓冲区中的接收下标。
    recvq    waitq  // list of recv waiters表示等待接收数据的 goroutine 队列。
    sendq    waitq  // list of send waiters表示等待发送数据的 goroutine 队列。

    // lock protects all fields in hchan, as well as several
    // fields in sudogs blocked on this channel.
    //
    // Do not change another G's status while holding this lock
    // (in particular, do not ready a G), as this can deadlock
    // with stack shrinking.
    lock mutex保护 hchan 结构体中的所有字段，确保并发操作下的安全性。
}
```

## 其中waitq结构体如下

```
type waitq struct {//waitq 的核心作用是作为等待队列，在同步操作（如 channel 操作或锁操作）时，当资源（如缓冲区或锁）不可用时，阻塞的 goroutine 被挂起并加入 waitq 中。一
    first *sudog//first 指向当前最早进入等待队列的 goroutine，通常是最先被唤醒的 goroutine。
    last  *sudog//last 指向当前最后一个进入等待队列的 goroutine，表示队列中最晚加入等待的 goroutine。当有新的 goroutine 进入等待时，它会被加到 last 的后面。
}
```

## 其中mutex结构如下

```
type mutex struct {
    // Empty struct if lock ranking is disabled, otherwise includes the lock rank
    lockRankStruct//lockRankStruct 是一个用于实现锁排序的结构体。它包含了锁的优先级或等级信息，帮助防止死锁。当锁排序（lock ranking）被启用时，lockRankStruct 可能会包含锁的排序信息（rank），否则它是一个空的结构体。锁排序通过确保锁获取的顺序来避免死锁。
    // Futex-based impl treats it as uint32 key,
    // while sema-based impl as M* waitm.
    // Used to be a union, but unions break precise GC.
    key uintptr//uintptr 类型的字段，表示锁的实现细节。
}
```

# 创建chan

```
func makechan(t *chantype, size int) *hchan {//接受chan类型指针和大小
```

chantype结构如下所示；

```
type ChanType struct {
    Type // 继承 Type 类型，表示它是一个 Go 类型
    Elem *Type// 表示 channel 内传输的数据类型
    Dir  ChanDir// 表示 channel 的方向，发送、接收或双向
}
```

```
var c *hchan //创建
c.elemsize = uint16(elem.Size_)//填充字段
c.elemtype = elem//if elem.Size_ >= 1<<16 {抛出错误，编译器会检查大于64kb，写一个安全点
c.dataqsiz = uint(size)
lockInit(&c.lock, lockRankHchan)//lockInit 是 Go 运行时中的一个函数，用于初始化一个锁对象。这个函数的作用是设置锁的初始状态，并配置必要的参数。具体的
func lockInit(l *mutex, rank int) {
    // 初始化锁的实现
}
//l *mutex：指向要初始化的 mutex 对象（或锁对象）。在这里，c.lock 是一个 mutex 类型的对象。
//rank int：锁的优先级或排序等级。这个参数用于避免死锁（通过防止锁的嵌套顺序混乱）。
```



# 关闭

```
func closechan(c *hchan) {//输入一个chan指针
```

```
lock(&c.lock)
//释放所有读者
sg := c.recvq.dequeue()//从 recvq 队列中取出一个等待接收数据的 sudog。sudog 是 Go 运行时用来存储等待状态的 goroutine 信息的结构体。
gp := sg.g//从 sudog 中获取对应的 goroutine 指针 gp。
glist.push(gp)//将 goroutine gp 推入 glist，这个列表存储了被唤醒的 goroutine。

//释放所有写者
sg := c.sendq.dequeue()//sendq队列里取出
。。。

unlock(&c.lock)//解锁

//处理唤醒 goroutine 的逻辑。它的目的是从 glist 中逐个取出 goroutine，并将它们标记为可执行的状态，加入调度队列。
for !glist.empty() {//这个循环会一直执行，直到 glist 为空。glist 是一个存储了 goroutine 列表的数据结构，包含所有需要被唤醒的 goroutine。
		gp := glist.pop()
		gp.schedlink = 0//将当前 goroutine 链接到调度器中的其他 goroutine。这里将其设置为 0，表示当前 goroutine 的调度链接已被清除，防止出现不必要的链表链接。
		goready(gp, 3)//将 goroutine gp 标记为可执行状态，并将它放入到调度队列中。
	}
```



# 接收

```
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
typedmemclr(c.elemtype, ep)//用于清除（清空）指定类型内存区域的函数。在接收数据之前，首先清空 ep 指向的内存空间，以避免潜在的数据污染。
recv(c, sg, ep, func() { unlock(&c.lock) }, 3)//从 channel 中取出数据并复制到 ep 指向的内存位置。
mysg.elem = ep// mysg 是当前 goroutine 的 sudog 结构体，mysg.elem 表示接收到的数据存放的位置。在接收数据后，将元素的指针 ep 存入 sudog 中，表示该 goroutine 已经接收到数据。
```



## 其中recv如下

```
func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
recvDirect(c.elemtype, sg, ep)//将数据直接从 channel 发送方的 sudog 中提取到接收方的内存地址 ep 中。
typedmemmove(c.elemtype, ep, qp)//一个类型安全的内存移动函数，用于将元素从 channel 的缓冲区 qp 复制到 ep 指向的内存位置。
unlockf()
goready(gp, skip+1)//让 gp 这个 goroutine 准备好执行的函数，唤醒这个等待接收数据的 goroutine 以继续执行。skip+1 是一个栈帧的偏移量

```

### 其中recvDirect如下

```
func recvDirect(t *_type, sg *sudog, dst unsafe.Pointer) {
    src := sg.elem//src 是发送方的数据指针。sg.elem 指向了 sudog 中的 elem 字段，这是发送方传递的数据的位置。
    typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.Size_)//在不同内存位置之间进行写屏障检查，特别是当涉及到垃圾收集器（GC）的时候。如果 channel 元素包含指针或者其他需要跟踪的内存，写屏障能够确保 GC 正确处理数据的复制。
    memmove(dst, src, t.Size_)//内存复制函数，负责将 src 中的数据移动到 dst，大小为 t.Size_。这里使用的是一个底层的内存操作
}
```

# 发送

```
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {//其中ep指数据元素，可能指向某个栈。callerpc uintptr: 调用者的程序计数器地址，通常用于跟踪调用位置或用于调试。
send(c, sg, ep, func() { unlock(&c.lock) }, 3)//主要作用是将元素 ep 发送到 channel c 中，并在完成后执行一些操作。具体参数说明：
typedmemmove(c.elemtype, qp, ep)//这是一个内存操作函数 typedmemmove，它负责将类型化的数据从源指针 ep 复制到目标指针 qp（qp 是 channel 的缓冲区位置）。这是一种低级的内存复制操作，通常用于将数据从栈或其他位置复制到 channel 缓冲区中，保证数据的正确传输。
mysg.elem = ep//mysg 是 sudog（当前 goroutine 在 channel 上的等待记录）。这里将 ep（待发送的数据元素指针）存储到 mysg.elem 中。这样做的目的是让 sudog 知道正在发送的数据是什么，以便在需要的时候（如恢复 goroutine 时）能够访问这个数据。
KeepAlive(ep)//防止 Go 垃圾回收器过早回收 ep 指向的数据的函数
```

## 其中send如下：

```
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
sendDirect(c.elemtype, sg, ep)//sendDirect 函数的主要作用是直接将数据传递给接收方的 sudog，而不是将数据放入 channel 的缓冲区中。
```

### 其中senddirect如下

```
func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
    dst := sg.elem//sg.elem 是接收方的内存位置（目标存储槽），将其赋值给 dst。
    typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.Size_)//typeBitsBulkBarrier 用于设置内存屏障。内存屏障是一种防止编译器或 CPU 对某些指令进行重排序的机制，在并发环境下保证内存操作的顺序一致。这里的 typeBitsBulkBarrier 会确保对 dst 和 src 的操作按照预期顺序进行，并且防止类型相关的内存管理错误。
    memmove(dst, src, t.Size_)//memmove 是一个内存复制函数，用于将 src（源地址）的数据复制到 dst（目标地址）。
}
```

# memo具体实现如下：

src\runtime\memmove_amd64.s

```
// func memmove(to, from unsafe.Pointer, n uintptr)
// ABIInternal for performance.
//用于将源地址到目标地址的数据移动操作。
TEXT runtime·memmove<ABIInternal>(SB), NOSPLIT, $0-24
    // AX = to
    // BX = from
    // CX = n
    MOVQ   AX, DI//将 AX 寄存器中的值（目标地址 to）移动到 DI 寄存器中。DI 寄存器是 x86-64 架构中的一个通用寄存器，用于存储目标地址。
    MOVQ   BX, SI//将 BX 寄存器中的值（源地址 from）移动到 SI 寄存器中。SI 寄存器是 x86-64 架构中的一个通用寄存器，用于存储源地址。
    MOVQ   CX, BX//将 CX 寄存器中的值（字节数 n）移动到 BX 寄存器中。这里的 BX 寄存器被重新用来存储字节数 n。
```