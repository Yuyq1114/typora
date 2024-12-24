- 并发编程中，main方法结束就直接结束，time.sleep让main等待几秒。
- 不同目录下函数名大写表示权限为公共
- 单引号 byte、rune，双引号 字符串，反引号 展示输出
- 空切片和nil切片区别，nil切片是没有分配任何底层数组的切片。空切片是分配了底层数组但不包含任何元素的切片。
- `for i:=range{ }`是一种用于遍历切片、数组、字符串、映射（map）或通道（channel）的语法。
- `array:=make([]int,length)`这样，`var array [length]int`会报错
- 遍历输出字符串的方式：1. `printf（ %q ）`go语法安全转义 2.`string()` 类型转换
- range中的i是int,不是元素
- 链表中，变量声明：`var :=1 new()`用来声明指针变量 `make()` 用来初始化slice，map，channel。方法：`head.addNode(v)`  非方法：`addNode(head,v)`。
- nil只能判断函数，管道，切片，映射，接口，指针。
- 函数内部不可声明函数，但可以用闭包函数`:=` 
- &p 取p的地址，*p，p为地址时取p指向的值。`var p *int`：声明p存储int类型变量的地址。
- `for     ；（条件）；{ }` 满足条件退出，类似while
- 切片和数组区别：初始化时：`切片make，[]` `数组 [3]~ , [...]`传递：`切片指针，函数值`
- `return（true，false）`可直接用条件，如a==1
- 返回值名可在函数上定义，直接return
- 不换行，print（a ' '）
- 协程中，main方法结束程序就结束
- go range通道前先关闭
- defer是栈结构
- 函数可以多个返回值
- new传指针，make传值
- panic recover
- main init函数
- `func（）{}`定义一个匿名函数  `func（）{}（）`定义一个匿名函数并执行
- 通过`m:="string"  n=fmt.Sprintf("~",m)` 进行格式化
- `q=append(q,edges[x]...)`  将edges中的元素逐一添加，...表示将切片展开，追加到另一个上面。
- `var a func (int) int`   , `a:=func (b int) int {}`将函数作为参数传递
- chan输出前未关闭chan会报错，deadlock！因为循环输出时无值会死锁 close（ch）
- `wg.Wait()`有时候会在一个闭包内部，有时候在主携程
- `fmt.Scanf` 跳过输入的原因通常与 `fmt.Scanf` 对输入的解析方式有关。特别是在读取字符串时，`fmt.Scanf("%s", &name)` 和 `fmt.Scanf("%s", &message)` 使用 `%s` 格式化字符串，它会读取直到第一个空白字符（如空格、制表符或换行符）。因为它还会读取到输入缓冲区中的换行符
- go使用`github.com/confluentinc/confluent-kafka-go/v2/kafka`中的`kafka.NewProducer`是，一直显示未定义，这是由于go未开启`$env:CGO_ENABLED = "1"`，让go可以引用c的编译。且需要安装gcc编译器。
- 函数参数定义时使用 `...` 代表可变参数，也就是参数数量不定。可变参数允许传递任意数量的某种类型的参数给函数。可变参数必须是函数的最后一个参数。如果函数有其他参数，可变参数必须放在最后。可变参数实际上是一个切片（slice），在函数内部可以像操作切片那样使用它。
- 可以实现主线程等待子线程信号的几种方法“

```
1.无缓冲管道
func main() {
	wait := make(chan struct{})
	go func() {
		for {
			i := rand.Intn(10)
			fmt.Println(i)
			if i == 5 {
				close(wait)
				break
			}
		}
	}()
	<-wait
}
2.使用sync.WaitGroup
func main() {
	var wg sync.WaitGroup
	wg.Add(1) // 增加等待计数

	go func() {
		defer wg.Done() // 减少等待计数
		i := rand.Intn(10)
		fmt.Println(i)
	}()

	wg.Wait() // 主线程等待所有goroutine完成
	fmt.Println("Main goroutine continues.")
}
3.使用 sync.Mutex 和 sync.Cond
func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	done := false

	go func() {
		mu.Lock()
		i := rand.Intn(10)
		fmt.Println(i)
		done = true
		cond.Signal() // 通知主线程继续
		mu.Unlock()
	}()

	mu.Lock()
	for !done {
		cond.Wait() // 主线程等待信号
	}
	mu.Unlock()

	fmt.Println("Main goroutine continues.")
}
4.使用 context 包
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		i := rand.Intn(10)
		fmt.Println(i)
		cancel() // 通知主线程继续
	}()

	<-ctx.Done() // 主线程等待信号
	fmt.Println("Main goroutine continues.")
}
5.使用带缓冲的通道（Buffered Channel）
func main() {
	done := make(chan struct{}, 2) // 缓冲大小为2，允许两个goroutine同时完成

	go func() {
		i := rand.Intn(10)
		fmt.Println(i)
		done <- struct{}{} // 通知主线程继续
	}()

	go func() {
		i := rand.Intn(10)
		fmt.Println(i)
		done <- struct{}{} // 通知主线程继续
	}()

	<-done // 主线程等待第一个goroutine完成
	<-done // 主线程等待第二个goroutine完成
	fmt.Println("Main goroutine continues.")
}
6.使用 sync.Once
func main() {
	var once sync.Once
	done := make(chan struct{})

	go func() {
		once.Do(func() {
			i := rand.Intn(10)
			fmt.Println(i)
			close(done)
		})
	}()

	<-done
	fmt.Println("Main goroutine continues.")
}
```

