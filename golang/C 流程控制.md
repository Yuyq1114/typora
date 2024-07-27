### 1. 条件语句

条件语句用于根据布尔表达式的结果来选择执行不同的代码块。

#### `if` 语句

基本形式：

```
if condition {
    // 如果 condition 为 true，则执行这里的代码
}
```

带 `else` 的 `if` 语句：

```
if condition {
    // 如果 condition 为 true，则执行这里的代码
} else {
    // 如果 condition 为 false，则执行这里的代码
}
```

带 `else if` 的 `if` 语句：

```
if condition1 {
    // 如果 condition1 为 true，则执行这里的代码
} else if condition2 {
    // 如果 condition2 为 true，则执行这里的代码
} else {
    // 如果以上条件都为 false，则执行这里的代码
}
```

条件初始化语句：

```
if x := computeValue(); x < 0 {
    fmt.Println("x is negative")
} else if x == 0 {
    fmt.Println("x is zero")
} else {
    fmt.Println("x is positive")
}
```

#### `switch` 语句

`switch` 语句用于选择多个分支之一执行，通常比多个 `if-else` 更清晰。

基本形式：

```
switch expression {
case value1:
    // 如果 expression == value1，则执行这里的代码
case value2:
    // 如果 expression == value2，则执行这里的代码
default:
    // 如果没有匹配到任何 case，则执行这里的代码
}
```

不带表达式的 `switch`：

```
switch {
case condition1:
    // 如果 condition1 为 true，则执行这里的代码
case condition2:
    // 如果 condition2 为 true，则执行这里的代码
default:
    // 如果没有匹配到任何条件，则执行这里的代码
}
```

### 2. 循环语句

循环语句用于反复执行某段代码。

#### `for` 循环

`for` 循环是Go语言唯一的循环结构，但它可以用于实现多种循环模式。

基本形式：

```
for initialization; condition; post {
    // 循环体
}
```

示例：

```
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

省略初始化和后处理：

```
i := 0
for ; i < 10; {
    fmt.Println(i)
    i++
}
```

省略条件（无限循环）：

```
for {
    fmt.Println("Infinite loop")
    break  // 避免实际无限循环
}
```

#### `range` 关键字

`range` 用于迭代数组、切片、映射、字符串或通道。

迭代数组或切片：

```
numbers := []int{1, 2, 3, 4, 5}
for index, value := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}
```

迭代映射：

```
person := map[string]string{"Name": "Alice", "Age": "30"}
for key, value := range person {
    fmt.Printf("%s: %s\n", key, value)
}
```

迭代字符串：

```
str := "hello"
for index, runeValue := range str {
    fmt.Printf("Index: %d, Rune: %c\n", index, runeValue)
}
```

### 3. 跳转语句

跳转语句用于在循环或其他控制结构中改变程序的执行顺序。

#### `break` 语句

`break` 语句用于终止最近的 `for` 循环或 `switch` 语句。

```
for i := 0; i < 10; i++ {
    if i == 5 {
        break
    }
    fmt.Println(i)
}
```

#### `continue` 语句

`continue` 语句用于跳过当前循环的剩余语句，并继续下一次循环。

```
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue
    }
    fmt.Println(i)
}
```

#### `goto` 语句

`goto` 语句用于无条件跳转到被标记的代码位置。尽量避免使用 `goto`，以免影响代码的可读性。

```C
func main() {
    i := 0
Loop:
    fmt.Println(i)
    i++
    if i < 5 {
        goto Loop
    }
}
```