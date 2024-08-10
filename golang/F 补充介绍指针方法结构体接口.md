### 指针

#### 指针基础

指针是指向变量内存地址的变量，可以用来直接访问和修改内存中的值。

```
package main

import "fmt"

func main() {
    var a int = 10
    var p *int = &a  // 指针 p 指向 a 的地址

    fmt.Println(a)   // 输出：10
    fmt.Println(p)   // 输出：变量 a 的内存地址
    fmt.Println(*p)  // 输出：10（通过指针解引用访问变量 a 的值）
}
```

#### 使用指针修改变量

通过指针参数修改函数外部的变量。

```
package main

import "fmt"

func modify(a *int) {
    *a = 20
}

func main() {
    var x int = 10
    modify(&x)
    fmt.Println(x)  // 输出：20
}
```

### 结构体

#### 结构体基础

结构体是将多个不同类型的字段聚合在一起的数据结构。

```
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    fmt.Println(p)  // 输出：{Alice 30}
}
```

#### 结构体作为函数参数

将结构体作为参数传递给函数时，传递的是结构体的副本。如果需要在函数中修改结构体的字段，应该传递结构体指针。

```
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func updatePerson(p Person) {
    p.Age = 40
}

func updatePersonPointer(p *Person) {
    p.Age = 40
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    
    updatePerson(p)
    fmt.Println(p)  // 输出：{Alice 30}（未修改）

    updatePersonPointer(&p)
    fmt.Println(p)  // 输出：{Alice 40}（已修改）
}
```

#### 结构体指针

结构体指针可以直接指向结构体，并且允许修改结构体的字段。

```
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    p := &Person{Name: "Alice", Age: 30}
    p.Age = 40
    fmt.Println(p)  // 输出：&{Alice 40}
}
```

#### 匿名结构体

匿名结构体可以用于定义只在特定范围内使用的临时数据结构。

```
package main

import "fmt"

func main() {
    p := struct {
        Name string
        Age  int
    }{
        Name: "Alice",
        Age:  30,
    }

    fmt.Println(p)  // 输出：{Alice 30}
}
```

### 方法

#### 方法定义

方法是附加在某个类型上的函数，方法接收者可以是**值类型或者指针类型**。

```
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() {
    fmt.Printf("Hello, my name is %s and I am %d years old.\n", p.Name, p.Age)
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    p.Greet()  // 输出：Hello, my name is Alice and I am 30 years old.
}
```

#### 指针接收者方法

使用指针接收者可以修改结构体实例的字段，同时避免值类型方法传递时的内存开销。

```
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func (p *Person) SetAge(age int) {
    p.Age = age
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    p.SetAge(35)
    fmt.Println(p.Age)  // 输出：35
}
```

#### 方法的注意事项

1. **值接收者与指针接收者的选择**：如果方法不需要修改接收者的值，使用值接收者。如果需要修改接收者的值或接收者是一个大的结构体，使用指针接收者。
2. **方法与函数的区别**：方法附加在类型上，函数则是独立的。

### 接口

#### 接口定义

接口是一组方法签名的集合，用于定义行为。

```
package main

import "fmt"

type Speaker interface {
    Speak() string
}

type Person struct {
    Name string
}

func (p Person) Speak() string {
    return "Hello, my name is " + p.Name
}

func main() {
    var s Speaker
    s = Person{Name: "Alice"}
    fmt.Println(s.Speak())  // 输出：Hello, my name is Alice
}
```

#### 接口组合

接口可以嵌套组合，形成更复杂的接口。

```
package main

import "fmt"

type Reader interface {
    Read() string
}

type Writer interface {
    Write(string)
}

type ReadWriter interface {
    Reader
    Writer
}

type Data struct {
    content string
}

func (d *Data) Read() string {
    return d.content
}

func (d *Data) Write(content string) {
    d.content = content
}

func main() {
    var rw ReadWriter = &Data{}
    rw.Write("Hello, Go!")
    fmt.Println(rw.Read())  // 输出：Hello, Go!
}
```

#### 空接口

空接口 `interface{}` 可以表示任何类型的值。

```
package main

import "fmt"

func PrintValue(value interface{}) {
    fmt.Println(value)
}

func main() {
    PrintValue(42)         // 输出：42
    PrintValue("Hello")    // 输出：Hello
    PrintValue(3.14)       // 输出：3.14
}
```

#### 类型断言

类型断言用于从接口类型转换回具体类型。

```
package main

import "fmt"

func main() {
    var value interface{} = "Hello, Go!"
    if str, ok := value.(string); ok {
        fmt.Println(str)  // 输出：Hello, Go!
    } else {
        fmt.Println("Not a string")
    }
}
```

#### 类型选择

类型选择用于根据接口值的具体类型执行不同的操作。

```
package main

import "fmt"

func PrintType(value interface{}) {
    switch v := value.(type) {
    case int:
        fmt.Println("Integer:", v)
    case string:
        fmt.Println("String:", v)
    case float64:
        fmt.Println("Float64:", v)
    default:
        fmt.Println("Unknown type")
    }
}

func main() {
    PrintType(42)         // 输出：Integer: 42
    PrintType("Hello")    // 输出：String: Hello
    PrintType(3.14)       // 输出：Float64: 3.14
}
```

### 高级用法

#### 使用接口实现多态

多态允许一个接口表示多个不同的具体类型，通过这种方式可以实现灵活的设计。

```
package main

import "fmt"

type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14 * c.Radius * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func main() {
    shapes := []Shape{Circle{Radius: 5}, Rectangle{Width: 4, Height: 5}}
    for _, shape := range shapes {
        fmt.Println("Area:", shape.Area())
    }
    // 输出：
    // Area: 78.5
    // Area: 20
}
```

#### 使用闭包和高阶函数

闭包和高阶函数允许函数作为一等公民，提升代码的抽象程度。

```
package main

import "fmt"

func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

func main() {
    pos, neg := adder(), adder()
    for i := 0; i < 10; i++ {
        fmt.Println(
            pos(i),
            neg(-2*i),
        )
    }
}
```