- 对于 `JsonRes` 结构体中的 `Data` 字段，选择使用 `interface{}` 是为了让这个字段能够灵活地存储不同类型的返回数据。这种设计通常用于通用的 API 响应结构，因为 API 的返回值可能会根据具体的请求内容而不同。
- ### 使用 `map` 加锁和直接使用 `sync.Map` 的区别：

  #### 1. **设计目的**

  - **`map` 加锁**：手动在普通的 `map` 上通过锁机制（`sync.Mutex` 或 `sync.RWMutex`）实现并发安全。适用于你想要更精细地控制锁的行为，比如读取时可以使用读锁以提高效率。
  - **`sync.Map`**：专门为并发访问而设计，内置了优化的并发安全机制。特别适用于频繁读写的场景，且减少了手动加锁的复杂度。

  #### 2. **性能**

  - **`map` 加锁**：手动控制锁的开销取决于锁的粒度。使用 `sync.RWMutex` 时，读写性能可以分开优化，读多写少的场景性能可能比 `sync.Map` 更好。但锁的开销仍然存在。
  - **`sync.Map`**：在某些场景下（尤其是频繁读写的并发场景），`sync.Map` 内部做了多种优化，比如分片锁机制（类似于读写锁的机制）以及无锁的原子操作，因此对于频繁读写的小型对象集合，它的性能可能会优于 `map` 加锁。然而，`sync.Map` 在大量写操作时性能会受到影响，适合读多写少的场景。

  #### 3. **使用难度**

  - **`map` 加锁**：需要开发者显式管理锁的生命周期和范围，容易出现死锁、锁粒度控制不当等问题，需要更加小心地编写代码。
  - **`sync.Map`**：开发者不需要关心加锁、解锁的问题，代码更加简洁，降低了出错的可能性。

  #### 4. **特性支持**

  - **`map` 加锁**：普通的 `map` 支持所有 `map` 的标准操作（遍历、键值对操作等）。但在遍历时需要小心处理加锁问题，以避免出现竞争条件。
  - **`sync.Map`**：它的操作接口不同于普通的 `map`，不支持索引访问（`m[key]`），而是通过 `Load`、`Store`、`LoadOrStore` 和 `Delete` 等方法进行操作。此外，它的遍历不保证一致性，也就是说在遍历过程中，其他 goroutine 对 `sync.Map` 的修改可能会影响遍历结果。

  #### 5. **遍历**

  - **`map` 加锁**：在遍历一个普通 `map` 时，需要锁住整个 `map`，否则可能会出现并发写入的问题。使用写锁（`sync.Mutex`）时，所有的并发访问都会被阻塞。
  - **`sync.Map`**：遍历时不会锁定整个 `map`，但遍历时的快照不保证是最新的内容，因为 `sync.Map` 可能在遍历过程中被其他 goroutine 修改。

  ### 什么时候选择 `map` 加锁或 `sync.Map`？

  - **读多写少的并发场景**：`sync.Map` 会表现得更好，因为它的读操作经过高度优化，适合频繁读取和少量写入的场景。
  - **写操作频繁或需要更灵活控制锁**：使用带锁的普通 `map` 会更合适。你可以根据读写的比例选择不同的锁策略（如读写锁），从而在写操作较多时提升性能。
- go里面只有值传递，没有引用传递。比如传一个大型结构体时，实际上是吧指针值拷贝一份过去，而不是把整个
- 闭包特点
- | 第一段代码 | 20, 10 | `defer` 在声明时即捕获了参数的值。两个 `defer` 分别捕获了 `x` 的不同值（10 和 20），并按照逆序执行。 |
  | ---------- | ------ | ------------------------------------------------------------ |
  |            |        |                                                              |

  | 第二段代码 | 20, 20 | 闭包捕获的是变量的引用，`defer` 延迟执行闭包，此时 `x` 的值已经被修改为 `20`。 |
  | ---------- | ------ | ------------------------------------------------------------ |
  |            |        |                                                              |
  
- go run -gcflags="-m" main.go进行逃逸分析















