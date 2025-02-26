## 一 IO相关模块

### 1. io 包

`io`包定义了基本的IO接口，用于通用的数据读取和写入操作。

#### 基本接口

- `io.Reader`：读取数据的接口，定义了`Read`方法。
- `io.Writer`：写入数据的接口，定义了`Write`方法。
- `io.Closer`：关闭资源的接口，定义了`Close`方法。
- `io.Seeker`：定位操作的接口，定义了`Seek`方法。

这些接口被广泛应用于Go语言中的各种IO操作，使得不同类型的数据源（文件、网络连接、内存等）可以统一通过相同的方式进行处理。

#### 示例：使用`io.Reader`和`io.Writer`

```
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    // 使用strings.NewReader创建一个io.Reader
    reader := strings.NewReader("Hello, Go!")

    // 从reader中读取数据并写入到标准输出
    io.Copy(os.Stdout, reader)
}
```

#### 错误处理

`io`包中的函数通常返回一个`error`类型，用于表示操作是否成功，需要进行错误处理。

```
go复制代码package main

import (
    "fmt"
    "io/ioutil"
)

func main() {
    // 读取文件内容
    data, err := ioutil.ReadFile("file.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    fmt.Println("File contents:", string(data))
}
```

### 2. bufio 包

`bufio`包实现了带缓冲区的读写操作，可以提高IO操作的效率，特别是在处理大量小数据块的场景下。

#### 缓冲读取器

- `bufio.NewReader(rd io.Reader) *Reader`：创建一个新的带缓冲的Reader。
- `Reader.ReadLine() ([]byte, error)`：读取一行数据。
- `Reader.ReadBytes(delim byte) ([]byte, error)`：读取直到遇到分隔符的数据。

```
package main

import (
    "bufio"
    "fmt"
    "strings"
)

func main() {
    // 使用strings.NewReader创建一个io.Reader
    reader := strings.NewReader("Line 1\nLine 2\nLine 3\n")

    // 创建一个带缓冲的Reader
    bufReader := bufio.NewReader(reader)

    // 逐行读取数据
    for {
        line, err := bufReader.ReadString('\n')
        if err != nil {
            break
        }
        fmt.Print(line)
    }
}
```

#### 缓冲写入器

- `bufio.NewWriter(wr io.Writer) *Writer`：创建一个新的带缓冲的Writer。
- `Writer.WriteString(s string) (int, error)`：写入字符串到缓冲区。
- `Writer.Flush() error`：将缓冲区的数据写入底层的io.Writer。

```
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // 创建一个带缓冲的Writer
    file, err := os.Create("output.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    // 写入数据到缓冲区
    _, err = writer.WriteString("Hello, Go!\n")
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    // 将缓冲区的数据写入文件
    err = writer.Flush()
    if err != nil {
        fmt.Println("Error flushing buffer:", err)
        return
    }
}
```

### 3. ioutil 包

`ioutil`包提供了一些便利的IO函数，尤其是对于一次性简单的IO操作。

- `ioutil.ReadFile(filename string) ([]byte, error)`：读取整个文件的内容。
- `ioutil.WriteFile(filename string, data []byte, perm os.FileMode) error`：将数据写入文件。

```
package main

import (
    "fmt"
    "io/ioutil"
)

func main() {
    // 读取文件内容
    data, err := ioutil.ReadFile("file.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    fmt.Println("File contents:", string(data))

    // 写入数据到文件
    err = ioutil.WriteFile("output.txt", []byte("Hello, Go!"), 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }
}
```



## 二 bytes

Go语言标准库中的`bytes`包提供了对字节切片（`[]byte`）的操作，它包含了许多用于操作字节切片的函数和方法。`bytes`包特别适合于需要高效操作字节数据的场景，例如数据的拼接、分割、搜索、替换等操作。下面详细介绍`bytes`包中常用的函数和方法。

### 1. 基本函数和变量

- `func Contains(b, subslice []byte) bool`：判断字节切片 `b` 是否包含子切片 `subslice`。
- `func Count(s, sep []byte) int`：计算字节切片 `s` 中子切片 `sep` 的非重叠实例的数量。
- `func Equal(a, b []byte) bool`：比较两个字节切片 `a` 和 `b` 是否相等。
- `func Index(s, sep []byte) int`：返回子切片 `sep` 在字节切片 `s` 中第一次出现的索引，如果未找到返回 -1。
- `func Join(s [][]byte, sep []byte) []byte`：连接多个字节切片 `s` 并用 `sep` 分隔。
- `func Repeat(b []byte, count int) []byte`：将字节切片 `b` 重复 `count` 次。
- `func Replace(s, old, new []byte, n int) []byte`：替换字节切片 `s` 中前 `n` 个 `old` 子切片为 `new` 子切片。
- `func Split(s, sep []byte) [][]byte`：根据子切片 `sep` 分割字节切片 `s`。
- `func Trim(s []byte, cutset string) []byte`：去掉字节切片 `s` 开头和结尾处在 `cutset` 中的所有字符。

