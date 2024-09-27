string底层是一个byte切片，byte切片是uint8的别名，由于字符串不可变，因此并发安全。

**栈上存储**: 栈上存储的数据通常包括局部变量和函数参数。这些数据的生命周期与函数调用的栈帧相关联。栈上的存储是临时的，数据在函数返回时会被销毁。

典型的栈上存储的数据类型包括：

- 基本数据类型（如 `int`, `float64`, `bool`, `byte` 等）
- 小型结构体（如果结构体的大小比较小且不涉及复杂的内存分配）

**堆上存储**: 堆上存储的数据的生命周期与对象的引用相关联。对象在堆上分配内存，直到没有更多的引用指向该对象为止，垃圾回收器会回收这些内存。

典型的堆上存储的数据类型包括：

- 大型结构体（尤其是当它们被传递给函数或方法时）
- 动态分配的内存（如通过 `new` 或 `make` 分配的内存）
- 字符串数据（虽然字符串的引用在栈上，但实际的字符串数据存储在堆上）
- 切片（切片本身在栈上，但切片的底层数组通常在堆上分配）

# go运行时string的concat

具体位置src\runtime\string.go

```
func concatstrings(buf *tmpBuf, a []string) string {//将多个字符串拼接成一个新的字符串。
    idx := 0
    l := 0
    count := 0
    for i, x := range a {//计算总长度和非空字符串数量
       n := len(x)
       if n == 0 {
          continue
       }
       if l+n < l {
          throw("string concatenation too long")
       }
       l += n
       count++
       idx = i
    }
    if count == 0 {
       return ""
    }

    // If there is just one string and either it is not on the stack
    // or our result does not escape the calling frame (buf != nil),
    // then we can return that string directly.
    //如果只有一个非空字符串，并且它要么在临时缓冲区中（buf != nil），要么不在栈上（!stringDataOnStack(a[idx])），则直接返回该字符串。这样可以避免不必要的内存分配和复制。
    if count == 1 && (buf != nil || !stringDataOnStack(a[idx])) {
       return a[idx]
    }
    s, b := rawstringtmp(buf, l)//调用 rawstringtmp 函数分配一个长度为 l 的临时内存，用于存储拼接后的结果
    for _, x := range a {//遍历字符串切片 a，将每个字符串复制到缓冲区 b 中。每次复制后，更新 b 的位置。
       copy(b, x)
       b = b[len(x):]
    }
    return s
}
```

## 具体的rawstringtmp格式如下：

```
func rawstringtmp(buf *tmpBuf, l int) (s string, b []byte) {
    if buf != nil && l <= len(buf) {//检查 buf 是否不为 nil，并且缓冲区的大小是否足够容纳长度为 l 的字节切片。
       b = buf[:l]
       s = slicebytetostringtmp(&b[0], len(b))//使用 slicebytetostringtmp 函数将字节切片 b 转换为字符串 s
    } else {
       s, b = rawstring(l)//调用 rawstring(l) 函数分配新的内存。rawstring 函数会返回一个新的字符串和字节切片，
    }
    return
}
```

### 具体的rawstring结构如下：

```
func rawstring(size int) (s string, b []byte) {
    p := mallocgc(uintptr(size), nil, false)//调用smallocgc分配内存
    return unsafe.String((*byte)(p), size), unsafe.Slice((*byte)(p), size)//返回不安全string和不安全切片
}
```



# strings库中的string分割

src\strings\strings.go

```
func genSplit(s, sep string, sepSave, n int) []string {//将字符串 s 按照分隔符 sep 切分成若干部分，并返回一个字符串切片。
    if n == 0 {
       return nil//
    }
    if sep == "" {
       return explode(s, n)
    }
    if n < 0 {
       n = Count(s, sep) + 1
    }

    if n > len(s)+1 {
       n = len(s) + 1
    }
    a := make([]string, n)//创建a用于存储返回的string
    n--
    i := 0
    //在循环中，查找分隔符 sep 在字符串 s 中的位置 m。如果没有找到分隔符（即 m < 0），退出循环。
    for i < n {
       m := Index(s, sep)//
       if m < 0 {
          break
       }
       a[i] = s[:m+sepSave]//将字符串 s 的前 m+sepSave 个字符（包括分隔符，如果 sepSave 为 1）赋值给 a[i]。
       s = s[m+len(sep):]//更新 s 为分隔符之后的部分。
       i++
    }
    a[i] = s
    return a[:i+1]
}
```