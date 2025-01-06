# 1. 数组

数组是一个固定长度的、同类型元素的序列。

#### 声明和初始化

```
// 声明一个长度为3的整数数组
var arr [3]int

// 声明并初始化
arr := [3]int{1, 2, 3}

// 使用...让编译器推断数组长度
arr := [...]int{1, 2, 3, 4, 5}
```

#### 访问和修改数组元素

```
arr[0] = 10      // 修改第一个元素
fmt.Println(arr) // 输出：[10 2 3]
```

#### 遍历数组

```
for i := 0; i < len(arr); i++ {
    fmt.Println(arr[i])
}

for index, value := range arr {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}
```

# 2. 切片

切片是动态数组，是对数组的一个视图，提供了更灵活、强大的接口。

#### 声明和初始化

```
// 声明一个空切片
var s []int

// 通过数组生成切片
arr := [5]int{1, 2, 3, 4, 5}
s = arr[1:4] // 包含索引1到3的元素

// 使用内置函数 make 创建切片
s = make([]int, 5)       // 创建一个长度和容量均为5的切片
s = make([]int, 3, 5)    // 创建一个长度为3，容量为5的切片

// 直接初始化切片
s = []int{1, 2, 3}
```

#### 访问和修改切片元素

```
s[0] = 10       // 修改第一个元素
fmt.Println(s)  // 输出：[10 2 3]
```

#### 追加元素

```
s = append(s, 4, 5)
fmt.Println(s)  // 输出：[10 2 3 4 5]
```

#### 遍历切片

```
for i := 0; i < len(s); i++ {
    fmt.Println(s[i])
}

for index, value := range s {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}
```

# 3. 映射（Map）

映射是键值对的集合，提供一种键值映射的无序数据结构。

#### 声明和初始化

```
// 使用 make 创建映射
var m map[string]int
m = make(map[string]int)

// 声明并初始化
m = map[string]int{"foo": 1, "bar": 2}
```

#### 访问和修改映射

```
m["foo"] = 10            // 修改键为 "foo" 的值
fmt.Println(m["foo"])    // 输出：10

// 检查键是否存在
value, ok := m["baz"]
if ok {
    fmt.Println("Key exists:", value)
} else {
    fmt.Println("Key does not exist")
}
```

#### 删除键值对

```
delete(m, "foo")
```

#### 遍历映射

```
for key, value := range m {
    fmt.Printf("Key: %s, Value: %d\n", key, value)
}
```

## map堆栈

在 Go 中，**`map` 通常会被分配到堆上**，原因如下：

### 1. **`map` 是一个引用类型**

- `map` 是 Go 的引用类型，和 `slice`、`channel` 类似。它的底层数据结构是通过指针间接访问的。
- 当你定义一个 `map` 变量时，实际上分配的是一个指向底层哈希表的指针，而哈希表的内存通常分配在堆上。

### 2. **逃逸分析的影响**

- 编译器通过逃逸分析决定变量是分配在栈上还是堆上。如果一个变量的生命周期超出了函数作用域（例如，返回值或被其他 goroutine 使用），它就会分配到堆上。
- `map` 的底层数据是动态增长的，而动态内存分配通常需要分配在堆上，即使 `map` 的变量本身可能存储在栈上。

### 3. **例子说明**

```
go复制代码package main

func main() {
    m := make(map[string]int)
    m["key"] = 42
}
```

- 在上面的例子中，`m` 是一个指向底层哈希表的指针，`m` 本身可能分配在栈上，但底层哈希表的内存分配在堆上。
- 编译器会通过逃逸分析来判断 `m` 的具体分配位置。如果 `m` 被函数返回或跨 goroutine 使用，则 `m` 本身也会分配到堆上。

### 4. **如何查看分配位置**

你可以使用 `go build` 的 `-gcflags="-m"` 参数查看逃逸分析的结果。例如：

```
bash


复制代码
go build -gcflags="-m" main.go
```

输出类似于：

```
css


复制代码
main.go:4:6: moved to heap: m
```

这表明 `m` 的底层数据被分配到了堆上。

### 5. **栈 vs 堆的简单规则**

- **栈**：分配小、生命周期明确且局限在当前作用域的变量。
- **堆**：用于动态分配、跨作用域或生命周期不确定的对象。



## map并发

### 使用 `map` 加锁和直接使用 `sync.Map` 的区别：

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





# 4. 结构体（Struct）

结构体是将多个不同类型的字段聚合在一起的数据结构。

#### 声明和初始化

```
// 声明结构体类型
type Person struct {
    Name string
    Age  int
}

// 创建结构体实例
var p Person
p = Person{"Alice", 30}
p = Person{Name: "Bob", Age: 25}

// 使用 new 创建结构体指针
pPtr := new(Person)
pPtr.Name = "Charlie"
pPtr.Age = 35
```

#### 访问和修改结构体字段

```
p.Name = "David"
fmt.Println(p.Name) // 输出：David
```

#### 嵌套结构体

```
type Address struct {
    City, State string
}

type Person struct {
    Name    string
    Age     int
    Address Address
}

p := Person{
    Name: "Alice",
    Age:  30,
    Address: Address{
        City:  "New York",
        State: "NY",
    },
}
```

# 5. 字符串

字符串是不可变的字节序列，用于表示文本。

#### 声明和初始化

```
var str string
str = "Hello, World!"

str := "Hello, Go!"
```

#### 访问字符串

```
fmt.Println(str[0]) // 输出：72 (H 的 ASCII 码)
fmt.Println(str[:5]) // 输出：Hello
```

#### 字符串长度

```
length := len(str)
fmt.Println(length) // 输出：13
```

#### 字符串拼接

```
s1 := "Hello"
s2 := "World"
s := s1 + ", " + s2 + "!"
fmt.Println(s) // 输出：Hello, World!
```

#### 字符串遍历

```
for i := 0; i < len(str); i++ {
    fmt.Printf("%c ", str[i])
}
fmt.Println()

for index, runeValue := range str {
    fmt.Printf("Index: %d, Rune: %c\n", index, runeValue)
}
```

# 6. 通道（Channel）

通道是用于 goroutine 之间通信的管道。

#### 声明和初始化

```
// 创建一个无缓冲的通道
var ch chan int
ch = make(chan int)

// 创建一个有缓冲的通道
ch = make(chan int, 2)
```

#### 发送和接收数据

```
ch <- 10     // 发送数据
value := <-ch // 接收数据
fmt.Println(value) // 输出：10
```

#### 关闭通道

```
close(ch)
```

#### 遍历通道

```
// 在一个 goroutine 中发送数据
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}()

// 在主 goroutine 中接收数据
for value := range ch {
    fmt.Println(value)
}
```