### 示例

```
package main

import (
    "bytes"
    "fmt"
)

func main() {
    // Contains示例
    fmt.Println("Contains:", bytes.Contains([]byte("seafood"), []byte("foo"))) // true

    // Count示例
    fmt.Println("Count:", bytes.Count([]byte("cheese"), []byte("e"))) // 3

    // Equal示例
    fmt.Println("Equal:", bytes.Equal([]byte("Go"), []byte("Go"))) // true

    // Index示例
    fmt.Println("Index:", bytes.Index([]byte("gopher"), []byte("he"))) // 3

    // Join示例
    fmt.Println("Join:", string(bytes.Join([][]byte{[]byte("a"), []byte("b"), []byte("c")}, []byte(", ")))) // a, b, c

    // Repeat示例
    fmt.Println("Repeat:", string(bytes.Repeat([]byte("na"), 2))) // nana

    // Replace示例
    fmt.Println("Replace:", string(bytes.Replace([]byte("oink oink oink"), []byte("oink"), []byte("moo"), 2))) // moo moo oink

    // Split示例
    fmt.Printf("Split: %q\n", bytes.Split([]byte("a,b,c"), []byte(","))) // ["a" "b" "c"]

    // Trim示例
    fmt.Printf("Trim: %s\n", bytes.Trim([]byte("¡¡¡Hello, Gophers!!!"), "!¡")) // Hello, Gophers
}
```

### 2. Buffer 类型和方法

除了上述函数外，`bytes`包中还定义了`Buffer`类型，提供了更灵活和高效的字节缓冲区操作。

#### Buffer 类型方法

- `func NewBuffer(buf []byte) *Buffer`：创建一个新的字节缓冲区，可以指定初始内容 `buf`。
- `func (b *Buffer) Bytes() []byte`：返回缓冲区的内容作为一个字节切片。
- `func (b *Buffer) String() string`：返回缓冲区的内容作为一个字符串。
- `func (b *Buffer) Len() int`：返回缓冲区中未读取数据的长度。
- `func (b *Buffer) Cap() int`：返回缓冲区的容量。
- `func (b *Buffer) Reset()`：重置缓冲区，清空所有内容但保留底层存储空间。
- `func (b *Buffer) Write(p []byte) (n int, err error)`：将字节切片 `p` 写入缓冲区。
- `func (b *Buffer) WriteByte(c byte) error`：将单个字节 `c` 写入缓冲区。
- `func (b *Buffer) WriteString(s string) (n int, err error)`：将字符串 `s` 写入缓冲区。

#### 示例：使用 Buffer 类型

```
package main

import (
    "bytes"
    "fmt"
)

func main() {
    // 创建一个新的 Buffer
    var buf bytes.Buffer

    // 写入数据到 Buffer
    buf.WriteString("Hello, ")
    buf.WriteByte('G')
    buf.WriteByte('o')
    buf.WriteString("pher!")

    // 输出 Buffer 的内容
    fmt.Println("Buffer contents:", buf.String())
}
```

## 三 sort

### 基本函数和类型

#### 1. 切片排序

- `func Ints(a []int)`：对整数切片 `a` 进行升序排序。
- `func Float64s(a []float64)`：对 `float64` 类型的切片 `a` 进行升序排序。
- `func Strings(a []string)`：对字符串切片 `a` 进行升序排序。
- `func Reverse(data Interface)`：对实现了 `Interface` 接口的切片进行降序排序。

#### 示例：对切片进行排序

```
package main

import (
    "fmt"
    "sort"
)

func main() {
    // 对整数切片进行排序
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
    sort.Ints(nums)
    fmt.Println("Sorted ints:", nums)

    // 对字符串切片进行排序
    strs := []string{"apple", "orange", "banana", "pear"}
    sort.Strings(strs)
    fmt.Println("Sorted strings:", strs)

    // 对自定义类型切片进行排序
    type Person struct {
        Name string
        Age  int
    }

    people := []Person{
        {"Alice", 25},
        {"Bob", 30},
        {"Charlie", 20},
    }

    // 按照年龄升序排序
    sort.Slice(people, func(i, j int) bool {
        return people[i].Age < people[j].Age
    })

    fmt.Println("Sorted people by age:", people)
}
```

