# golang性能分析工具

## golang pprof

当你的golang程序在运行过程中消耗了超出你理解的内存时，你就需要搞明白，到底是 程序中哪些代码导致了这些内存消耗。此时golang编译好的程序对你来说是个黑盒，该 如何搞清其中的内存使用呢？幸好golang已经内置了一些机制来帮助我们进行分析和追踪。

此时，通常我们可以采用golang的pprof来帮助我们分析golang进程的内存使用。

### pprof 实例

通常我们采用http api来将pprof信息暴露出来以供分析，我们可以采用net/http/pprof 这个package。下面是一个简单的示例：

// pprof 的init函数会将pprof里的一些handler注册到http.DefaultServeMux上
// 当不使用http.DefaultServeMux来提供http api时，可以查阅其init函数，自己注册handler
```go
import _ "net/http/pprof"

go func() {
    http.ListenAndServe("0.0.0.0:8080", nil)
}()
```
此时我们可以启动进程，然后访问http://localhost:8080/debug/pprof/可以看到一个简单的 页面，页面上显示: 注意: 以下的全部数据，包括go tool pprof 采集到的数据都依赖进程中的pprof采样率，默认512kb进行 一次采样，当我们认为数据不够细致时，可以调节采样率runtime.MemProfileRate，但是采样率越低，进 程运行速度越慢。
```
/debug/pprof/

profiles:
0         block
136840    goroutine
902       heap
0         mutex
40        threadcreate

full goroutine stack dump
```

上面简单暴露出了几个内置的Profile统计项。例如有136840个goroutine在运行，点击相关链接 可以看到详细信息。

当我们分析内存相关的问题时，可以点击heap项，进入http://127.0.0.1:8080/debug/pprof/heap?debug=1 可以查看具体的显示：

heap profile: 3190: 77516056 [54762: 612664248] @ heap/1048576
1: 29081600 [1: 29081600] @ 0x89368e 0x894cd9 0x8a5a9d 0x8a9b7c 0x8af578 0x8b4441 0x8b4c6d 0x8b8504 0x8b2bc3 0x45b1c1
#    0x89368d    github.com/syndtr/goleveldb/leveldb/memdb.(*DB).Put+0x59d
#    0x894cd8    xxxxx/storage/internal/memtable.(*MemTable).Set+0x88
#    0x8a5a9c    xxxxx/storage.(*snapshotter).AppendCommitLog+0x1cc
#    0x8a9b7b    xxxxx/storage.(*store).Update+0x26b
#    0x8af577    xxxxx/config.(*config).Update+0xa7
#    0x8b4440    xxxxx/naming.(*naming).update+0x120
#    0x8b4c6c    xxxxx/naming.(*naming).instanceTimeout+0x27c
#    0x8b8503    xxxxx/naming.(*naming).(xxxxx/naming.instanceTimeout)-fm+0x63

......

# runtime.MemStats
# Alloc = 2463648064
# TotalAlloc = 31707239480
# Sys = 4831318840
# Lookups = 2690464
# Mallocs = 274619648
# Frees = 262711312
# HeapAlloc = 2463648064
# HeapSys = 3877830656
# HeapIdle = 854990848
# HeapInuse = 3022839808
# HeapReleased = 0
# HeapObjects = 11908336
# Stack = 655949824 / 655949824
# MSpan = 63329432 / 72040448
# MCache = 38400 / 49152
# BuckHashSys = 1706593
# GCSys = 170819584
# OtherSys = 52922583
# NextGC = 3570699312
# PauseNs = [1052815 217503 208124 233034 1146462 456882 1098525 530706 551702 419372 768322 596273 387826 455807 563621 587849 416204 599143 572823 488681 701731 656358 2476770 12141392 5827253 3508261 1715582 1295487 908563 788435 718700 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
# NumGC = 31
# DebugGC = false
其中显示的内容会比较多，但是主体分为2个部分: 第一个部分打印为通过runtime.MemProfile()获取的runtime.MemProfileRecord记录。 其含义为：

