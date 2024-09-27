# Slices

路径：go-master\src\slices\slices.go

git地址：https://github.com/golang/go

## 创建切片

src\builtin\builtin.go中的make函数，调用src\runtime\malloc.go中的mallocgc函数分配内存，大对象（>32kb）分配到heap上,小对象分配到P 上。

### 具体的mallogc函数：

```
mp := acquirem() //获取当前g的m
mp.mallocing = 1 //检查没bug将mallocign置为1
c := getMCache(mp) //创建一个mini的catch
//如果不需要扫描内存，且分配的小于Tinysize
///判断是否是8，4，2的倍数，用于内存对齐
////如果不是8的倍数
return (n + a - 1) &^ (a - 1)//用于对齐数值到某个对齐边界。具体来说，这段代码通常用于将一个数 n 向上对齐到最接近的 a 的倍数。如果10，分配16个内存
v := nextFreeFast(span)//返回下一个空闲内存return gclinkptr(uintptr(result)*s.elemsize + s.base())预分配的
//如果没有
v, span, shouldhelpgc = c.nextFree(tinySpanClass)//返回系统申请的
//分配的内存较大到堆上
span = c.allocLarge(size, noscan)//分配大内存
///计算分配不一致，内部状态
///更新堆活跃内存量
mheap_.central[spc].mcentral.fullSwept(mheap_.sweepgen).push(s)
s.limit = s.base() + size
s.initHeapBits(false)
//将内存块 s 推入到 mcentral 对象的内存块列表中。这通常是在内存块被回收后，将其重新放入内存池以供未来的分配使用。
//设置内存块 s 的限制地址，使得 s 可以管理的内存范围从 s.base() 开始，到 s.base() + size 结束。
//初始化内存块 s 的堆位标记。
//判断其他的参数
releasem(mp)//释放gp.stackguard0 = stackPreempt
```

底层汇编上runtime/asm_*.s中的分配内存。

**`mheap`**：全局的堆结构体，负责管理堆内存。

**`mcentral`**：用于管理不同大小的内存块，每个 `mcentral` 负责管理一类大小的内存块。

**`mspan`**：管理一段连续的内存区域。`mspan` 是内存分配的最小单位，一个 `mspan` 可以包含多个小对象，也可以用于大对象的分配。

make函数调用mallocgc则通过cmd/compile/internal/gc看到相关实现

```
// arguments in registers.
TEXT callRet<>(SB), NOSPLIT, $40-0
	NO_LOCAL_POINTERS
	MOVQ	DX, 0(SP)
	MOVQ	DI, 8(SP)
	MOVQ	SI, 16(SP)
	MOVQ	CX, 24(SP)
	MOVQ	R12, 32(SP)
	CALL	runtime·reflectcallmove(SB)
	RET
```



## 切片增加

 插入切片时，`func Insert[S ~[]E, E any](s S, i int, v ...E) S {`通过泛型实现，`_ = s[i:] `

### 1. **触发副作用**

虽然这种形式的代码并没有使用返回的结果，但它有时用来触发一些**副作用**。副作用可以包括：

- **延迟计算**：某些操作可能不会立即执行，而是在你访问到特定的切片部分时才会执行。
- **避免编译器优化**：Go 编译器有时会优化掉未使用的代码。使用 `_ = s[i:]` 可以告诉编译器，某段内存（切片）是有用的，防止编译器优化掉它。

### 2. **避免未使用变量的警告**

在 Go 中，如果一个变量被声明但没有被使用，编译器会抛出一个编译错误。为了避免这种错误，你可以通过使用空标识符 `_` 来显式地忽略该变量。

- `_ = s[i:]` 意味着虽然我们从 `i` 开始切片，但我们并不关心返回的结果，避免编译器的未使用变量警告。

