### 1. Goroutine

#### Goroutine基础

Goroutine是Go并发模型的核心。它是由Go运行时管理的轻量级线程。

```
package main

import (
    "fmt"
    "time"
)

func say(s string) {
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(s)
    }
}

func main() {
    go say("world") // 启动一个新的goroutine
    say("hello")    // 主goroutine
}
```

#### Goroutine的特点

- 轻量级：相比于系统线程，goroutine占用的内存非常少，通常只有几KB。
- 并发调度：由Go运行时调度，无需程序员显式创建和管理线程。

### 2. Channel

#### Channel基础

Channel是Go语言提供的用于在多个goroutine之间进行通信的管道。它们通过类型安全的方式来传递数据。

```
package main

import "fmt"

func main() {
    messages := make(chan string) // 创建一个channel

    go func() {
        messages <- "ping" // 发送数据到channel
    }()

    msg := <-messages // 从channel接收数据
    fmt.Println(msg)  // 输出：ping
}
```

#### 带缓冲的Channel

Channel可以是带缓冲的，这样发送方在缓冲区满之前不会被阻塞。

```
package main

import "fmt"

func main() {
    messages := make(chan string, 2) // 创建一个带缓冲的channel

    messages <- "buffered"
    messages <- "channel"

    fmt.Println(<-messages) // 输出：buffered
    fmt.Println(<-messages) // 输出：channel
}
```

#### 关闭Channel

关闭channel表示不会再向其发送数据，这样接收方可以知道什么时候停止接收。

```
package main

import "fmt"

func main() {
    jobs := make(chan int, 5)
    done := make(chan bool)

    go func() {
        for {
            j, more := <-jobs
            if more {
                fmt.Println("received job", j)
            } else {
                fmt.Println("received all jobs")
                done <- true
                return
            }
        }
    }()

    for j := 1; j <= 3; j++ {
        jobs <- j
        fmt.Println("sent job", j)
    }
    close(jobs)
    fmt.Println("sent all jobs")

    <-done
}
```

### 3. Select

#### Select基础

`select`语句使得一个goroutine可以等待多个通信操作。

```
package main

import (
    "fmt"
    "time"
)

func main() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }
}
```

#### 超时处理

使用`select`和`time.After`可以实现超时机制。

```
package main

import (
    "fmt"
    "time"
)

func main() {
    c := make(chan string, 1)

    go func() {
        time.Sleep(2 * time.Second)
        c <- "result"
    }()

    select {
    case res := <-c:
        fmt.Println(res)
    case <-time.After(1 * time.Second):
        fmt.Println("timeout")
    }
}
```

#### 非阻塞通信

`select`的默认分支用于实现非阻塞通信。

```
package main

import "fmt"

func main() {
    messages := make(chan string)
    signals := make(chan bool)

    select {
    case msg := <-messages:
        fmt.Println("received message", msg)
    default:
        fmt.Println("no message received")
    }

    msg := "hi"
    select {
    case messages <- msg:
        fmt.Println("sent message", msg)
    default:
        fmt.Println("no message sent")
    }

    select {
    case msg := <-messages:
        fmt.Println("received message", msg)
    case sig := <-signals:
        fmt.Println("received signal", sig)
    default:
        fmt.Println("no activity")
    }
}
```

### 4. WaitGroup

#### WaitGroup基础

`sync.WaitGroup`用于等待一组goroutine完成。

```
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
}
```

### 5. Mutex

#### Mutex基础

`sync.Mutex`用于在多个goroutine间实现互斥锁。

```
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    v   map[string]int
    mux sync.Mutex
}

func (c *SafeCounter) Inc(key string) {
    c.mux.Lock()
    c.v[key]++
    c.mux.Unlock()
}

func (c *SafeCounter) Value(key string) int {
    c.mux.Lock()
    defer c.mux.Unlock()
    return c.v[key]
}

func main() {
    c := SafeCounter{v: make(map[string]int)}
    for i := 0; i < 1000; i++ {
        go c.Inc("somekey")
    }

    fmt.Println(c.Value("somekey"))
}
```

### 6. Once

#### Once基础

`sync.Once`用于确保某个操作只执行一次。

