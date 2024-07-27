# 一 简介

Go语言（又称Golang）是一种由Google开发的开源编程语言，设计用于简化高性能和高可靠性软件的开发。

**并发支持**: 通过goroutine（类似于线程但更轻量级）和channel（用于goroutine之间的通信）来实现并发编程，

**简洁明快**: Go语言设计简洁，语法清晰，支持垃圾回收

**跨平台**: Go语言支持多种操作系统，包括Windows、Linux、Mac OS等，

**静态类型和编译型语言**: Go语言是静态类型语言，编译型语言，

# 二 杂项

## 1 go常用命令

**go mod init**: 初始化一个新的模块（在Go 1.11及以后版本中使用）。例如：`go mod init example.com/myproject`

**go mod tidy**: 根据go.mod文件整理模块依赖，更新 `go.sum` 文件。例如：`go mod tidy`

**go build**: 编译当前目录下的所有Go文件或指定的Go文件。生成可执行文件。例如：`go build` 或 `go build filename.go`

**go run**: 编译并直接运行一个Go源码文件。例如：`go run filename.go`

**go install**: 编译并安装当前包及其依赖包。安装后生成可执行文件或库文件。例如：`go install`

**go mod vendor：**命令用于将项目中所有依赖的模块（包括直接依赖和间接依赖）复制到项目的 `vendor` 目录下。这些依赖将会在 `vendor` 目录中被维护，使得项目在构建时可以使用这些本地存储的依赖，而不需要从网络上下载。

**go get**: 下载并安装指定的包及其依赖包。例如：`go get github.com/example/package`

**go test**: 运行当前包中的测试文件。例如：`go test`

**go fmt**: 格式化源代码文件，使其符合Go语言规范。例如：`go fmt filename.go`

**go vet**: 分析Go源码中的静态错误。例如：`go vet`

**go version**: 查看当前安装的Go版本。例如：`go version`

## 2 快速代码段

**ff** `fmt.Printf("", var)`

**pkgm** `package main func main() {}`

**if** `if condition {}`

**iferr** `if err != nil { return nil, err }`

**for** `for i := 0; i < count; i++ {}`

## 3 步骤

**创建项目文件夹**：

`mkdir mygoapp
cd mygoapp`

**初始化Go模块**：

`go mod init mygoapp`

**创建main.go文件**：

`touch main.go`

**编辑main.go文件**：

**运行程序**：

`go run main.go`

# 三 第二章 标识符、关键字、命名规则

## 1标识符

**组成**：

- 标识符由字母、数字和下划线 `_` 组成。
- 第一个字符必须是字母（大小写均可）或下划线 `_`，不能是数字或其他字符。

**大小写敏感**：

- Go语言中的标识符是大小写敏感的，例如 `myVar` 和 `MyVar` 是两个不同的标识符。

**Unicode字符**：

- 标识符可以使用Unicode字符（比如非ASCII字符），但通常不建议过度使用，以保持代码的可读性和兼容性。

## 2 关键字

`break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var`

## 3 命名约定

**驼峰命名法**：

- Go语言推荐使用驼峰命名法（CamelCase）来命名变量、函数和方法。例如：`myVariableName`、`calculateInterestRate`。

**包名**：

- 包名通常采用小写字母，短小且有描述性，比如 `fmt`、`math`。

**导出标识符**：

- 如果标识符以大写字母开头，则被视为导出标识符（Exported Identifier），可以被外部包（package）访问和使用。例如：`MyPublicFunction`。

**非导出标识符**：

- 如果标识符以小写字母开头，则被视为非导出标识符（Unexported Identifier），只能在定义它的包内部使用。例如：`myPrivateFunction`。