heap profile: 3190(inused objects): 77516056(inused bytes) [54762(alloc objects): 612664248(alloc bytes)] @ heap/1048576(2*MemProfileRate)
1: 29081600 [1: 29081600] (前面4个数跟第一行的一样，此行以后是每次记录的，后面的地址是记录中的栈指针)@ 0x89368e 0x894cd9 0x8a5a9d 0x8a9b7c 0x8af578 0x8b4441 0x8b4c6d 0x8b8504 0x8b2bc3 0x45b1c1
#    0x89368d    github.com/syndtr/goleveldb/leveldb/memdb.(*DB).Put+0x59d 栈信息
第二部分就比较好理解，打印的是通过runtime.ReadMemStats()读取的runtime.MemStats信息。 我们可以重点关注一下

Sys 进程从系统获得的内存空间，虚拟地址空间。
HeapAlloc 进程堆内存分配使用的空间，通常是用户new出来的堆对象，包含未被gc掉的。
HeapSys 进程从系统获得的堆内存，因为golang底层使用TCmalloc机制，会缓存一部分堆内存，虚拟地址空间。
PauseNs 记录每次gc暂停的时间(纳秒)，最多记录256个最新记录。
NumGC 记录gc发生的次数。
相信，对pprof不了解的用户看了以上内容，很难获得更多的有用信息。因此我们需要引用更多工具来帮助 我们更加简单的解读pprof内容。

## go tool

我们可以采用go tool pprof -inuse_space http://127.0.0.1:8080/debug/pprof/heap命令连接到进程中 查看正在使用的一些内存相关信息，此时我们得到一个可以交互的命令行。

我们可以看数据top10来查看正在使用的对象较多的10个函数入口。通常用来检测有没有不符合预期的内存 对象引用。

(pprof) top10
1355.47MB of 1436.26MB total (94.38%)
Dropped 371 nodes (cum <= 7.18MB)
Showing top 10 nodes out of 61 (cum >= 23.50MB)
      flat  flat%   sum%        cum   cum%
  512.96MB 35.71% 35.71%   512.96MB 35.71%  net/http.newBufioWriterSize
  503.93MB 35.09% 70.80%   503.93MB 35.09%  net/http.newBufioReader
  113.04MB  7.87% 78.67%   113.04MB  7.87%  runtime.rawstringtmp
   55.02MB  3.83% 82.50%    55.02MB  3.83%  runtime.malg
   45.01MB  3.13% 85.64%    45.01MB  3.13%  xxxxx/storage.(*Node).clone
   26.50MB  1.85% 87.48%    52.50MB  3.66%  context.WithCancel
   25.50MB  1.78% 89.26%    83.58MB  5.82%  runtime.systemstack
   25.01MB  1.74% 91.00%    58.51MB  4.07%  net/http.readRequest
      25MB  1.74% 92.74%    29.03MB  2.02%  runtime.mapassign
   23.50MB  1.64% 94.38%    23.50MB  1.64%  net/http.(*Server).newConn
然后我们在用go tool pprof -alloc_space http://127.0.0.1:8080/debug/pprof/heap命令链接程序来查看 内存对象分配的相关情况。然后输入top来查看累积分配内存较多的一些函数调用:

(pprof) top
523.38GB of 650.90GB total (80.41%)
Dropped 342 nodes (cum <= 3.25GB)
Showing top 10 nodes out of 106 (cum >= 28.02GB)
      flat  flat%   sum%        cum   cum%
  147.59GB 22.68% 22.68%   147.59GB 22.68%  runtime.rawstringtmp
  129.23GB 19.85% 42.53%   129.24GB 19.86%  runtime.mapassign
   48.23GB  7.41% 49.94%    48.23GB  7.41%  bytes.makeSlice
   46.25GB  7.11% 57.05%    71.06GB 10.92%  encoding/json.Unmarshal
   31.41GB  4.83% 61.87%   113.86GB 17.49%  net/http.readRequest
   30.55GB  4.69% 66.57%   171.20GB 26.30%  net/http.(*conn).readRequest
   22.95GB  3.53% 70.09%    22.95GB  3.53%  net/url.parse
   22.70GB  3.49% 73.58%    22.70GB  3.49%  runtime.stringtoslicebyte
   22.70GB  3.49% 77.07%    22.70GB  3.49%  runtime.makemap
   21.75GB  3.34% 80.41%    28.02GB  4.31%  context.WithCancel
