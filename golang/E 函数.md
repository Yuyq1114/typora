## 1. 基本函数

#### 函数声明与定义

```
func functionName(parameterList) returnType {
    // 函数体
}
```

#### 示例

```
func add(a int, b int) int {
    return a + b
}

func main() {
    result := add(3, 4)
    fmt.Println(result) // 输出：7
}
```

#### 多个返回值

Go语言函数可以返回多个值。

```
func divide(a int, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

func main() {
    q, r := divide(10, 3)
    fmt.Printf("Quotient: %d, Remainder: %d\n", q, r) // 输出：Quotient: 3, Remainder: 1
}
```

#### 命名返回值

可以为返回值命名，使代码更清晰。

```
func divide(a int, b int) (quotient int, remainder int) {
    quotient = a / b
    remainder = a % b
    return
}

func main() {
    q, r := divide(10, 3)
    fmt.Printf("Quotient: %d, Remainder: %d\n", q, r)
}
```

## 2. 参数传递

#### 值传递

Go语言中的所有参数传递都是值传递，即传递的是参数的副本。

```
func modify(x int) {
    x = 10
}

func main() {
    a := 5
    modify(a)
    fmt.Println(a) // 输出：5，a 未改变
}
```

#### 引用传递

通过传递指针可以实现引用传递。

```
func modify(x *int) {
    *x = 10
}

func main() {
    a := 5
    modify(&a)
    fmt.Println(a) // 输出：10，a 被修改
}
```

## 3. 变长参数

可以使用 `...` 语法实现变长参数。

```
func sum(numbers ...int) int {
    total := 0
    for _, number := range numbers {
        total += number
    }
    return total
}

func main() {
    result := sum(1, 2, 3, 4, 5)
    fmt.Println(result) // 输出：15
}
```

## 4. 匿名函数和闭包

#### 匿名函数

匿名函数没有名字，可以在声明时直接调用。

```
func main() {
    func(msg string) {
        fmt.Println(msg)
    }("Hello, World!") // 输出：Hello, World!
}
```

#### 闭包

闭包是一个函数，它引用了其外部作用域中的变量。

```
func main() {
    a := 0
    increment := func() int {
        a++
        return a
    }
    fmt.Println(increment()) // 输出：1
    fmt.Println(increment()) // 输出：2
}
```

## 5. 函数类型和高阶函数

#### 函数类型

函数可以作为变量类型进行声明和赋值。

```
func add(a, b int) int {
    return a + b
}

func main() {
    var op func(int, int) int
    op = add
    fmt.Println(op(1, 2)) // 输出：3
}
```

#### 高阶函数

高阶函数可以接受函数作为参数或返回一个函数。

```
func apply(op func(int, int) int, a, b int) int {
    return op(a, b)
}

func multiply(a, b int) int {
    return a * b
}

func main() {
    result := apply(multiply, 3, 4)
    fmt.Println(result) // 输出：12
}
```

## 6. 方法

方法是附加在某个类型上的函数。

#### 定义方法

```
type Person struct {
    Name string
}

func (p Person) Greet() {
    fmt.Printf("Hello, my name is %s\n", p.Name)
}

func main() {
    p := Person{Name: "Alice"}
    p.Greet() // 输出：Hello, my name is Alice
}
```

#### 指针接收者方法

使用指针接收者可以修改接收者的值。

```
type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++
}

func main() {
    c := Counter{}
    c.Increment()
    fmt.Println(c.count) // 输出：1
}
```

## 7. 函数文档

使用注释可以为函数添加文档。

```
// Add returns the sum of two integers.
func Add(a, b int) int {
    return a + b
}
```

## 8.高阶技巧

### 将函数作为参数传入

将函数作为参数传入另一个函数，可以实现策略模式（Strategy Pattern）、回调（Callback）等功能。

#### 示例：策略模式

定义一个计算函数和不同的策略函数：

```
package main

import (
    "fmt"
)

// 定义一个计算函数类型
type Operation func(int, int) int

// 加法策略
func add(a, b int) int {
    return a + b
}

// 减法策略
func subtract(a, b int) int {
    return a - b
}

// 执行操作的高阶函数
func compute(a, b int, op Operation) int {
    return op(a, b)
}

func main() {
    fmt.Println("Addition:", compute(3, 4, add))        // 输出：Addition: 7
    fmt.Println("Subtraction:", compute(3, 4, subtract)) // 输出：Subtraction: -1
}
```