#### 2. 自定义排序

如果要对自定义类型进行排序，需要实现 `sort.Interface` 接口的三个方法：`Len()`、`Less(i, j int) bool` 和 `Swap(i, j int)`。

```
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

这使得我们可以对任何类型的数据进行排序，只要实现了上述三个方法。

#### 示例：自定义类型排序

```
package main

import (
    "fmt"
    "sort"
)

type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
    people := []Person{
        {"Alice", 25},
        {"Bob", 30},
        {"Charlie", 20},
    }

    // 使用自定义的排序
    sort.Sort(ByAge(people))
    fmt.Println("Sorted people by age:", people)
}
```

### 注意事项

- 对切片进行排序会改变切片本身的顺序，而不是返回一个新的切片。
- 自定义排序时，确保实现了 `sort.Interface` 接口的所有方法，并正确实现 `Less` 方法的比较逻辑。

## 四 math

### 常用函数

#### 1. 基本函数

- `func Abs(x float64) float64`：返回 `x` 的绝对值。
- `func Ceil(x float64) float64`：返回大于或等于 `x` 的最小整数。
- `func Floor(x float64) float64`：返回小于或等于 `x` 的最大整数。
- `func Round(x float64) float64`：返回 `x` 四舍五入的整数。
- `func Max(x, y float64) float64`：返回 `x` 和 `y` 中的最大值。
- `func Min(x, y float64) float64`：返回 `x` 和 `y` 中的最小值。
- `func Pow(x, y float64) float64`：返回 `x` 的 `y` 次幂。

#### 示例

```
package main

import (
    "fmt"
    "math"
)

func main() {
    x := -10.5
    y := 8.3

    // 绝对值
    fmt.Println("Abs:", math.Abs(x))

    // 向上取整
    fmt.Println("Ceil:", math.Ceil(x))

    // 向下取整
    fmt.Println("Floor:", math.Floor(x))

    // 四舍五入
    fmt.Println("Round:", math.Round(x))

    // 最大值和最小值
    fmt.Println("Max:", math.Max(x, y))
    fmt.Println("Min:", math.Min(x, y))

    // 幂运算
    fmt.Println("Pow:", math.Pow(x, y))
}
```

#### 2. 三角函数

- `func Sin(x float64) float64`：返回角度 `x` 的正弦值。
- `func Cos(x float64) float64`：返回角度 `x` 的余弦值。
- `func Tan(x float64) float64`：返回角度 `x` 的正切值。

#### 示例

```
package main

import (
    "fmt"
    "math"
)

func main() {
    angle := 45.0
    radians := angle * math.Pi / 180.0

    fmt.Printf("Sin(%f) = %f\n", angle, math.Sin(radians))
    fmt.Printf("Cos(%f) = %f\n", angle, math.Cos(radians))
    fmt.Printf("Tan(%f) = %f\n", angle, math.Tan(radians))
}
```

### 常量

`math`包还定义了一些常用的数学常量，如π和自然对数的底数：

- `math.Pi`：圆周率π的近似值。
- `math.E`：自然对数的底数e的近似值。

```
package main

import (
    "fmt"
    "math"
)

func main() {
    fmt.Println("Pi:", math.Pi)
    fmt.Println("E:", math.E)
}
```

### 注意事项

- 对于复杂的数学运算（如精确的浮点数计算），应该谨慎处理浮点数精度问题。
- 在使用三角函数时，需要将角度转换为弧度（弧度 = 角度 * π / 180）。

## 五 json

### 基本函数和类型

#### 1. 结构体和函数

- `func Marshal(v interface{}) ([]byte, error)`：将 Go 中的数据结构 `v` 编码为 JSON 格式的字节切片。
- `func Unmarshal(data []byte, v interface{}) error`：将 JSON 格式的字节切片 `data` 解码为 Go 中的数据结构 `v`。

这两个函数是 `encoding/json` 包中最基本和最常用的函数，分别用于 JSON 编码和解码。

#### 示例：编码和解码

```
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // 编码 JSON
    person := Person{"Alice", 30}
    encoded, err := json.Marshal(person)
    if err != nil {
        fmt.Println("JSON encode error:", err)
        return
    }
    fmt.Println("Encoded JSON:", string(encoded))

    // 解码 JSON
    var decoded Person
    if err := json.Unmarshal(encoded, &decoded); err != nil {
        fmt.Println("JSON decode error:", err)
        return
    }
    fmt.Println("Decoded Person:", decoded)
}
```

#### 2. 标签（Tags）

在结构体中，可以使用标签（Tags）来控制 JSON 编码和解码过程中的字段名称和行为。标签通过 `json:"tagname"` 的形式添加到结构体字段上，例如 `Name string `json:"name"`。这样可以在 JSON 数据和 Go 结构体之间进行灵活的映射。

