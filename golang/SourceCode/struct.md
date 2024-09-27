# 结构体类型的定义

**`rtype` 结构体**:

- Go 使用 `rtype` 结构体来表示所有类型的元数据，包括结构体类型。
- `rtype` 包含类型的基本信息，如大小、对齐、种类（kind）、方法集等。
- 结构体的具体信息存储在 `rtype` 的字段中，包括结构体字段的布局、对齐要求和方法集。

具体的：src\runtime\type.go

```
type rtype struct {
	*abi.Type // embedding is okay here (unlike reflect) because none of this is public
}
```

其中的Type定义为：src\cmd\compile\internal\types\type.go

```
type Type struct {
    //可以用来存储与类型相关的其他信息。具体用途可能取决于编译器的不同实现。
	extra interface{}

	// width is the width of this Type in bytes.表示该类型的宽度，以字节为单位。width 只有在对齐要求（align）大于 0 的情况下才有效。
	width int64 // valid if Align > 0

	// list of base methods (excluding embedding)包含了该类型的所有方法（不包括嵌入的类型中的方法）。
	methods fields
	// list of all methods (including embedding)包含了该类型的所有方法，包括嵌入类型中的方法。
	allMethods fields

	// canonical OTYPE node for a named type (should be an ir.Name node with same sym)对于命名类型（即用户定义的类型），obj 是一个 Object 类型的值，表示一个 IR（中间表示）节点，如 ir.Name 节点，这个节点与该类型的符号（symbol）相关联。
	obj Object
	// the underlying type (type literal or predeclared type) for a defined type对于定义类型（例如，用户定义的结构体或接口），underlying 指向该类型的基础类型。对于内置类型或字面类型，underlying 指向类型本身。
	underlying *Type

	// Cache of composite types, with this type being the element type.这个缓存用于快速查找与该类型相关的复合类型，如指针和切片。

	cache struct {
		ptr   *Type // *T, or nil指向该类型的指针类型。如果该类型是 *T，那么 ptr 指向 T。
		slice *Type // []T, or nil指向该类型的切片类型。如果该类型是 []T，那么 slice 指向 T。
	}

	kind  Kind  // kind of typeKind 是一个枚举类型，描述了类型的具体类别，例如结构体、数组、切片、接口等。
	align uint8 // the required alignment of this type, in bytes (0 means Width and Align have not yet been computed)表示该类型的对齐要求，以字节为单位。如果 align 为 0，则表示 width 和 align 尚未计算。

	intRegs, floatRegs uint8 // registers needed for ABIInternal表示该类型在内部 ABI（应用程序二进制接口）中需要的整数寄存器和浮点寄存器的数量。这用于决定如何在函数调用中传递参数。

	flags bitset8//这是一个位集，用于存储与类型相关的各种标志。具体标志取决于编译器的实现和类型的特性。
	alg   AlgKind // valid if Align > 0表示算法的种类。如果 align 大于 0，则 alg 记录了与该类型相关的算法（如哈希算法、比较算法等）。

	// size of prefix of object that contains all pointers. valid if Align > 0.
	// Note that for pointers, this is always PtrSize even if the element type
	// is NotInHeap. See size.go:PtrDataSize for details.表示对象前缀中包含所有指针的字节数。对于指针类型，这通常是 PtrSize，即指针的大小。这个字段在对齐要求（align）大于 0 时有效。
	ptrBytes int64
}
```

## 循环使用：

在 Go 语言的实现中，循环引用问题通常通过以下方法解决：

- **延迟初始化**:
  - Go 编译器在处理类型信息时会使用延迟初始化。具体来说，当编译器需要创建 `rtype` 的实例时，类型的完整定义可能尚未完成。为了处理这种情况，Go 语言的运行时和编译器会采用延迟填充或初始化的策略。
  - 在初始化类型时，编译器会先创建一个部分填充的 `rtype` 结构体，待类型的所有部分完成后，再填充完整。
- **类型指针**:
  - `type` 结构体中的 `ptrToThis` 字段是一个指向 `rtype` 自身的指针。这个字段在类型完全定义之后被填充，避免了直接循环引用的问题。
  - 通过指针引用，Go 可以在运行时动态地更新和访问类型信息，而不是在编译时就需要完成所有的引用。

# 结构体变量与单独定义的区别

**结构体**：字段在内存中是连续的，编译器可能会添加填充字节来满足对齐要求，字段间访问更高效，特别是在频繁访问多个字段的场景下缓存友好。

**独立变量**：没有连续性要求，可能占用更少的内存（没有填充字节），但在内存中的布局由编译器决定，访问效率可能不如结构体高。