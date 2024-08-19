Go 的 `testing` 模块是 Go 标准库的一部分，专门用于编写单元测试、性能基准测试和示例代码。该模块提供了开发者用于验证代码功能性、性能及行为的工具。

### 核心功能

- **单元测试**：验证函数或方法的正确性。
- **基准测试**：测试函数或方法的性能表现。
- **示例代码**：提供文档化的代码示例，可以直接运行并验证输出。

### 基本结构

`testing` 包的核心结构包括以下几种：

- **`testing.T`**: 用于单元测试，包含控制测试流程的方法。
- **`testing.B`**: 用于性能基准测试，提供用于控制和报告性能测试的机制。
- **`testing.M`**: 用于在单元测试或基准测试之前设置全局初始化。

### 单元测试

在 Go 中，单元测试文件的名称通常以 `_test.go` 结尾，并放置在与被测试代码相同的包中。测试函数的名称以 `Test` 开头，并接受一个 `*testing.T` 类型的参数。

#### 示例

```
package main

import "testing"

// 被测试的函数
func Add(a, b int) int {
    return a + b
}

// 单元测试函数
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

在这个例子中，`TestAdd` 是一个单元测试函数，用于验证 `Add` 函数的正确性。如果结果不符合预期，测试会报告错误。

### 运行测试

使用 Go 的测试工具，可以通过以下命令运行测试：

```
go test
```

它会自动查找以 `_test.go` 结尾的文件，运行其中的所有测试函数。

### 基准测试

基准测试用于评估代码的性能。基准测试函数的名称以 `Benchmark` 开头，并接受一个 `*testing.B` 类型的参数。基准测试函数中的代码会被重复执行 `b.N` 次，以便测量性能。

#### 示例

```
package main

import "testing"

// 基准测试函数
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

运行基准测试命令：

```
go test -bench=.
```

这会运行以 `Benchmark` 开头的函数，并报告性能结果。

### 示例代码

Go 的 `testing` 包还支持为文档提供示例代码，这些示例代码可以验证输出是否正确。示例代码函数的名称以 `Example` 开头，不需要传递参数。

#### 示例

```
package main

import "fmt"

func ExampleAdd() {
    fmt.Println(Add(2, 3))
    // Output: 5
}
```

在这个例子中，`ExampleAdd` 函数展示了如何使用 `Add` 函数，并通过注释 `// Output: 5` 验证示例的输出。

### 使用 `TestMain`

有时候在执行测试之前需要进行全局设置或清理工作，可以使用 `TestMain` 函数。它允许你在测试运行之前和之后执行代码。

#### 示例

```
package main

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // 测试前的初始化
    setup()

    // 运行测试
    code := m.Run()

    // 测试后的清理
    teardown()

    // 退出
    os.Exit(code)
}
```

### 常用的 `testing.T` 方法

- **`t.Error(args...)`**: 记录测试错误信息并继续执行。
- **`t.Fail()`**: 标记测试失败，但继续执行。
- **`t.FailNow()`**: 标记测试失败，并立即停止当前测试函数。
- **`t.Fatal(args...)`**: 记录错误信息并停止测试。
- **`t.Skip(args...)`**: 跳过测试。
- **`t.Parallel()`**: 允许测试并发执行。

### 表格驱动测试

表格驱动测试是一种在 Go 测试中常见的模式，适用于为同一个函数测试多种输入和输出。

#### 示例

```
package main

import "testing"

func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"add positive", 1, 2, 3},
        {"add zero", 0, 0, 0},
        {"add negative", -1, -2, -3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```