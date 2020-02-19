


# 2.协程Goroutines

## 2.1 goroutine占用内存资源

golang中，提升并发题型的主要是 goroutines。goroutine可以理解成轻量级的线程，其使用成本很低，比线程要小得多，可以认为它们几乎是没有成本的。

Go 运行时runtime是为运行数以万计的 goroutines 所设计的，即使可以同时开启上十万goroutines。

但是，每个 goroutine 确实消耗了 goroutine 栈的内存量，目前至少为 2k。如果一个goroutine按照2KB算，1000000个goroutine就是2GB：
2048byte * 1,000,000 goroutines == 2GB 内存，什么都不干的情况下。

这也许算多，可能也不算多，这取决于机器上其他耗费内存的应用。

## 2.2 goroutine 什么时候退出

虽然 goroutine 的启动和运行成本都很低，它们的内存占用也是有限的，但不可能创建无限数量的 goroutine。

每次在程序中使用go关键字启动 goroutine 时，需要知道这个 goroutine 何时退出，如何退出。如果使用不当，就可能造成潜在的内存泄漏。

实现此目的的一个好方法是利用如waitgroup, run.Group（https://github.com/oklog/run）， workgroup.Group (https://github.com/heptio/workgroup)这类的东西。

Peter Bourgon has a great presentation on the design behing run.Group from GopherCon EU


Go 对一些请求使用高效的网络轮询
Go 运行时使用高效的操作系统轮询机制（kqueue，epoll，windows IOCP等）处理网络IO。 许多等待的 goroutine 将由一个操作系统线程提供服务。

但是，对于本地文件IO（channel 除外），Go 不实现任何 IO 轮询。每一个*os.File在运行时都消耗一个操作系统线程。

大量使用本地文件IO会导致程序产生数百或数千个线程；这可能会超过操作系统的最大值限制。

您的磁盘子系统可能处理不数百或数千个并发IO请求。

注意程序中的 IO 复杂度
如果你写的是服务端程序，那么其主要工作是复用网络连接客户端和存储在应用程序中的数据。

大多数服务端程序都是接受请求，进行一些处理，然后返回结果。这听起来很简单，但有的时候，这样做会让客户端在服务器上消耗大量（可能无限制）的资源。下面有一些注意事项：

每个请求的IO操作数量；单个客户端请求生成多少个IO事件？ 如果使用缓存，则它可能平均为1，或者可能小于1。
服务查询所需的读取量；它是固定的？N + 1的？还是线性的（读取整个表格以生成结果的最后一页）？
如果内存都不算快，那么相对来说，IO操作就太慢了，你应该不惜一切代价避免这样做。 最重要的是避免在请求的上下文中执行IO——不要让用户等待磁盘子系统写入磁盘，甚至连读取都不要做。

使用流式 IO 接口
尽可能避免将数据读入[]byte 并传递使用它。

根据请求的不同，你最终可能会将兆字节（或更多）的数据读入内存。这会给GC带来巨大的压力，并且会增加应用程序的平均延迟。

作为替代，最好使用io.Reader和io.Writer构建数据处理流，以限制每个请求使用的内存量。

如果你使用了大量的io.Copy，那么为了提高效率，请考虑实现io.ReaderFrom / io.WriterTo。 这些接口效率更高，并避免将内存复制到临时缓冲区。

超时，超时，还是超时
永远不要在不知道需要多长时间才能完成的情况下执行 IO 操作。

你要在使用SetDeadline，SetReadDeadline，SetWriteDeadline进行的每个网络请求上设置超时。

您要限制所使用的阻塞IO的数量。 使用 goroutine 池或带缓冲的 channel 作为信号量。

var semaphore = make(chan struct{}, 10)

func processRequest(work *Work) {
        semaphore <- struct{}{} // 持有信号量
        // 执行请求
        <-semaphore // 释放信号量
}



# 1 内存优化

GO性能优化小结：https://www.cnblogs.com/zhangboyu/p/7456609.html

## 1.1 减少空间分配

## 1.1.1 小对象合并成结构体一次分配，减少内存分配次数

在C/C++程序中，小对象在堆上频繁地申请、释放，会造成内存碎片（有的叫空洞），导致分配大的对象时无法申请到连续的内存空间。为了解决这种情况，一般采用内存池的方式。Go runtime底层也采用内存池，但每个span大小为4k，同时维护一个cache。cache有一个0到n的list数组，list数组的每个单元挂载的是一个链表，链表的每个节点就是一块可用的内存，同一链表中的所有节点内存块都是大小相等的；但是不同链表的内存大小是不等的，也就是说list数组的一个单元存储的是一类固定大小的内存块，不同单元里存储的内存块大小是不等的。这就说明cache缓存的是不同类大小的内存对象，当然想申请的内存大小最接近于哪类缓存内存块时，就分配哪类内存块。当cache不够再向spanalloc中分配。

