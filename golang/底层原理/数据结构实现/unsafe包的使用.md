# golang中的unsafe包详解

链接：https://www.jianshu.com/p/10b8870a9e8e
来源：简书

import "unsafe"

unsafe包提供了一些跳过go语言类型安全限制的操作。

## 各种类型所占内存大小

1.bit(位): 二进制数中的一个数位，可以是0或者1，是计算机中数据的最小单位。二进制的一个“0”或一个“1”叫一位
2.Byte(字节): 计算机中数据的基本单位，每8位组成一个字节

int8:   8位，也就是占用一个字节
int16: 16位，也就是2个字节
int32: 4个字节
int64: 8个字节

float32: 4个字节
float64: 8 个字节

int 比较特殊，占用多大取决于你的cpu：32位cpu 就是 4个字节，64位 就是 8 个字节

float32: 4个字节
float64: 8个字节

string
英文的ascii的string 1个英文字符或英文标点占1byte
中文的string 一个文字占用 3 byte

bool: 占用1byte


## 一、unsafe 作用

从golang的定义来看，unsafe 是涉及类型安全的操作。顾名思义，unsafe包提供了一些跳过go语言类型安全限制的操作，它应该被非常谨慎地使用。unsafe可能很危险，但也可能非常有用。例如，当使用系统调用和Go结构必须具有与C结构相同的内存布局时，这时候可能别无选择，只能使用unsafe。

type ArbitraryType int

ArbitraryType在本文档里表示任意一种类型，但并非一个实际存在于unsafe包的类型。

type Pointer *ArbitraryType
Pointer类型用于表示任意类型的指针。


关于指针操作，有4个特殊的只能用于Pointer类型的操作，在unsafe包官方定义里有如下四个描述：
```
1) 任意类型的指针可以转换为一个Pointer类型值
2) 一个Pointer类型值可以转换为任意类型的指针
3) 一个uintptr类型值可以转换为一个Pointer类型值
4) 一个Pointer类型值可以转换为一个uintptr类型值
```
额外在加上一个规则：指向不同类型数据的指针，是无法直接相互转换的，必须借助unsafe.Pointer(类似于C的 void指针)代理一下再转换也就是利用上述的1，2规则。

举例：
```go
func Float64bits(f float64) uint64 {
  // 无法直接转换，报错：Connot convert expression of type *float64 to type *uint64
  // return *(*uint64)(&f)   

 // 先把*float64 转成 Pointer(描述1)，再把Pointer转成*uint64(描述2)
  return *(*uint64)(unsafe.Pointer(&f)) 
}
```

## 二、unsafe的定义：

整体代码比较简单，2个类型定义和3个uintptr的返回函数
```go
package unsafe

// ArbitraryType仅用于文档目的，实际上并不是unsafe包的一部分，它表示任意Go表达式的类型。
type ArbitraryType int

// 任意类型的指针，类似于C的*void
type Pointer *ArbitraryType

// Sizeof返回类型v本身数据所占用的字节数。返回值是“顶层”的数据占有的字节数。例如，若v是一个切片，它会返回该切片描述符的大小，而非该切片底层引用的内存的大小。
func Sizeof(x ArbitraryType) uintptr            // 结构x在内存中占用的确切大小

// 返回结构体中某个field的偏移量
// Offsetof返回类型x所代表的结构体中的字段在结构体中的偏移量，它必须为结构体类型x的字段f的形式。换句话说，它返回该结构体x起始处与该字段f起始处之间的字节数。
func Offsetof(x ArbitraryType) uintptr

// Alignof返回类型x的对齐方式（即类型x在内存中占用的字节数）；若是结构体类型的字段的形式，它会返回字段在该结构体中的对齐方式。
func Alignof(x ArbitraryType) uintptr           // 返回结构体中某个field的对其值（字节对齐的原因）
```

看一个栗子：
```go
package main

import (
    "fmt"
    "unsafe"
)

type Human struct {
    sex  bool           // bool 占1个字节
    age  uint8          // 占用8位，即1个字节
    min  int            // int取决于操作系统，32位系统占4个字节，64位占8个字节
    name string         // 英文字符，一个字符占1个字节；中文，一个字符占3个字节
}

func main() {
    h := Human{
        true,
        30,
        1,
        "hello",
    }
    i := unsafe.Sizeof(h)
    j := unsafe.Alignof(h.age)
    k := unsafe.Offsetof(h.name)
    fmt.Println(i, j, k)
    fmt.Printf("%p\n", &h)
    var p unsafe.Pointer
    p = unsafe.Pointer(&h)
    fmt.Println(p)
}
```

//输出
//32 1 16
//0xc00000a080
//0xc00000a080
// 32：string 占16字节，所以16+16 =32；1 是因为age前是bool，占用1个字节；8是name的偏移是int 占8个字
节

### 三、Pointer使用

前面已经说了，pointer是任意类型的指针，可以指向任意类型数据。参照Float64bits的转换和上述例子的unsafe.Pointer(&h)，所以主要用于转换各种类型

### 四、uintptr

在golang中uintptr的定义是：
```go
type uintptr uintptr 
```
uintptr是golang的内置类型，是一个能存储指针的无符号整型值。

根据描述3，一个unsafe.Pointer指针也可以被转化为uintptr类型，然后保存到指针型数值变量中（注：这只是和当前指针相同的一个数字值，并不是一个指针），然后用以做必要的指针数值运算。（uintptr是一个无符号的整型数，足以保存一个地址）

这种转换虽然也是可逆的，但是将uintptr转为unsafe.Pointer指针可能会破坏类型系统，因为并不是所有的数字都是有效的内存地址。
许多将unsafe.Pointer指针转为uintptr，然后再转回为unsafe.Pointer类型指针的操作也是不安全的。比如下面的例子需要将变量x的地址加上b字段地址偏移量转化为*int16类型指针，然后通过该指针更新x.b：
```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {

    var x struct {
        a bool
        b int16
        c []int
    }

    /**
    unsafe.Offsetof 函数的参数必须是一个字段 x.f, 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞.
    */

    /**
    uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
    指针的运算
    */
    // 和 pb := &x.b 等价
    pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
    *pb = 42
    fmt.Println(x.b) // "42"
}
```
上面的写法尽管很繁琐，但在这里并不是一件坏事，因为这些功能应该很谨慎地使用。不要试图引入一个uintptr类型的临时变量，因为它可能会破坏代码的安全性（注：这是真正可以体会unsafe包为何不安全的例子）。

下面段代码是错误的：
```go
// NOTE: subtly incorrect!
tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
```
产生错误的原因很微妙。

有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动GC。当一个变量被移动，所有的保存发生改变的旧地址的指针必须同时被更新为变量移动后的新地址。从垃圾收集器的视角来看，一个unsafe.Pointer是一个指向变量的指针，因此当变量被移动时对应的指针也必须被更新；但是uintptr类型的临时变量只是一个普通的数字，所以其值不应该被改变。上面错误的代码，由于引入一个非指针的临时变量tmp，导致垃圾收集器无法正确识别这是一个指向变量x的指针。当第二个语句执行时，变量x可能已经被转移，这时候临时变量tmp也就不再是现在的&x.b地址。第三个语句指向向之前无效地址空间的赋值语句将彻底摧毁整个程序！



