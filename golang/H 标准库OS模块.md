在Go语言标准库中，`os`模块是一个非常重要的包，提供了对操作系统底层功能的访问和控制。它涵盖了文件操作、进程管理、环境变量等多个方面，是进行系统级编程不可或缺的工具。下面详细介绍`os`模块中的主要功能和使用方法。

### 文件和目录操作

`os`模块提供了多种文件和目录操作函数，用于创建、打开、删除文件，以及管理目录。

#### 创建和打开文件

- `os.Create(name string) (*File, error)`: 创建一个文件，返回一个文件对象。
- `os.Open(name string) (*File, error)`: 打开一个文件进行读取操作。
- `os.OpenFile(name string, flag int, perm FileMode) (*File, error)`: 根据指定的标志打开文件，支持更多的选项。

```
package main

import (
    "fmt"
    "os"
)

func main() {
    // 创建文件
    f, err := os.Create("test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()

    // 打开文件进行读取
    f, err = os.Open("test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
}
```

#### 文件信息和操作

- `os.Stat(name string) (FileInfo, error)`: 返回一个文件对象的信息。
- `os.Remove(name string) error`: 删除指定的文件。
- `os.RemoveAll(path string) error`: 递归删除指定的目录及其所有内容。

```
package main

import (
    "fmt"
    "os"
)

func main() {
    // 获取文件信息
    info, err := os.Stat("test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("File size:", info.Size())

    // 删除文件
    err = os.Remove("test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

#### 目录操作

- `os.Mkdir(name string, perm FileMode) error`: 创建一个目录。
- `os.MkdirAll(path string, perm FileMode) error`: 创建一个路径中所有不存在的目录。

```
package main

import (
    "fmt"
    "os"
)

func main() {
    // 创建目录
    err := os.Mkdir("mydir", 0755)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 创建多层目录
    err = os.MkdirAll("mydir/subdir", 0755)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

### 环境变量和进程管理

`os`模块还提供了访问环境变量和执行进程的功能。

#### 环境变量

- `os.Getenv(key string) string`: 获取指定环境变量的值。
- `os.Setenv(key, value string) error`: 设置指定环境变量的值。
- `os.Unsetenv(key string) error`: 删除指定环境变量。

```
package main

import (
    "fmt"
    "os"
)

func main() {
    // 获取环境变量
    path := os.Getenv("PATH")
    fmt.Println("Path:", path)

    // 设置环境变量
    err := os.Setenv("MY_VAR", "my_value")
    if err != nil {
        fmt.Println(err)
        return
    }

    // 删除环境变量
    err = os.Unsetenv("MY_VAR")
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

#### 进程管理

- `os.Exit(code int)`: 终止当前进程。
- `os.Getpid() int`: 获取当前进程的PID。
- `os.Getppid() int`: 获取父进程的PID。

```
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("PID:", os.Getpid())
    fmt.Println("PPID:", os.Getppid())

    // 退出当前进程
    os.Exit(1)
}
```

### 其他功能

除了上述功能之外，`os`模块还提供了其他一些有用的函数，例如：

- `os.Hostname() (name string, err error)`: 获取主机名。
- `os.Getwd() (dir string, err error)`: 获取当前工作目录。
- `os.Chdir(dir string) error`: 改变当前工作目录。

```
package main

import (
    "fmt"
    "os"
)

func main() {
    // 获取主机名
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Hostname:", hostname)

    // 获取当前工作目录
    wd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Current working directory:", wd)
}
```