建议：小对象合并成结构体一次分配，示意如下：
```go
for k, v := range m {
    k, v := k, v // copy for capturing by the goroutine
    go func() {
        // using k & v
    }()
}
替换为：

for k, v := range m {
    x := struct {k , v string} {k, v} // copy for capturing by the goroutine
    go func() {
        // using x.k & x.v
    }()
}
``` 

### 1.1.2 接口方法API中减少给调用方增加垃圾

考虑这两个 Read 方法

func (r *Reader) Read() ([]byte, error)
func (r *Reader) Read(buf []byte) (int, error)
第一个 Read 方法不带参数，并将一些数据作为[]byte返回。 第二个采用[]byte缓冲区并返回读取的字节数。

第一个 Read 方法总是会分配一个缓冲区，这会给 GC 带来压力。 第二个填充传入的缓冲区。

## 1.2 缓存区内容一次分配足够大小空间，并适当复用

在协议编解码时，需要频繁地操作[]byte，可以使用bytes.Buffer或其它byte缓存区对象。

建议：bytes.Buffert等通过预先分配足够大的内存，避免当添加元素超过指定长度后时（扩容，进行）动态申请内存，这样可以减少内存分配次数。同时对于byte缓存区对象考虑适当地复用。

对于切片，虽然使用append方法添加元素很方便，但是有代价。

切片的增长在元素到达 1024 个之前一直是两倍地变化，在到达 1024 个之后，大约是 25% 地增长。

下面的程序，声明一个长度和容量为1024的切片，然后append一个元素，可以看到append之后的容量并没有增加一倍，而是25%左右。
```go
func main() {
        b := make([]int, 1024)
        fmt.Println("len:", len(b), "cap:", cap(b))
        b = append(b, 99)
        fmt.Println("len:", len(b), "cap:", cap(b))
}
output:
len: 1024 cap: 1024
len: 1025 cap: 1280
```
如果使用 append，有可能超过预分配的长度，造成成倍地扩容空间，复制大量数据并产生大量垃圾。

如果事先知道切片的长度，最好预先分配大小以避免复制，并确保目标的大小完全正确。
```go
var s []string                  //这种没有给定切片长度的声明方式，可能造成空间的扩容，大量数据复制并产生垃圾。
s := make([]string, len(vals))  // 预先分配特定长度的空间，就避免了中间垃圾
```


## 1.3 slice和map通过make方式进行创建时，预估大小指定容量

slice和map与数组不一样，不存在固定空间大小，可以根据增加元素来动态扩容。

slice初始会指定一个长度（没有指定的话或者make的时候指定len=0，在添加元素后，其默认长度是切片中元素个数。如果当前切片中元素个数超过2^n到小于1024，则其容量cap会扩大一倍），对slice进行append等操作时，如果容量不够，会自动扩容，扩容机制为：
```markdown
如果新的大小是当前大小2倍以上，则容量增长为新的大小；
否则，循环以下操作：如果当前容量小于1024，按2倍增加；大于1024，每次按当前容量的1/4增长，直到增长的容量超过或等于新大小。
```

map的扩容比较复杂，每次扩容会增加到上次容量的2倍。它的结构体中有一个buckets和oldbuckets，用于实现增量扩容：
```markdown
正常情况下，直接使用buckets，oldbuckets为空；
如果正在扩容，则oldbuckets不为空，buckets是oldbuckets的2倍。
```

建议：初始化时预估大小指定容量。
```go
m := make(map[string]string, 100)
s := make([]string, 0, 100) // 注意：对于slice进行make声明时，make中的第二个参数是初始大小，第三个参数才是容量
```

## 1.4 长调用栈避免申请较多的临时对象
goroutine的调用栈默认大小是4K（1.7修改为2K），它采用连续栈机制，当栈空间不够时，Go runtime会自动扩容：
```markdown
当栈空间不够时，按2倍增加，原有栈中的数据直接copy到新的栈空间，变量指针指向新的空间地址；
退栈会释放栈空间的占用，GC时发现栈空间占用不到1/4时，则栈空间减少一半。
```
比如栈的最终大小2M，则极端情况下，就会有10次的扩栈操作，这会带来性能下降。

建议：

控制调用栈和函数的复杂度，不要在一个goroutine中做完所有的逻辑；
如果的确需要长调用栈，而考虑goroutine池化，避免频繁创建goroutine带来栈空间的变化。

## 1.5 避免频繁创建临时对象
Go在GC时会引发stop the world，即整个情况暂停。虽1.7版本已大幅优化GC性能，1.8甚至在最坏情况下GC为100us。但暂停时间还是取决于临时对象的个数，临时对象数量越多，暂停时间可能越长，并消耗CPU。

建议：GC优化方式是尽可能地减少临时对象的个数：

尽量使用局部变量
将多个局部变量合并一个大的结构体或数组，减少扫描对象的次数，一次回收尽可能多的内存。


## 1.6 切片使用注意事项

