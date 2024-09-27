# 原子操作

src\sync\atomic\value.go

**参数检查**

```
if val == nil {
    panic("sync/atomic: store of nil value into Value")
}
```

如果传入的 `val` 是 `nil`，则抛出一个 panic，因为 `Value` 不允许存储 `nil` 值。

**类型转换**

```
vp := (*efaceWords)(unsafe.Pointer(v))
vlp := (*efaceWords)(unsafe.Pointer(&val))
```

`efaceWords` 是一个内部结构，用于表示 `interface{}` 的数据。`unsafe.Pointer` 用于将 `v` 和 `val` 转换为 `efaceWords` 结构，以便直接访问底层字段。

**循环处理**

```
for {
    typ := LoadPointer(&vp.typ)
    if typ == nil {
        // Attempt to start first store.
        // Disable preemption so that other goroutines can use
        // active spin wait to wait for completion.
        runtime_procPin()
        if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(&firstStoreInProgress)) {
            runtime_procUnpin()
            continue
        }
        // Complete first store.
        StorePointer(&vp.data, vlp.data)
        StorePointer(&vp.typ, vlp.typ)
        runtime_procUnpin()
        return
    }
    if typ == unsafe.Pointer(&firstStoreInProgress) {
        // First store in progress. Wait.
        // Since we disable preemption around the first store,
        // we can wait with active spinning.
        continue
    }
    // First store completed. Check type and overwrite data.
    if typ != vlp.typ {
        panic("sync/atomic: store of inconsistently typed value into Value")
    }
    StorePointer(&vp.data, vlp.data)
    return
}
```

- **获取类型**：

  ```
  typ := LoadPointer(&vp.typ)
  ```

  从 `vp` 中加载类型指针。如果 `typ` 是 `nil`，说明 `Value` 对象尚未被初始化或正在被其他 goroutine 初始化。

- **首次存储尝试**：

  ```
  runtime_procPin()
  if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(&firstStoreInProgress)) {
      runtime_procUnpin()
      continue
  }
  ```

  如果 `typ` 为 `nil`，尝试使用 `CompareAndSwapPointer` 将其设置为 `firstStoreInProgress`。在此过程中，调用 `runtime_procPin()` 禁用抢占，以确保当前 goroutine 可以活跃地等待完成。这是为了避免在并发环境下多个 goroutine 同时初始化 `Value` 对象。

- **完成存储**：

  ```
  StorePointer(&vp.data, vlp.data)
  StorePointer(&vp.typ, vlp.typ)
  runtime_procUnpin()
  ```

  设置 `vp.data` 和 `vp.typ`，完成首次存储。`runtime_procUnpin()` 重新启用抢占。

- **等待首次存储完成**：

  ```
  if typ == unsafe.Pointer(&firstStoreInProgress) {
      continue
  }
  ```

  如果 `typ` 为 `firstStoreInProgress`，则表示首次存储正在进行，当前 goroutine 需要等待。通过 `continue` 重新进入循环进行等待。

- **类型检查和数据存储**：

  ```
  if typ != vlp.typ {
      panic("sync/atomic: store of inconsistently typed value into Value")
  }
  StorePointer(&vp.data, vlp.data)
  return
  ```

  如果首次存储已完成，检查类型是否一致。如果类型不一致，则抛出 panic。否则，将数据存储到 `vp.data` 中，并退出函数。

# WaitGroup

src\sync\waitgroup.go

```
type WaitGroup struct {
    noCopy noCopy

    state atomic.Uint64 // high 32 bits are counter, low 32 bits are waiter count.
    sema  uint32
}
```

其中atomic.Uint64如下：

```
type Uint64 struct {
    _ noCopy//_ noCopy：这个字段的类型是 noCopy，通常它是一个用于禁止结构体实例通过复制方式使用的技巧。noCopy 可能是一个空的结构体，配合 go vet 工具使用，以帮助检测误用，比如检测到通过复制传递该结构体时发出警告。这样可以防止开发者错误地复制这个类型的值。
    _ align64//_ align64：这个字段用于对齐。Go 语言中的 align64 是一种用于内存对齐的技巧（通常定义为一个具有特定对齐属性的结构体）。它确保 v uint64 字段在内存中是 64 位对齐的，这对于某些架构下的性能优化，或处理并发时保证变量的对齐性非常有用。尤其是在并发访问的时候，未对齐的访问可能会导致性能下降甚至引发不一致性问题。
    v uint64
}
```

## add方法

```
state := wg.state.Add(uint64(delta) << 32)  //高 32 位：表示当前 WaitGroup 中的计数器，也就是 wg.counter，它用于跟踪有多少个 goroutine 在执行中。低 32 位：表示 WaitGroup 是否在等待 goroutine 结束，也就是 wg.waiters。将 delta 左移 32 位，表示将其值放入 64 位整数的高 32 位。
v := int32(state >> 32)  //将 state 的高 32 位提取出来。存储的是 WaitGroup 中的计数器值。
w := uint32(state)//将 state 的低 32 位提取出来，并转换为无符号 32 位整数。这部分用于跟踪 WaitGroup 中是否有 goroutine 在等待（通过 Wait() 函数等待其他 goroutine 完成）。
```

## Done方法

```
func (wg *WaitGroup) Done() {
    wg.Add(-1)//add-1
}
```

## Wait方法

```
if v == 0 {//不需要等
if wg.state.CompareAndSwap(state, state+1) {//CompareAndSwap 是一个原子操作，用于检查 wg.state 是否等于当前的 state 值。如果相等，它将 wg.state 更新为 state + 1。
			if race.Enabled && w == 0 {
				// Wait must be synchronized with the first Add.
				// Need to model this is as a write to race with the read in Add.
				// As a consequence, can do the write only for the first waiter,
				// otherwise concurrent Waits will race with each other.
				race.Write(unsafe.Pointer(&wg.sema))//向 race detector 提供信息，它告诉 race detector，在 Wait() 过程中对 wg.sema 做了写操作。wg.sema 是用于信号量机制的变量，控制等待的 goroutine 数量。
			}
			runtime_Semacquire(&wg.sema)//runtime_Semacquire 是 Go 运行时中的一个底层函数，用于阻塞当前 goroutine，直到信号量 wg.sema 被释放（也就是其他 goroutine 通过 Add() 或 Done() 操作减少 WaitGroup 计数器时）。
			if wg.state.Load() != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))//告诉 race detector，当前 goroutine 拥有对 WaitGroup 的控制权。
			}
			return
		}
```