可以看出string-[]byte相互转换、分配map、bytes.makeSlice、encoding/json.Unmarshal等调用累积分配的内存较多。 此时我们就可以review代码，如何减少这些相关的调用，或者优化相关代码逻辑。

当我们不明确这些调用时是被哪些函数引起的时，我们可以输入top -cum来查找，-cum的意思就是，将函数调用关系 中的数据进行累积，比如A函数调用的B函数，则B函数中的内存分配量也会累积到A上面，这样就可以很容易的找出调用链。

(pprof) top20 -cum
322890.40MB of 666518.53MB total (48.44%)
Dropped 342 nodes (cum <= 3332.59MB)
Showing top 20 nodes out of 106 (cum >= 122316.23MB)
      flat  flat%   sum%        cum   cum%
         0     0%     0% 643525.16MB 96.55%  runtime.goexit
 2184.63MB  0.33%  0.33% 620745.26MB 93.13%  net/http.(*conn).serve
         0     0%  0.33% 435300.50MB 65.31%  xxxxx/api/server.(*HTTPServer).ServeHTTP
 5865.22MB  0.88%  1.21% 435300.50MB 65.31%  xxxxx/api/server/router.(*httpRouter).ServeHTTP
         0     0%  1.21% 433121.39MB 64.98%  net/http.serverHandler.ServeHTTP
         0     0%  1.21% 430456.29MB 64.58%  xxxxx/api/server/filter.(*chain).Next
   43.50MB 0.0065%  1.21% 429469.71MB 64.43%  xxxxx/api/server/filter.TransURLTov1
         0     0%  1.21% 346440.39MB 51.98%  xxxxx/api/server/filter.Role30x
31283.56MB  4.69%  5.91% 175309.48MB 26.30%  net/http.(*conn).readRequest
         0     0%  5.91% 153589.85MB 23.04%  github.com/julienschmidt/httprouter.(*Router).ServeHTTP
         0     0%  5.91% 153589.85MB 23.04%  github.com/julienschmidt/httprouter.(*Router).ServeHTTP-fm
         0     0%  5.91% 153540.85MB 23.04%  xxxxx/api/server/router.(*httpRouter).Register.func1
       2MB 0.0003%  5.91% 153117.78MB 22.97%  xxxxx/api/server/filter.Validate
151134.52MB 22.68% 28.58% 151135.02MB 22.68%  runtime.rawstringtmp
         0     0% 28.58% 150714.90MB 22.61%  xxxxx/api/server/router/naming/v1.(*serviceRouter).(git.intra.weibo.com/platform/vintage/api/server/router/naming/v1.service)-fm
         0     0% 28.58% 150714.90MB 22.61%  xxxxx/api/server/router/naming/v1.(*serviceRouter).service
         0     0% 28.58% 141200.76MB 21.18%  net/http.Redirect
132334.96MB 19.85% 48.44% 132342.95MB 19.86%  runtime.mapassign
      42MB 0.0063% 48.44% 125834.16MB 18.88%  xxxxx/api/server/router/naming/v1.heartbeat
         0     0% 48.44% 122316.23MB 18.35%  xxxxxx/config.(*config).Lookup

如上所示，我们就很容易的查找到这些函数是被哪些函数调用的。

根据代码的调用关系，filter.TransURLTov1会调用filter.Role30x，但是他们之间的cum%差值有12.45%，因此 我们可以得知filter.TransURLTov1内部自己直接分配的内存量达到了整个进程分配内存总量的12.45%，这可是一个 值得大大优化的地方。

然后我们可以输入命令web，其会给我们的浏览器弹出一个.svg图片，其会把这些累积关系画成一个拓扑图，提供给 我们。或者直接执行go tool pprof -alloc_space -cum -svg http://127.0.0.1:8080/debug/pprof/heap > heap.svg来生 成heap.svg图片。

