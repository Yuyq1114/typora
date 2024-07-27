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