如果大于切片的cap，则append。append定义在`src\builtin\builtin.go`中，由于append` 函数属于内置函数，通常直接在 Go 编译器中实现，而不是在 Go 标准库的源码中以普通函数形式定义的。因此不能直接跳转，可以在`src\runtime\slice.go`中的`func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {`中看到具体实现，这是运行时在编译时生成底层代码来处理切片的扩容和元素追加，它的职责是创建一个更大的底层数组，并将旧的元素复制过去。runtime中定义了slice结构体，包含

```
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```

。growslice中，如果可竞争，先获取`getcallerpc()`，不同架构实现位置不同，如amd64在`src\runtime\asm_amd64.s`用于获取调用者的程序计数器（Program Counter，PC）。

```
TEXT _rt0_amd64(SB),NOSPLIT,$-8
	MOVQ	0(SP), DI	// argc
	LEAQ	8(SP), SI	// argv
	JMP	runtime·rt0_go(SB)
```

代码细节

#### 1. `TEXT _rt0_amd64(SB),NOSPLIT,$-8`

- `TEXT _rt0_amd64(SB),NOSPLIT,$-8`

  :

  - **`TEXT`**：这是 Go 汇编语言中的伪指令，用于定义一个新的函数或代码块。
  - **`_rt0_amd64(SB)`**：这是函数的名称，`_rt0_amd64`。`SB` 是符号表的标识符，表示函数 `_rt0_amd64` 在符号表中的位置。
  - **`NOSPLIT`**：这个标志表示函数不进行栈分割。栈分割用于处理函数调用时的栈空间管理，但在启动代码中通常不需要，因为我们在这个阶段还没有进入正常的 Go 运行时环境。
  - **`$-8`**：这个参数指定了栈空间的大小，这里是 `-8` 字节，表示需要为函数分配 8 字节的栈空间。

#### 2. `MOVQ 0(SP), DI`

- `MOVQ 0(SP), DI`

  :

  - **`MOVQ`**：这是一个将数据从一个位置移动到另一个位置的汇编指令。`Q` 表示操作的是 64 位数据（即 `quad`）。
  - **`0(SP)`**：表示从栈顶 `SP` 开始的 0 偏移量处读取数据。栈顶存储了程序的启动参数。
  - **`DI`**：这是一个寄存器，用于存储 `argc`，即命令行参数的数量。这个指令将栈顶的 `argc` 值移动到 `DI` 寄存器中。

#### 3. `LEAQ 8(SP), SI`

- `LEAQ 8(SP), SI`

  :

  - **`LEAQ`**：这是加载有效地址的指令，用于将计算出的地址存储到寄存器中。
  - **`8(SP)`**：表示从栈顶 `SP` 开始的 8 偏移量处的地址。这个偏移量指向命令行参数数组 `argv` 的起始位置。
  - **`SI`**：这是一个寄存器，用于存储 `argv`，即命令行参数的地址。这个指令将 `argv` 的地址加载到 `SI` 寄存器中。

#### 4. `JMP runtime·rt0_go(SB)`

- `JMP runtime·rt0_go(SB)`

  :

  - **`JMP`**：这是一个无条件跳转指令。
  - **`runtime·rt0_go(SB)`**：这是 Go 运行时的一个函数，`rt0_go`。这个函数是 Go 程序启动的核心函数，它会完成进一步的初始化工作。
  - 通过跳转到 `runtime·rt0_go`，代码将控制权转移给 Go 运行时的初始化函数，完成程序的启动过程。

然后调用`func racereadrangepc(p uintptr, n uintptr, pc uintptr, fn abi.FuncPCABI)` 用于在数据竞争检测过程中读取数据。**`p`**: 内存区域的起始地址。

**`n`**: 内存区域的大小（以字节为单位）。

**`pc`**: 调用此函数的程序计数器（指向函数的指针）。

**`fn`**: 代表被调用的函数的 ABI（Application Binary Interface）标识符。进行一系列检测，确认没问题调用`p = mallocgc(capmem, et, true)`，**`size uintptr`**:

- 这是要分配的内存的大小（以字节为单位）。`capmem` 在你的代码中是表示需要分配的内存的大小的参数。

**`typ \*_type`**:

- 这是一个指向类型信息的指针，表示分配的内存块的类型信息。`et` 在你的代码中是表示类型的变量，它指向一个 `_type` 结构体，这个结构体包含了有关类型的信息，比如元素的大小和对齐方式。

**`needszero bool`**:

- 这是一个布尔值，指示是否需要将分配的内存初始化为零。`true` 表示内存需要被初始化为零，`false` 则表示不需要。通常，在 Go 中，如果分配的内存是某种结构体或切片的底层数组，`needszero` 为 `true` 可以确保内存被清零，这对于防止使用未初始化的数据是很重要的。

`src\runtime\malloc.go`获取分配内存块的指针。

```
type notInHeapSlice struct {
    array *notInHeap
    len   int
    cap   int
}
```

`notInHeapSlice` 结构体在 Go 运行时系统中可能用于以下情况：

1. **非堆内存切片**：
   - `notInHeapSlice` 结构体可能表示一个不在堆内存中分配的切片。这可能是为了优化内存管理，或者用于某些特定的运行时操作。
   - 在 Go 的运行时，内存管理的细节可能会有所不同，尤其是在涉及到垃圾回收和内存分配的内部实现中。
2. **特殊内存区域**：
   - `notInHeap` 可能表示一种特殊的内存区域，可能是栈内存、全局内存或其他非堆内存区域。
   - 这种结构可能用于在这些特殊内存区域中管理切片，而不是使用标准的堆内存分配。
3. **性能优化**：
   - 在 Go 运行时，特别是在高性能或低延迟的场景下，使用非堆内存切片可以减少垃圾回收的开销或避免堆内存的额外开销。