```
package main

import (
    "fmt"
    "sync"
)

var once sync.Once

func initialize() {
    fmt.Println("Initialized")
}

func main() {
    for i := 0; i < 10; i++ {
        go once.Do(initialize)
    }
}
```

### 7. 条件变量

#### 条件变量基础

`sync.Cond`用于在某些条件下进行等待或唤醒。

```
go复制代码package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    mutex := sync.Mutex{}
    cond := sync.NewCond(&mutex)
    ready := false

    go func() {
        time.Sleep(1 * time.Second)
        mutex.Lock()
        ready = true
        cond.Broadcast() // 通知所有等待的goroutine
        mutex.Unlock()
    }()

    cond.L.Lock()
    for !ready {
        cond.Wait()
    }
    fmt.Println("Goroutine is ready")
    cond.L.Unlock()
}
```

### 8原子操作的概念

原子操作保证了在多线程或多goroutine环境下对共享资源的安全访问，它具有以下特点：

- **不可分割性**：原子操作在执行过程中不能被中断。
- **原子性**：要么执行完整个操作，要么不执行，不会出现部分执行的情况。
- **不可见性**：原子操作对其他线程或goroutine是不可见的，直到操作完成。

在并发编程中，原子操作通常涉及对共享变量的读取、修改和写入。

#### 原子函数和原子类型

Go语言提供了`sync/atomic`包，其中定义了一系列原子操作函数，用于对基本类型进行原子操作，如`Add`, `CompareAndSwap`, `Load`, `Store`等。

```
package main

import (
    "fmt"
    "sync/atomic"
)

func main() {
    var count int32 = 0

    // 原子的增加操作
    atomic.AddInt32(&count, 1)

    // 原子的比较并交换操作
    atomic.CompareAndSwapInt32(&count, 1, 2)

    // 原子的加载操作
    val := atomic.LoadInt32(&count)

    // 原子的存储操作
    atomic.StoreInt32(&count, 3)

    fmt.Println("Count:", count)  // 输出：Count: 3
}
```

#### 原子操作的应用场景

1. **计数器和标志位**：用于多goroutine共享的计数器或标志位，如线程池中的工作计数器。
2. **并发队列**：保证对队列的入队和出队操作的原子性。
3. **单例模式**：保证多个goroutine访问单例对象的安全性。

#### 原子操作的注意事项

1. **性能开销**：原子操作通常会比普通的非原子操作要慢，因为它们需要使用底层的处理器指令或锁机制来确保操作的原子性。
2. **适用性**：只有在确实需要对共享资源进行并发访问控制时才使用原子操作，避免过度使用导致性能下降。
3. **数据竞争**：虽然原子操作可以避免竞态条件，但并不是万能的解决方案，需要结合锁和其他并发控制机制使用。

### 高级用法

#### 使用goroutine池

为了控制并发数量，可以使用goroutine池。

```
package main

import (
    "fmt"
    "sync"
    "time"
)

type Pool struct {
    work chan func()
    wg   sync.WaitGroup
}

func NewPool(maxGoroutines int) *Pool {
    p := Pool{
        work: make(chan func()),
    }

    p.wg.Add(maxGoroutines)
    for i := 0; i < maxGoroutines; i++ {
        go func() {
            for work := range p.work {
                work()
            }
            p.wg.Done()
        }()
    }

    return &p
}

func (p *Pool) Run(work func()) {
    p.work <- work
}

func (p *Pool) Shutdown() {
    close(p.work)
    p.wg.Wait()
}

func main() {
    p := NewPool(3)

    for i := 0; i < 10; i++ {
        i := i
        p.Run(func() {
            fmt.Printf("Processing task %d\n", i)
            time.Sleep(time.Second)
        })
    }

    p.Shutdown()
}
```

#### 使用Context进行并发控制

`context`包用于控制多个goroutine的生命周期，如取消、超时等。

```
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        time.Sleep(2 * time.Second)
        cancel()
    }()

    select {
    case <-ctx.Done():
        fmt.Println("Context canceled")
    }
}
```

#### 使用信号量控制并发

信号量可以用来限制资源的并发访问。

```
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, sem *sync.WaitGroup) {
    defer sem.Done()
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var sem sync.WaitGroup
    sem.Add(3)

    for i := 1; i <= 3; i++ {
        go worker(i, &sem)
    }

    sem.Wait()
}
```