下面我们取一个图片中的一个片段进行分析:
```
golang-memory-pprof.png
```
每一个方块为pprof记录的一个函数调用栈，指向方块的箭头上的数字是记录的该栈累积分配的内存向，从方块指出的 箭头上的数字为该函数调用的其他函数累积分配的内存。他们之间的差值可以简单理解为本函数除调用其他函数外，自 身分配的。方块内部的数字也体现了这一点，其数字为:(自身分配的内存 of 该函数累积分配的内存)。

--inuse/alloc_space --inuse/alloc_objects区别
通常情况下：

用--inuse_space来分析程序常驻内存的占用情况;
用--alloc_objects来分析内存的临时分配情况，可以提高程序的运行速度。
go-torch
除了直接使用go tool pprof外，我们还可以使用更加直观了火焰图 。因此我们可以直接使用go-torch来生成golang程序的火焰图，该工具也直接 依赖pprof/go tool pprof等。该工具的相关安装请看该项目的介绍。该软件的a4daa2b 以后版本才支持内存的profiling。

我们可以使用

go-torch -alloc_space http://127.0.0.1:8080/debug/pprof/heap --colors=mem
go-torch -inuse_space http://127.0.0.1:8080/debug/pprof/heap --colors=mem
注意:-alloc_space/-inuse_space参数与-u/-b等参数有冲突，使用了-alloc_space/-inuse_space后请将pprof的 资源直接追加在参数后面，而不要使用-u/-b参数去指定，这与go-torch的参数解析问题有关，看过其源码后既能明白。 同时还要注意，分析内存的URL一定是heap结尾的，因为默认路径是profile的，其用来分析cpu相关问题。

通过上面2个命令，我们就可以得到alloc_space/inuse_space含义的2个火焰图，例如 alloc_space.svg/inuse_space.svg。 我们可以使用浏览器观察这2张图，这张图，就像一个山脉的截面图，从下而上是每个函数的调用栈，因此山的高度跟函数 调用的深度正相关，而山的宽度跟使用/分配内存的数量成正比。我们只需要留意那些宽而平的山顶，这些部分通常是我们 需要优化的地方。

testing
当我们需要对go test中某些test/benchmark进行profiling时，我们可以使用类似的方法。例如我们可以先使用go test 内置的参数生成pprof数据，然后借助go tool pprof/go-torch来分析。

生成cpu、mem的pprof文件
go test -bench=BenchmarkStorageXXX -cpuprofile cpu.out -memprofile mem.out
此时会生成一个二进制文件和2个pprof数据文件，例如
storage.test cpu.out mem.out
然后使用go-torch来分析，二进制文件放前面
#分析cpu
go-torch storage.test cpu.out
#分析内存
go-torch --colors=mem -alloc_space storage.test mem.out
go-torch --colors=mem -inuse_space storage.test mem.out

## 优化建议
Debugging performance issues in Go programs 提供了一些常用的优化建议:

### 1 将多个小对象合并成一个大的对象

### 2 减少不必要的指针间接引用，多使用copy引用

例如使用bytes.Buffer代替*bytes.Buffer，因为使用指针时，会分配2个对象来完成引用。

### 3 局部变量逃逸时，将其聚合起来

这一点理论跟1相同，核心在于减少object的分配，减少gc的压力。 例如，以下代码
```go
for k, v := range m {
	k, v := k, v   // copy for capturing by the goroutine
	go func() {
		// use k and v
	}()
}
```
可以修改为:
```go
for k, v := range m {
	x := struct{ k, v string }{k, v}   // copy for capturing by the goroutine
	go func() {
		// use x.k and x.v
	}()
}
```
修改后，逃逸的对象变为了x，将k，v2个对象减少为1个对象。

### 4 []byte的预分配

当我们比较清楚的知道[]byte会到底使用多少字节，我们就可以采用一个数组来预分配这段内存。 例如:

type X struct {
    buf      []byte
    bufArray [16]byte // Buf usually does not grow beyond 16 bytes.
}