#### 3. 其他函数

- `func NewDecoder(r io.Reader) *Decoder`：创建一个新的解码器，从 `io.Reader` 中读取 JSON 数据。
- `func NewEncoder(w io.Writer) *Encoder`：创建一个新的编码器，将 JSON 数据写入到 `io.Writer` 中。

### 注意事项

- JSON 解码时，要确保目标结构体字段是可导出的（即首字母大写），以便 `encoding/json` 包能够访问和设置这些字段。
- 在处理嵌套和复杂的 JSON 结构时，需要仔细设计和检查数据的结构和类型，以确保正确的解码和编码。

## 六 xml

### 基本函数和类型

#### 1. 结构体和函数

- `func Marshal(v interface{}) ([]byte, error)`：将Go语言中的数据结构 `v` 编码为XML格式的字节切片。
- `func Unmarshal(data []byte, v interface{}) error`：将XML格式的字节切片 `data` 解码为Go语言中的数据结构 `v`。

这两个函数与`encoding/json`包中的函数功能类似，用于XML的编码和解码。

#### 示例：编码和解码

```
package main

import (
	"encoding/xml"
	"fmt"
)

type Person struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
}

func main() {
	// 编码 XML
	person := Person{Name: "Alice", Age: 30}
	encoded, err := xml.Marshal(person)
	if err != nil {
		fmt.Println("XML encode error:", err)
		return
	}
	fmt.Println("Encoded XML:", string(encoded))

	// 解码 XML
	var decoded Person
	if err := xml.Unmarshal(encoded, &decoded); err != nil {
		fmt.Println("XML decode error:", err)
		return
	}
	fmt.Println("Decoded Person:", decoded)
}
```

#### 2. 标签（Tags）

与JSON的标签类似，XML编码和解码中也可以使用标签来指定XML元素的名称和属性。

- `xml:"elementname"`：指定XML元素的名称。
- `xml:",attr"`：将字段作为XML元素的属性。

#### 示例：使用标签

```
type Book struct {
	XMLName  xml.Name `xml:"book"`
	Title    string   `xml:"title"`
	Author   string   `xml:"author"`
	Pages    int      `xml:"pages"`
	Published xml.CharData `xml:"published"`
}

func main() {
	book := Book{
		Title:    "The Go Programming Language",
		Author:   "Alan Donovan & Brian Kernighan",
		Pages:    380,
		Published: xml.CharData("2015-10-26"),
	}

	encoded, err := xml.MarshalIndent(book, "", "  ")
	if err != nil {
		fmt.Println("XML encode error:", err)
		return
	}

	fmt.Println(string(encoded))
}
```

#### 3. Decoder 和 Encoder 类型

除了 `Marshal` 和 `Unmarshal` 函数外，`encoding/xml` 包还提供了 `Decoder` 和 `Encoder` 类型，可以用于更灵活的 XML 数据流操作。

- `func NewDecoder(r io.Reader) *Decoder`：创建一个新的 XML 解码器，从 `io.Reader` 中读取 XML 数据。
- `func NewEncoder(w io.Writer) *Encoder`：创建一个新的 XML 编码器，将 XML 数据写入到 `io.Writer` 中。

### 注意事项

- 在使用 `encoding/xml` 包进行编码和解码时，确保目标结构体字段是可导出的（即首字母大写），以便包可以访问和设置这些字段。
- 要注意处理嵌套和复杂的 XML 结构时，需要仔细设计和检查数据的结构和类型，以确保正确的解码和编码。

## 七 time库

### 基本函数和类型

#### 1. 时间类型

- `type Time`：表示一个具体的时刻，包含年、月、日、时、分、秒、纳秒和时区信息。

#### 2. 获取当前时间

- `func Now() Time`：返回当前本地时间。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("Current time:", now)
}
```

#### 3. 格式化时间

- `func (t Time) Format(layout string) string`：根据指定的格式 `layout` 将时间格式化为字符串。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("Formatted time:", now.Format("2006-01-02 15:04:05"))
}
```

