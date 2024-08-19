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