### 1.6.1 strings vs []bytes
Go 语言中 string 是不可改变的，而 []byte 是可变的。

通常简短的程序喜欢使用 string，而大多数 IO 操作更喜欢使用 []byte。

在使用过程中，要尽可能避免 []byte 到 string 的转换。对于一个值来说，最好选定一种表示方式，要么是[]byte，要么是string。 通常情况下，如果是从网络或磁盘读取数据，将使用[]byte 表示。

bytes 包也有一些和 strings 包相同的操作函数，如Split， Compare， HasPrefix， Trim等。

实际上， strings 使用和 bytes 包底层使用相同的汇编原语。

### 1.6.2 使用 []byte 当做 map 的 key时的优化

使用 string 作为 map 的 key 是很常见的，但有时拿到的是一个 []byte，编译器为这种情况实现了特定的优化方式：
```go
var bytes []byte
var m map[string]string
v, ok := m[string(bytes)]
```
如上面这样写，编译器会避免将字节切片转换为字符串到 map 中查找，这是非常特定的细节。如果像下面这样写，这个优化就会失效：
```go
key := string(bytes)
val, ok := m[key]
```

优化字符串连接操作
Go 的字符串是不可变的。连接两个字符串就会生成（中间变量）第三个字符串。


# 2 并发优化

## 2.1 高并发的任务处理使用goroutine池

goroutine虽轻量，但对于高并发的轻量任务处理，通过频繁地创建goroutine来执行，执行效率并不会太高效：
```markdown
过多的goroutine创建，会影响go运行时runtime对goroutine的调度，以及GC消耗；
高并发时，若出现调用异常造成阻塞积压，这些大量的goroutine短时间积压可能导致程序崩溃。
```

## 2.2 避免高并发调用同步系统接口
goroutine的实现，是通过同步来模拟异步操作。如下操作不会阻塞go运行时runtime的线程调度：
```markdown
网络IO
锁
channel
time.sleep
基于底层系统异步调用的Syscall
```

下面这些对go运行时的阻塞，会创建新的调度线程：
```markdown
本地IO调用；
基于底层系统同步调用的Syscall；
CGo方式调用C语言动态库中的调用IO或其它阻塞。
```
网络IO可以基于epoll的异步机制（或kqueue等异步机制），但对于一些系统函数并没有提供异步机制。例如常见的posix api中，对文件的操作就是同步操作。虽有开源的fileepoll来模拟异步文件操作。但Go的Syscall还是依赖底层的操作系统的API。系统API没有异步，Go也做不了异步化处理。

建议：把涉及到同步调用的goroutine，隔离到可控的goroutine中，而不是直接高并的goroutine调用。

## 2.3 高并发时避免共享对象互斥
传统多线程编程时，当并发冲突在4~8线程时，性能可能会出现拐点。Go中的推荐是不要通过共享内存来通讯，Go创建goroutine非常容易，当大量goroutine共享同一互斥对象时，也会在某一数量的goroutine出在拐点。

建议：goroutine尽量独立，无冲突地执行；若goroutine间存在冲突，则可以采分区来控制goroutine的并发个数，减少同一互斥对象冲突并发数。

# 3 其它优化

## 3.1 避免使用CGO或者减少CGO调用次数

GO可以调用C库函数，带有垃圾收集器且Go的栈是动态增长的，但这些无法与C无缝地对接。Go的环境转入C代码执行前，必须为C创建一个新的调用栈，把栈变量赋值给C调用栈，调用结束现拷贝回来。而这个调用开销也非常大，需要维护Go与C的调用上下文，两者调用栈的映射。相比直接的GO调用栈，单纯的调用栈可能有2个甚至3个数量级以上。

建议：尽量避免使用CGO，无法避免时，要减少跨CGO的调用次数。

## 3.2 减少[]byte与string之间转换，尽量采用[]byte来字符串处理
GO里面的string类型是一个不可变类型，不像c++中std:string，可以直接char*取值转化，指向同一地址内容；而GO中[]byte与string底层两个不同的结构，他们之间的转换存在实实在在的值对象拷贝，所以尽量减少这种不必要的转化

建议：存在字符串拼接等处理，尽量采用[]byte，例如：

func Prefix(b []byte) []byte {
    return append([]byte("hello", b...))
}

## 3.3 字符串的拼接优先考虑bytes.Buffer
由于string类型是一个不可变类型，通过连接符"+"拼接会创建新的string。GO中字符串拼接常见有如下几种方式：
```markdown
string + 操作 ：导致多次对象的分配与值拷贝
fmt.Sprintf ：会动态解析参数，效率好不哪去
strings.Join ：内部是[]byte的append
bytes.Buffer ：可以预先分配大小，减少对象分配与拷贝
```
建议：出于高性能要求，优先考虑bytes.Buffer，预先分配大小。另外，fmt.Sprintf可以简化不同类型转换与拼接。

参考：
1. Go语言内存分配器-FixAlloc
2. https://blog.golang.org/strings