在格式化时间时，`layout` 中的特定格式组合（如 "2006-01-02 15:04:05"）是 Go 语言中的约定形式，用于表示年、月、日、时、分、秒的固定格式。

#### 4. 解析时间字符串

- `func Parse(layout, value string) (Time, error)`：根据指定的格式 `layout` 解析时间字符串 `value`。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	str := "2024-07-15 10:30:00"
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	fmt.Println("Parsed time:", t)
}
```

### 其他常用函数

- `func Sleep(d Duration)`：休眠指定的时间段 `d`。
- `func Since(t Time) Duration`：返回从时间 `t` 到当前时间的时间段。
- `func Until(t Time) Duration`：返回从当前时间到时间 `t` 的时间段。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	time.Sleep(2 * time.Second)
	elapsed := time.Since(start)
	fmt.Println("Elapsed time:", elapsed)
}
```

### 时间间隔（Duration）

- `type Duration`：表示时间段，单位为纳秒，类型为 `int64`。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	d := 10 * time.Second
	fmt.Println("Duration:", d)
}
```

### 时区和位置

`time` 包中的时间类型 `Time` 包含了时区信息，可以进行时区转换和处理。

### 注意事项

- 在处理时间和日期时，尤其是跨时区的应用中，应当注意使用正确的时区和格式化方法，以避免潜在的错误和混淆。



## 八 Context包

### 1. **背景与设计目的**

Go 的 `context` 包是为了解决以下问题而设计的：

- **控制协程的生命周期**：在多协程环境下，能够优雅地通知某些协程结束。
- **超时与取消**：支持设置超时或取消信号，以避免协程无限运行或被挂起。
- **传递请求范围内的数据**：通过 `context` 可以在线程之间传递少量上下文信息，比如用户 ID、请求标识等。

------

### 2. **`context` 的核心接口**

`context` 是一个接口，其定义如下：

```
type Context interface {
    Deadline() (deadline time.Time, ok bool) // 返回上下文被取消的时间
    Done() <-chan struct{}                  // 返回一个通道，当上下文被取消或超时时关闭
    Err() error                             // 返回上下文结束的原因（取消或超时）
    Value(key any) any                      // 获取上下文中存储的值
}
```

#### **主要方法解析**

1. **`Deadline`**：用于获取上下文的截止时间。超出截止时间，`context` 会自动取消。
2. **`Done`**：返回一个只读通道，表示上下文的取消信号。当通道关闭时，所有依赖此 `context` 的操作都应该停止。
3. **`Err`**：当 `Done` 通道关闭时，可以通过此方法判断取消的原因（如超时、主动取消等）。
4. **`Value`**：用来存储和获取与上下文相关的值，适合存放一些与请求相关的小型全局变量。

------

### 3. **创建 `context` 的方式**

`context` 包中提供了以下四种方法来生成 `context`：

#### (1) `context.Background()`

- **用途**：返回一个空的上下文，通常作为顶层父上下文使用。
- **特点**：不会被取消，也没有超时时间，适合作为根上下文。

```
ctx := context.Background()
```

------

#### (2) `context.TODO()`

- **用途**：返回一个空的上下文，表示当前不知道要使用什么 `context`。
- **特点**：通常在代码未完善时作为占位符，待确定后再替换为适当的 `context`。

```
ctx := context.TODO()
```

------

#### (3) `context.WithCancel(parent)`

- **用途**：创建一个可取消的上下文。
- **特点**：返回的 `context` 可以通过调用 `cancel` 函数主动取消。

```
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // 确保在适当的地方调用取消函数

go func() {
    <-ctx.Done() // 监听取消信号
    fmt.Println("Context canceled")
}()
```

------

#### (4) `context.WithTimeout(parent, duration)`

- **用途**：创建一个具有超时时间的上下文。
- **特点**：超时后，`context` 会自动取消，并关闭其 `Done` 通道。

```
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-time.After(6 * time.Second):
    fmt.Println("Task completed")
case <-ctx.Done():
    fmt.Println("Context timeout:", ctx.Err()) // 超时后返回错误
}
```

------

#### (5) `context.WithDeadline(parent, deadline)`

- **用途**：类似于 `WithTimeout`，但允许指定具体的截止时间。
- **特点**：适合需要精确控制时间点的场景。

```
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

------

#### (6) `context.WithValue(parent, key, value)`

- **用途**：在上下文中存储一个键值对。
- **特点**：可以传递一些少量的请求范围内的元信息，但不适合存储大量数据。