#### 示例：回调函数

使用回调函数来处理操作完成后的行为：

```
package main

import (
    "fmt"
)

// 定义一个回调函数类型
type Callback func(string)

// 处理操作的高阶函数
func process(data string, callback Callback) {
    // 模拟一些处理操作
    fmt.Println("Processing:", data)
    // 调用回调函数
    callback(data)
}

// 回调函数
func printResult(result string) {
    fmt.Println("Result:", result)
}

func main() {
    process("example data", printResult) // 输出：Processing: example data, Result: example data
}
```

### 将函数作为返回值

将函数作为返回值，可以实现工厂模式（Factory Pattern）、闭包（Closure）等功能。

#### 示例：工厂函数

工厂函数返回不同的策略函数：

```
package main

import (
    "fmt"
)

// 定义一个操作函数类型
type Operation func(int, int) int

// 返回加法函数的工厂函数
func getAddOperation() Operation {
    return func(a, b int) int {
        return a + b
    }
}

// 返回减法函数的工厂函数
func getSubtractOperation() Operation {
    return func(a, b int) int {
        return a - b
    }
}

func main() {
    add := getAddOperation()
    subtract := getSubtractOperation()

    fmt.Println("Addition:", add(3, 4))        // 输出：Addition: 7
    fmt.Println("Subtraction:", subtract(3, 4)) // 输出：Subtraction: -1
}
```

#### 示例：闭包

闭包是一种函数，引用了其外部作用域中的变量。通过闭包可以创建带状态的函数。

```
package main

import (
    "fmt"
)

// 返回一个计数器函数
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    c := counter()
    fmt.Println(c()) // 输出：1
    fmt.Println(c()) // 输出：2
    fmt.Println(c()) // 输出：3
}
```

## 闭包

### 闭包的基本概念

在Go语言中，任何匿名函数都可以成为闭包。匿名函数是没有名字的函数，可以在函数内部定义，并可以捕获其外部函数中的变量。

### 闭包的特性

1. **捕获环境**：闭包可以捕获并保存其定义时所在环境的变量。
2. **延长变量生命周期**：闭包可以延长局部变量的生命周期，确保在闭包被调用时这些变量依然存在。
3. **状态保持**：闭包可以用来创建拥有内部状态的函数，这些状态可以在多次调用中保持。

### 示例

以下是一个简单的闭包示例：

```
package main

import "fmt"

func main() {
    adder := createAdder(5)
    fmt.Println(adder(2)) // 输出 7
    fmt.Println(adder(10)) // 输出 15
}

func createAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}
```

### 解释

1. **`createAdder` 函数**：这个函数返回一个闭包。它接收一个整数参数 `x`，并返回一个匿名函数。匿名函数接收一个整数参数 `y`，并返回 `x` 和 `y` 的和。
2. **捕获环境**：在 `createAdder` 函数返回的闭包中，`x` 是 `createAdder` 函数的参数。即使 `createAdder` 函数已经执行完毕，`x` 仍然会被闭包捕获并保存在闭包中。
3. **使用闭包**：在 `main` 函数中，我们调用 `createAdder(5)` 得到一个闭包 `adder`。当我们调用 `adder(2)` 时，它会将 `5` 和 `2` 相加并返回 `7`。同理，调用 `adder(10)` 会返回 `15`。

### 闭包在实际开发中的应用

闭包在实际开发中非常有用，特别是在以下场景中：

1. **回调函数**：闭包可以作为回调函数，捕获并保存执行上下文中的变量。
2. **延迟计算**：闭包可以用来延迟计算，将计算逻辑和执行环境一起保存下来。
3. **数据封装**：闭包可以用来封装数据和方法，形成更高级的抽象。

### 示例：使用闭包实现一个简单的计数器

```
package main

import "fmt"

func main() {
    counter := createCounter()
    fmt.Println(counter()) // 输出 1
    fmt.Println(counter()) // 输出 2
    fmt.Println(counter()) // 输出 3
}

func createCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

### 解释

1. **`createCounter` 函数**：这个函数返回一个闭包。闭包内部定义了一个局部变量 `count`，并返回一个匿名函数。每次调用匿名函数时，`count` 都会递增并返回。
2. **状态保持**：闭包捕获并保存了 `count` 变量，因此每次调用 `counter` 函数时，`count` 的值会递增。