func MakeX() *X {
    x := &X{}
    // Preinitialize buf with the backing array.
    x.buf = x.bufArray[:0]
    return x
}

### 5 尽可能使用字节数少的类型

当我们的一些const或者计数字段不需要太大的字节数时，我们通常可以将其声明为int8类型。

### 6 减少不必要的指针引用

当一个对象不包含任何指针（注意：strings，slices，maps 和chans包含隐含的指针），时，对gc的扫描影响很小。 比如，1GB byte 的slice事实上只包含有限的几个object，不会影响垃圾收集时间。 因此，我们可以尽可能的减少指针的引用。

### 7 使用sync.Pool来缓存常用的对象

注意
go1.9、go1.9.2之间的版本go tool pprof引入了一个BUG，会导致上面的内存分析命令失败。 下面给出一种修复办法：
```
cd $GOROOT/src/cmd/vendor/github.com/google
rm pprof
git clone https://github.com/google/pprof.git #确保在版本`e82ee9addc1b6c8e1d667ed6de0194241e1e03b5`之后
rm $GOROOT/pkg/darwin_amd64/cmd/vendor/github.com/google/pprof
cd $GOROOT/src/cmd/pprof
go build
mv pprof $GOROOT/pkg/tool/darwin_amd64/pprof
```

# golang内存分析工具

分析golang内存占用情况：https://www.cnblogs.com/eliotgolang/p/9671589.html

性能优化是个永恒的话题，而很多时候我们在作性能优化的时候，往往基于代码上面的直觉，把所有能想到的优化都优化了一遍，不错过任何小的优化点，结果整个代码的逻辑变得极其复杂，而性能上面并没有太大的提升。事实上，性能问题往往集中在某些小点，有时候很小的改动就能有巨大的提升，所以问题的关键是是怎么去找出这些优化点，幸运的是 golang 在设计的时候就考虑了这个问题，原生提供了性能分析的工具，可以很方便地帮我们找到性能瓶颈

## pprof 简介

import "runtime/pprof"

pprof包以pprof可视化工具期望的格式书写运行时剖面数据。

golang 的性能分析库在 runtime/pprof 里，主要提供下面几个接口：
```go
import "runtime/pprof"

// 堆栈分析
func WriteHeapProfile(w io.Writer) error

// cpu分析
func StartCPUProfile(w io.Writer) error
func StopCPUProfile()
```
使用上面的方法比较简单，只需要将文件指针传给对应的函数即可，性能数据将写入到文件中，然后可以使用 golang 自带的 pprof 工具生成 svg、pdf 的可视化图，然后就可以很直观地从这些图里面看到主要的性能消耗了。

## 示例

### 代码示例

首先需要一个程序，然后在程序里面注入 pprof 代码，下面是一段示例代码，完整代码在：https://github.com/hatlonely/hellogolang/blob/master/cmd/pprof_runtime.go
```go
func main() {
    go doSomething1()
    go doSomething2()
    go doSomething3()

    if err := pprof.PPCmd("cpu 10s"); err != nil {
        panic(err)
    }

    if err := pprof.PPCmd("mem"); err != nil {
        panic(err)
    }
}
```
编译，运行上面代码会生成两个 pprof 文件，cpu.pprof.yyyymmddhhmmss 和 mem.pprof.yyyymmddhhmmss，编译运行的方法如下：
```
cd $GOPATH/src
git clone git@github.com:hatlonely/hellogolang.git
cd hellogolang
glide install
go build cmd/pprof_runtime.go
./pprof_runtime
```

### pprof 文件分析


pprof 文件是二进制的，不是给人读的，需要翻译一下，而 golang 原生就给我们提供了分析工具，直接执行下面命令即可，会生成一张很直观的 svg 图片，直接用 chrome 就可以打开，当然也可以生成别的格式（pdf，png 都可以），可以用 go tool pprof -h 命令查看支持的输出类型
```
go tool pprof -svg ./pprof_runtime cpu.pprof.201801301415 > cpu.svg
```
注意这个工具依赖于 graphviz 工具，Mac 上可用 brew install graphviz，centos系统使用命令 yum install graphviz 安装即可。