```
ctx := context.WithValue(context.Background(), "userID", 12345)

userID := ctx.Value("userID").(int) // 获取存储的值
fmt.Println("UserID:", userID)
```

> **注意**：`context.Value` 不应该用于传递业务逻辑中的主要数据，只适合传递少量元数据。

------

### 4. **常见应用场景**

#### (1) **超时控制**

避免某些操作（如网络请求、数据库查询）无限等待，通过设置超时来控制任务的执行时间。

```
func fetchData(ctx context.Context) {
    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Data fetched")
    case <-ctx.Done():
        fmt.Println("Operation canceled:", ctx.Err())
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    fetchData(ctx)
}
```

------

#### (2) **协程间的取消控制**

```
func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopped\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    for i := 0; i < 3; i++ {
        go worker(ctx, i)
    }

    time.Sleep(3 * time.Second)
    cancel() // 通知所有协程停止
    time.Sleep(1 * time.Second)
}
```

------

#### (3) **传递请求范围数据**

在处理 HTTP 请求时，可以用 `context` 将请求范围内的元数据传递给下游函数。

```
func handler(ctx context.Context) {
    userID := ctx.Value("userID").(int)
    fmt.Println("Handling request for user:", userID)
}

func main() {
    ctx := context.WithValue(context.Background(), "userID", 42)
    handler(ctx)
}
```

------

### 5. **注意事项**

1. **及时调用 `cancel`**：
   - 当使用 `WithCancel`、`WithTimeout` 或 `WithDeadline` 时，一定要确保调用对应的 `cancel` 函数以释放资源。
2. **不要滥用 `WithValue`**：
   - `context` 是用于传递上下文信息，而不是数据存储工具。避免在上下文中传递过多的数据。
3. **嵌套上下文的管理**：
   - 上下文可以嵌套，子上下文的取消会级联到父上下文，但父上下文的取消会传播到所有子上下文。

## 九 log包

`log` 包是 Go 标准库中用于记录日志的包。它提供了简单而强大的日志记录功能，可以将日志信息输出到不同的目的地（如控制台、文件等），并支持设置日志的前缀和时间格式。

### 基本用法

#### 创建日志记录器

`log` 包提供了一个默认的日志记录器，可以直接使用 `log.Print`、`log.Println`、`log.Printf` 等方法记录日志。

```
log.Print("This is a log message.")
log.Println("This is a log message with a new line.")
log.Printf("This is a formatted log message: %s", "formatted")
```

#### 设置日志前缀和标志

可以使用 `log.SetPrefix` 设置日志的前缀，使用 `log.SetFlags` 设置日志的标志（时间戳、文件名、行号等）。

```
log.SetPrefix("INFO: ")
log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

log.Println("This is a log message with a prefix and flags.")
```

常用的标志包括：

- `log.Ldate`：日期（2009/01/23）
- `log.Ltime`：时间（01:23:23）
- `log.Lmicroseconds`：微秒级时间（01:23:23.123123）
- `log.Llongfile`：完整文件名和行号（/a/b/c/d.go:23）
- `log.Lshortfile`：文件名和行号（d.go:23）
- `log.LUTC`：使用 UTC 时间

### 高级用法

#### 自定义日志记录器

可以创建自定义的日志记录器，以便将日志信息输出到不同的目的地（如文件、网络等）。

```
package main

import (
    "log"
    "os"
)

func main() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %s", err)
    }
    defer file.Close()

    logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

    logger.Println("This is a log message written to a file.")
}
```

#### 并发安全

`log` 包中的默认日志记录器是并发安全的，可以在多个 goroutine 中安全地使用。如果创建自定义的日志记录器，默认情况下也是并发安全的。

#### 设置输出目标

可以使用 `log.SetOutput` 方法更改默认日志记录器的输出目标。

```
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatalf("Failed to open log file: %s", err)
}
defer file.Close()

log.SetOutput(file)

log.Println("This is a log message written to a file.")
```

#### 自定义日志格式

如果需要自定义日志的输出格式，可以实现 `io.Writer` 接口，并使用 `log.SetOutput` 将其设置为日志记录器的输出目标。

```
package main

import (
    "fmt"
    "log"
    "os"
    "time"
)

type customWriter struct {
    file *os.File
}

func (cw *customWriter) Write(p []byte) (n int, err error) {
    logMessage := fmt.Sprintf("%s: %s", time.Now().Format(time.RFC3339), string(p))
    return cw.file.Write([]byte(logMessage))
}

func main() {
    file, err := os.OpenFile("custom.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %s", err)
    }
    defer file.Close()

    cw := &customWriter{file: file}
    log.SetOutput(cw)

    log.Println("This is a custom formatted log message.")
}
```

### 日志级别

`log` 包本身不支持日志级别（如 INFO、DEBUG、WARN、ERROR），但可以通过创建不同的日志记录器来实现日志级别管理。

```
package main

import (
    "log"
    "os"
)

var (
    infoLogger  *log.Logger
    errorLogger *log.Logger
)

func init() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %s", err)
    }

    infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
    infoLogger.Println("This is an info message.")
    errorLogger.Println("This is an error message.")
}
```

### 示例项目

以下是一个完整的示例项目，演示了如何使用 `log` 包记录日志，设置日志前缀和标志，以及自定义日志记录器。

```
package main

import (
    "log"
    "os"
)

func main() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %s", err)
    }
    defer file.Close()

    log.SetOutput(file)
    log.SetPrefix("INFO: ")
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    log.Println("This is a log message written to a file.")

    errorLogger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger.Println("This is an error message.")
}
```

在这个示例中，日志记录器被配置为将日志信息输出到文件 `app.log`，并设置了日志前缀和标志。还创建了一个自定义的错误日志记录器，用于记录错误信息。

### 其他

| 日志级别 | 描述                                                         | 适用场景                                                     |
| -------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `DEBUG`  | **调试级别日志**：提供最详细的信息，通常用于开发和调试。     | - 输出调试信息，例如变量值、程序流程状态。 - 适合在开发阶段使用，生产环境中通常会关闭。 |
| `INFO`   | **信息级别日志**：提供普通操作的运行状态信息。               | - 输出正常运行状态的信息，如程序启动、配置加载成功等。       |
| `WARN`   | **警告级别日志**：表示潜在问题或非正常情况，但程序还能继续运行。 | - 未预期但可恢复的问题，例如配置文件中缺少非关键字段。 - 网络连接临时中断但会自动重连。 |
| `ERROR`  | **错误级别日志**：表示严重问题，但不会立即导致程序终止。     | - 无法访问外部服务。 - 操作失败导致某个功能不可用。          |
| `FATAL`  | **致命错误日志**：表示不可恢复的错误，程序通常会在记录后终止。 | - 无法连接到数据库。 - 必需的资源文件缺失。                  |



## errgroup

`errgroup` 是 Go 的一个标准库扩展包，位于 `golang.org/x/sync/errgroup` 中，主要用于处理一组并发任务并收集其中的第一个错误。它简化了并发任务的管理，尤其是当多个 goroutine 需要同步地处理任务且其中任何一个任务失败时，都需要取消剩余任务的情况。

### 核心特性

- **错误传播**: `errgroup` 能够捕获并返回并发任务中最早发生的错误。
- **任务取消**: 任务一旦出现错误，`errgroup` 会自动取消其余未完成的任务。
- **同步等待**: 使用 `errgroup.Wait()` 方法等待所有任务完成或某个任务失败。

### 安装

`errgroup` 属于 Go 扩展包，首先需要安装：

```
go get golang.org/x/sync/errgroup
```

### 使用方法

1. **`errgroup.Group`**: 主要的结构体，负责管理多个并发任务。
2. **`Go(func() error)`**: 启动一个 goroutine，函数返回一个错误，如果出现错误将被 `errgroup` 捕获。
3. **`Wait()`**: 等待所有 goroutine 完成，如果任何一个 goroutine 返回错误，`Wait()` 会返回该错误。

### 示例

以下是一个简单的示例，展示如何使用 `errgroup` 来管理多个并发任务，并在其中某个任务失败时取消所有任务：

```
package main

import (
    "context"
    "errors"
    "fmt"
    "golang.org/x/sync/errgroup"
    "time"
)

func main() {
    // 创建 context 和 errgroup
    ctx := context.Background()
    g, ctx := errgroup.WithContext(ctx)

    // 启动多个并发任务
    g.Go(func() error {
        time.Sleep(2 * time.Second)
        return nil // 成功
    })

    g.Go(func() error {
        time.Sleep(1 * time.Second)
        return errors.New("task 2 failed") // 返回错误
    })

    g.Go(func() error {
        <-ctx.Done() // 检查 context 是否被取消
        return ctx.Err()
    })

    // 等待所有任务完成
    if err := g.Wait(); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("All tasks completed successfully")
    }
}
```

### 解读

- `errgroup.WithContext(ctx)`：创建了一个 `errgroup.Group` 和关联的上下文 `ctx`，如果其中一个任务失败，其他正在执行的任务会通过上下文感知到取消信号。
- `g.Go(func() error)`：启动多个 goroutine 并处理错误。
- `ctx.Done()`：如果某个 goroutine 返回了错误，其余任务会在 `ctx.Done()` 中收到取消信号并停止执行。
- `g.Wait()`：等待所有任务完成，并返回第一个遇到的错误。

### 典型应用场景

- **批量请求**: 处理多个 HTTP 请求或数据库查询，并在任何一个请求失败时终止所有请求。
- **并发任务控制**: 需要处理多个并发任务，并确保所有任务成功完成或者某个任务失败后可以立刻停止其他任务。

### 优势

- **简洁性**: 简化了 goroutine 的错误处理和同步控制。
- **自动取消**: 任务失败后自动取消其他任务，避免不必要的资源消耗。

## testing

### 1. 定义一个简单的模型

假设我们有一个 `User` 模型，并且有一些方法对该模型进行操作。

```
// user.go
package model

type User struct {
    Name  string
    Age   int
    Email string
}

// IsAdult 方法判断用户是否成年
func (u *User) IsAdult() bool {
    return u.Age >= 18
}

// UpdateEmail 更新用户的邮箱
func (u *User) UpdateEmail(newEmail string) {
    u.Email = newEmail
}
```

### 2. 编写测试代码

接下来，我们为这个 `User` 模型编写单元测试。测试文件通常以 `_test.go` 结尾，以下是测试代码：

```
// user_test.go
package model

import (
    "testing"
)

// 测试 IsAdult 方法
func TestIsAdult(t *testing.T) {
    user := User{Name: "Alice", Age: 20, Email: "alice@example.com"}
    
    if !user.IsAdult() {
        t.Errorf("Expected %s to be an adult, but got not adult", user.Name)
    }

    user2 := User{Name: "Bob", Age: 16, Email: "bob@example.com"}
    
    if user2.IsAdult() {
        t.Errorf("Expected %s to be a child, but got adult", user2.Name)
    }
}

// 测试 UpdateEmail 方法
func TestUpdateEmail(t *testing.T) {
    user := User{Name: "Alice", Age: 20, Email: "alice@example.com"}
    newEmail := "newalice@example.com"
    
    user.UpdateEmail(newEmail)
    
    if user.Email != newEmail {
        t.Errorf("Expected email to be updated to %s, but got %s", newEmail, user.Email)
    }
}
```

### 3. 运行测试

使用 `go test` 来运行测试。

```
go test -v
```

## 压力测试

Go 的 `testing` 包提供了基准测试功能，可以用于测量函数在特定条件下的性能表现。尽管基准测试主要用于性能评估，但也可以用于简单的压力测试场景。

Go 语言推荐测试文件和源代码文件放在一块，测试文件以 `_test.go` 结尾。比如，当前 package 有 `calc.go` 一个文件，我们想测试 `calc.go` 中的 `Add` 和 `Mul` 函数，那么应该新建 `calc_test.go` 作为测试文件。

- 测试用例名称一般命名为 `Test` 加上待测试的方法名。
- 测试用的参数有且只有一个，在这里是 `t *testing.T`。
- 基准测试(benchmark)的参数是 `*testing.B`，TestMain 的参数是 `*testing.M` 类型。
- `go test -v`，`-v` 参数会显示每个用例的测试结果，另外 `-cover` 参数可以查看覆盖率。

### 示例：基准测试一个处理函数

假设我们有一个简单的 HTTP 服务器，提供一个处理请求的函数 `HandleRequest`。我们希望通过基准测试来评估其在高并发情况下的性能。

#### 1. 定义处理函数

```
// server.go
package main

import (
    "fmt"
    "net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // 模拟处理逻辑
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", HandleRequest)
    http.ListenAndServe(":8080", nil)
}
```

#### 2. 编写基准测试

```
// server_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func BenchmarkHandleRequest(b *testing.B) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        b.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HandleRequest)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        handler.ServeHTTP(rr, req)
    }
}
```

#### 3. 运行基准测试

在终端中运行以下命令：

```
go test -bench=.
```

#### 4. 解释结果

基准测试会输出类似以下的结果：

```
goos: darwin
goarch: amd64
BenchmarkHandleRequest-8   	10000000	       200 ns/op
PASS
ok  	path/to/your/package	2.345s
```

- **BenchmarkHandleRequest-8**：表示在 8 个 CPU 核心上运行的基准测试。
- **10000000**：表示函数被调用的次数。
- **200 ns/op**：每次操作的平均耗时。