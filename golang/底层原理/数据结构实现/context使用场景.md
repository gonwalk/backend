# golang中Context的数据结构及使用场景

## Context背景

作者：吴德宝AllenWu
链接：https://www.jianshu.com/p/e5df3cd0708b

golang在1.6.2的时候还没有自己的context，在1.7的版本中就把golang.org/x/net/context包加入到了官方的库中。golang 的 Context包，是专门用来简化对于处理单个请求的多个goroutine之间与请求域相关的数据、取消信号、截止时间等操作，这些操作可能涉及多个 API 调用。

比如有一个网络请求Request，每个Request都需要开启一个goroutine做一些事情，这些goroutine又可能会开启其他的goroutine。这样的话， 我们就可以通过Context，来跟踪这些goroutine，并且通过Context达到控制他们的目的，这就是Go语言为我们提供的Context，中文可以称之为“上下文”。

另外一个实际例子是，在Go服务器程序中，每个请求都会有一个goroutine去处理。然而，处理程序往往还需要创建额外的goroutine去访问后端资源，比如数据库、RPC服务等。由于这些goroutine都是在处理同一个请求，所以它们往往需要访问一些共享的资源，比如用户身份信息、认证token、请求截止时间等。而且如果请求超时或者被取消后，所有的goroutine都应该马上退出并且释放相关的资源。这种情况也需要用Context来为我们取消掉所有goroutine。

如果要使用可以通过 go get golang.org/x/net/context 命令获取这个包。

## Context 数据结构定义

Context的主要数据结构是一种嵌套的结构或者说是单向的继承关系的结构（，比如最初的context是一个小盒子，里面装了一些数据，之后从这个context继承下来的children就像在原本的context中又套上了一个盒子，然后里面装着一些自己的数据）。或者说context是一种分层的结构，根据使用场景的不同，每一层context都具备有一些不同的特性，这种层级式的组织也使得context易于扩展，职责清晰。

context 包的核心是Context，声明如下：
```go
type Context interface {

Deadline() (deadline time.Time, ok bool)

Done() <-chan struct{}

Err() error

Value(key interface{}) interface{}

}
```

Context是一个接口interface，在golang里面，它一共有4个方法。interface是一个使用非常广泛的结构，它可以接纳任何类型。
```
Deadline方法是获取设置的截止时间的意思，第一个返回参数是截止时间，到了这个时间点，Context会自动发起取消请求；第二个返回值是一个bool值变量ok，当ok==false时表示没有设置截止时间，如果想要取消的话，需要调用取消函数进行取消。

Done方法返回一个只读的chan，类型为struct{}，在goroutine中，如果该方法返回的chan可以读取，则意味着parent context已经发起了取消请求，我们通过Done方法收到这个信号后，就应该做清理操作，然后退出goroutine，释放资源。之后，Err 方法会返回一个错误，告知为什么 Context 被取消。

Err方法返回取消的错误原因，即为什么Context被取消。

Value方法获取该Context上绑定的值，是一个键值对，所以要通过一个Key才可以获取对应的值，这个值一般是线程安全的。
```

### Context 的实现方法

Context 虽然是个接口，但是并不需要使用方实现，golang内置的context 包，已经帮我们实现了2个方法，一般在代码中，开始上下文的时候都是以这两个作为最顶层的parent context，然后再衍生出子context。这些 Context 对象形成一棵树：当一个 Context 对象被取消时，继承自它的所有 Context的子对象 都会被取消。这两个方法：Background()和TODO()的实现方式如下：
```go
var (
    background = new(emptyCtx)

    todo = new(emptyCtx)
)

func Background() Context {
    return background
}

func TODO() Context {
    return todo
}
```
一个是Background，主要用于main函数、初始化以及测试代码中，作为Context这个树结构的最顶层的Context，也就是根Context，它不能被取消。

一个是TODO，如果我们不知道该使用什么Context的时候，可以使用这个，但是实际应用中，直接使用这个TODO的场景还不多。

它们两个本质上都是emptyCtx结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的Context。
```go
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {

    return
}

func (*emptyCtx) Done() <-chan struct{} {

    return nil
}

func (*emptyCtx) Err() error {

    return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {

    return nil
}
```
### Context 的 继承

有了上面的根Context，那么是如何衍生更多的子Context的呢？这就要靠context包为我们提供的With系列的函数了。
```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

func WithValue(parent Context, key, val interface{}) Context
```
通过这些函数，就创建了一颗Context树，树的每个节点都可以有任意多个子节点，节点层级可以有任意多个。

WithCancel函数，传递一个父Context作为参数，返回子Context，以及一个取消函数用来取消Context。

WithDeadline函数，和WithCancel差不多，它会多传递一个截止时间参数，意味着到了这个时间点，会自动取消Context，当然我们也可以不等到这个时候，可以提前通过取消函数进行取消。

WithTimeout和WithDeadline基本上一样，这个表示是超时自动取消，是多少时间后自动取消Context的意思。

WithValue函数和取消Context无关，它是为了生成一个绑定了一个键值对数据的Context，这个绑定的数据可以通过Context.Value方法访问到，这是我们实际用经常要用到的技巧，一般我们想要通过上下文来传递数据时，可以通过这个方法，如我们需要tarce追踪系统调用栈的时候。

### With 系列函数详解

#### WithCancel

context.WithCancel生成了一个withCancel的实例以及一个cancelFuc，这个函数就是用来关闭ctxWithCancel中的 Done channel 函数。

下面来分析下源码实现，首先看看初始化，如下：
```go
func newCancelCtx(parent Context) cancelCtx {
    return cancelCtx{
        Context: parent,
        done:    make(chan struct{}),
    }
}

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    c := newCancelCtx(parent)
    propagateCancel(parent, &c)
    return &c, func() { c.cancel(true, Canceled) }
}
```
newCancelCtx返回一个初始化的cancelCtx，cancelCtx结构体继承了Context，实现了canceler方法：

//*cancelCtx 和 *timerCtx 都实现了canceler接口，实现该接口的类型都可以被直接canceled
```go
type canceler interface {
    cancel(removeFromParent bool, err error)
    Done() <-chan struct{}
}


type cancelCtx struct {
    Context
    done chan struct{} // closed by the first cancel call.
    mu       sync.Mutex
    children map[canceler]bool // set to nil by the first cancel call
    err      error             // 当其被cancel时将会把err设置为非nil
}

func (c *cancelCtx) Done() <-chan struct{} {
    return c.done
}

func (c *cancelCtx) Err() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.err
}

func (c *cancelCtx) String() string {
    return fmt.Sprintf("%v.WithCancel", c.Context)
}
```
//核心是关闭c.done
//同时会设置c.err = err, c.children = nil
//依次遍历c.children，每个child分别cancel
//如果设置了removeFromParent，则将c从其parent的children中删除
```go
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    c.mu.Lock()
    if c.err != nil {
        c.mu.Unlock()
        return // already canceled
    }
    c.err = err
    close(c.done)
    for child := range c.children {
        // NOTE: acquiring the child's lock while holding parent's lock.
        child.cancel(false, err)
    }
    c.children = nil
    c.mu.Unlock()

    if removeFromParent {
        removeChild(c.Context, c) // 从此处可以看到 cancelCtx的Context项是一个类似于parent的概念
    }
}
```
可以看到，所有的children都存在一个map中；Done方法会返回其中的done channel， 而另外的cancel方法会关闭Done channel并且逐层向下遍历，关闭children的channel，并且将当前canceler从parent中移除。

WithCancel初始化一个cancelCtx的同时，还执行了propagateCancel方法，最后返回一个cancel function。

propagateCancel 方法定义如下：
```go
// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent Context, child canceler) {
    if parent.Done() == nil {
        return // parent is never canceled
    }
    if p, ok := parentCancelCtx(parent); ok {
        p.mu.Lock()
        if p.err != nil {
            // parent has already been canceled
            child.cancel(false, p.err)
        } else {
            if p.children == nil {
                p.children = make(map[canceler]struct{})
            }
            p.children[child] = struct{}{}
        }
        p.mu.Unlock()
    } else {
        go func() {
            select {
            case <-parent.Done():
                child.cancel(false, parent.Err())
            case <-child.Done():
            }
        }()
    }
}
```
propagateCancel 的含义就是传递cancel，从当前传入的parent开始（包括该parent），向上查找最近的一个可以被cancel的parent， 如果找到的parent已经被cancel，则将方才传入的child树给cancel掉，否则，将child节点直接连接为找到的parent的children中（Context字段不变，即向上的父亲指针不变，但是向下的孩子指针变直接了）； 如果没有找到最近的可以被cancel的parent，即其上都不可被cancel，则启动一个goroutine等待传入的parent终止，则cancel传入的child树，或者等待传入的child终结。

#### WithDeadLine

在withCancel的基础上进行的扩展，如果时间到了之后就进行cancel的操作，具体的操作流程基本上与withCancel一致，只不过控制cancel函数调用的时机是有一个timeout的channel所控制的。

### Context 使用原则 和 技巧

不要把Context放在结构体中，要以参数的方式传递，parent Context一般为Background
应该要把Context作为第一个参数传递给入口请求和出口请求链路上的每一个函数，放在第一位，变量名建议都统一，如ctx。
给一个函数方法传递Context的时候，不要传递nil，否则在tarce追踪的时候，就会断了连接
Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递
Context是线程安全的，可以放心的在多个goroutine中传递
可以把一个 Context 对象传递给任意个数的 gorotuine，对它执行 取消 操作时，所有 goroutine 都会接收到取消信号。

### Context的常用方法实例

1.调用Context Done方法取消

```
func Stream(ctx context.Context, out chan<- Value) error {

    for {
        v, err := DoSomething(ctx)

        if err != nil {
            return err
        }
        select {
        case <-ctx.Done():

            return ctx.Err()
        case out <- v:
        }
    }
}
```

2.通过 context.WithValue 来传值

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    valueCtx := context.WithValue(ctx, key, "add value")

    go watch(valueCtx)
    time.Sleep(10 * time.Second)
    cancel()

    time.Sleep(5 * time.Second)
}

func watch(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            //get value
            fmt.Println(ctx.Value(key), "is cancel")

            return
        default:
            //get value
            fmt.Println(ctx.Value(key), "int goroutine")

            time.Sleep(2 * time.Second)
        }
    }
}
```
3.超时取消 context.WithTimeout

```go
package main

import (
    "fmt"
    "sync"
    "time"

    "golang.org/x/net/context"
)

var (
    wg sync.WaitGroup
)

func work(ctx context.Context) error {
    defer wg.Done()

    for i := 0; i < 1000; i++ {
        select {
        case <-time.After(2 * time.Second):
            fmt.Println("Doing some work ", i)

        // we received the signal of cancelation in this channel
        case <-ctx.Done():
            fmt.Println("Cancel the context ", i)
            return ctx.Err()
        }
    }
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
    defer cancel()

    fmt.Println("Hey, I'm going to do some work")

    wg.Add(1)
    go work(ctx)
    wg.Wait()

    fmt.Println("Finished. I'm going home")
}
```

4.截止时间 取消 context.WithDeadline
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    d := time.Now().Add(1 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), d)

    // Even though ctx will be expired, it is good practice to call its
    // cancelation function in any case. Failure to do so may keep the
    // context and its parent alive longer than necessary.
    defer cancel()

    select {
    case <-time.After(2 * time.Second):
        fmt.Println("oversleep")
    case <-ctx.Done():
        fmt.Println(ctx.Err())
    }
}
```

## Context使用场景

文章链接：https://www.cnblogs.com/yjf512/p/10399190.html
作者：叶剑峰，主页：http://www.cnblogs.com/yjf512/

context在Go1.7之后就进入标准库中了。它主要的用处如果用一句话来说，是在于控制goroutine的生命周期。当一个计算任务被goroutine承接了之后，由于某种原因（超时，或者强制退出）我们希望中止这个goroutine的计算任务，那么就用得到这个Context了。

关于Context的四种结构，CancelContext,TimeoutContext,DeadLineContext,ValueContext的使用在这一篇快速掌握 Golang context 包已经说的很明白了。

本文主要来盘一盘golang中context的一些使用场景：

### 场景一：RPC调用


在主goroutine上有4个RPC，RPC2/3/4是并行请求的，我们这里希望在RPC2请求失败之后，直接返回错误，并且让RPC3/4停止继续计算。这个时候，就使用的到Context。

这个的具体实现如下面的代码。
```go
package main

import (
    "context"
    "sync"
    "github.com/pkg/errors"
)

func Rpc(ctx context.Context, url string) error {
    result := make(chan int)
    err := make(chan error)

    go func() {
        // 进行RPC调用，并且返回是否成功，成功通过result传递成功信息，错误通过error传递错误信息
        isSuccess := true
        if isSuccess {
            result <- 1
        } else {
            err <- errors.New("some error happen")
        }
    }()

    select {
        case <- ctx.Done():
            // 其他RPC调用调用失败
            return ctx.Err()
        case e := <- err:
            // 本RPC调用失败，返回错误信息
            return e
        case <- result:
            // 本RPC调用成功，不返回错误信息
            return nil
    }
}


func main() {

    ctx, cancel := context.WithCancel(context.Background())     // func WithCancel(parent Context) (ctx Context, cancel CancelFunc)，其中返回的cacel是一个函数，在main函数的goroutine中，当Rpc()调用出错时被调用cacel()。

    // RPC1调用
    err := Rpc(ctx, "http://rpc_1_url")
    if err != nil {
        return
    }

    wg := sync.WaitGroup{}

    // RPC2调用
    wg.Add(1)
    go func(){
        defer wg.Done()
        err := Rpc(ctx, "http://rpc_2_url")
        if err != nil {
            cancel()
        }
    }()

    // RPC3调用
    wg.Add(1)
    go func(){
        defer wg.Done()
        err := Rpc(ctx, "http://rpc_3_url")
        if err != nil {
            cancel()
        }
    }()

    // RPC4调用
    wg.Add(1)
    go func(){
        defer wg.Done()
        err := Rpc(ctx, "http://rpc_4_url")
        if err != nil {
            cancel()
        }
    }()

    wg.Wait()
}
```
当然我这里使用了waitGroup来保证main函数在所有RPC调用完成之后才退出。

```go
    ctx, cancel := context.WithCancel(context.Background())     // func WithCancel(parent Context) (ctx Context, cancel CancelFunc)，其中返回的cacel是一个函数，在main函数的goroutine中，当Rpc()调用出错时被调用cacel()。Background returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline. It is typically used by the main function, initialization, and tests, and as the top-level Context for incoming requests.
```

在Rpc函数中（即Rpc(ctx context.Context, url string)），第一个参数是一个Context，对于这个Context，形象的说就像一个传话筒，在创建Context的时候（context.Background()），返回了一个听声器（ctx）和话筒（cancel函数）。所有的goroutine都拿着这个听声器（ctx），当主goroutine想要告诉所有goroutine要结束的时候，通过cancel函数把结束的信息告诉给所有的goroutine。当然所有的goroutine都需要内置处理这个听声器结束信号的逻辑（<- ctx.Done()）。我们可以看Rpc函数内部，通过一个select来判断ctx的done和当前的rpc调用哪个先结束。

这个waitGroup和其中一个RPC调用就通知所有RPC的逻辑，其实有一个包已经帮我们做好了————[errorGroup](https://godoc.org/golang.org/x/sync/errgroup)，这个errorGroup包的使用可以看该包的test例子。

有人可能会担心我们这里的cancel()会被多次调用，context包的cancel调用是幂等的。可以放心多次调用。

我们不妨品一下，对于上面的Rpc函数，实际上是一个“阻塞式”的请求，这个请求如果是使用http.Get或者http.Post来实现，实际上Rpc函数的Goroutine结束了，内部的那个实际的http.Get却没有结束。所以，需要理解下，这里的函数最好是“非阻塞”的，比如是http.Do，然后可以通过某种方式进行中断。比如像这篇文章[Cancel http.Request using Context](https://medium.com/@ferencfbin/golang-cancel-http-request-using-context-1f45aeba6464)中的这个例子：
```go
func httpRequest(
  ctx context.Context,
  client *http.Client,
  req *http.Request,
  respChan chan []byte,
  errChan chan error
) {
  req = req.WithContext(ctx)
  tr := &http.Transport{}
  client.Transport = tr
  go func() {
    resp, err := client.Do(req)
    if err != nil {
      errChan <- err
    }
    if resp != nil {
      defer resp.Body.Close()
      respData, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        errChan <- err
      }
      respChan <- respData
    } else {
      errChan <- errors.New("HTTP request failed")
    }
  }()
  for {
    select {
    case <-ctx.Done():
      tr.CancelRequest(req)
      errChan <- errors.New("HTTP request cancelled")
      return
    case <-errChan:
      tr.CancelRequest(req)
      return
    }
  }
}
```
它使用了http.Client.Do，然后接收到ctx.Done的时候，通过调用transport.CancelRequest来进行结束。
我们还可以参考net/dail/DialContext
换而言之，如果你希望你实现的包是“可中止/可控制”的，那么你在你包实现的函数里面，最好是能接收一个Context函数，并且处理了Context.Done。

场景二：PipeLine
pipeline模式就是流水线模型，流水线上的几个工人，有n个产品，一个一个产品进行组装。其实pipeline模型的实现和Context并无关系，没有context我们也能用chan实现pipeline模型。但是对于整条流水线的控制，则是需要使用上Context的。这篇文章Pipeline Patterns in Go的例子是非常好的说明。这里就大致对这个代码进行下说明。

runSimplePipeline的流水线工人有三个，lineListSource负责将参数一个个分割进行传输，lineParser负责将字符串处理成int64,sink根据具体的值判断这个数据是否可用。他们所有的返回值基本上都有两个chan，一个用于传递数据，一个用于传递错误。（<-chan string, <-chan error）输入基本上也都有两个值，一个是Context，用于传声控制的，一个是(in <-chan)输入产品的。

我们可以看到，这三个工人的具体函数里面，都使用switch处理了case <-ctx.Done()。这个就是生产线上的命令控制。

func lineParser(ctx context.Context, base int, in <-chan string) (
    <-chan int64, <-chan error, error) {
    ...
    go func() {
        defer close(out)
        defer close(errc)

        for line := range in {

            n, err := strconv.ParseInt(line, base, 64)
            if err != nil {
                errc <- err
                return
            }

            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out, errc, nil
}
场景三：超时请求
我们发送RPC请求的时候，往往希望对这个请求进行一个超时的限制。当一个RPC请求超过10s的请求，自动断开。当然我们使用CancelContext，也能实现这个功能（开启一个新的goroutine，这个goroutine拿着cancel函数，当时间到了，就调用cancel函数）。

鉴于这个需求是非常常见的，context包也实现了这个需求：timerCtx。具体实例化的方法是 WithDeadline 和 WithTimeout。

具体的timerCtx里面的逻辑也就是通过time.AfterFunc来调用ctx.cancel的。

官方的例子：

package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
    defer cancel()

    select {
    case <-time.After(1 * time.Second):
        fmt.Println("overslept")
    case <-ctx.Done():
        fmt.Println(ctx.Err()) // prints "context deadline exceeded"
    }
}
在http的客户端里面加上timeout也是一个常见的办法

uri := "https://httpbin.org/delay/3"
req, err := http.NewRequest("GET", uri, nil)
if err != nil {
    log.Fatalf("http.NewRequest() failed with '%s'\n", err)
}

ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
req = req.WithContext(ctx)

resp, err := http.DefaultClient.Do(req)
if err != nil {
    log.Fatalf("http.DefaultClient.Do() failed with:\n'%s'\n", err)
}
defer resp.Body.Close()
在http服务端设置一个timeout如何做呢？

package main

import (
    "net/http"
    "time"
)

func test(w http.ResponseWriter, r *http.Request) {
    time.Sleep(20 * time.Second)
    w.Write([]byte("test"))
}


func main() {
    http.HandleFunc("/", test)
    timeoutHandler := http.TimeoutHandler(http.DefaultServeMux, 5 * time.Second, "timeout")
    http.ListenAndServe(":8080", timeoutHandler)
}
我们看看TimeoutHandler的内部，本质上也是通过context.WithTimeout来做处理。

func (h *timeoutHandler) ServeHTTP(w ResponseWriter, r *Request) {
  ...
        ctx, cancelCtx = context.WithTimeout(r.Context(), h.dt)
        defer cancelCtx()
    ...
    go func() {
    ...
        h.handler.ServeHTTP(tw, r)
    }()
    select {
    ...
    case <-ctx.Done():
        ...
    }
}
场景四：HTTP服务器的request互相传递数据
context还提供了valueCtx的数据结构。

这个valueCtx最经常使用的场景就是在一个http服务器中，在request中传递一个特定值，比如有一个中间件，做cookie验证，然后把验证后的用户名存放在request中。

我们可以看到，官方的request里面是包含了Context的，并且提供了WithContext的方法进行context的替换。

package main

import (
    "net/http"
    "context"
)

type FooKey string

var UserName = FooKey("user-name")
var UserId = FooKey("user-id")

func foo(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), UserId, "1")
        ctx2 := context.WithValue(ctx, UserName, "yejianfeng")
        next(w, r.WithContext(ctx2))
    }
}

func GetUserName(context context.Context) string {
    if ret, ok := context.Value(UserName).(string); ok {
        return ret
    }
    return ""
}

func GetUserId(context context.Context) string {
    if ret, ok := context.Value(UserId).(string); ok {
        return ret
    }
    return ""
}

func test(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("welcome: "))
    w.Write([]byte(GetUserId(r.Context())))
    w.Write([]byte(" "))
    w.Write([]byte(GetUserName(r.Context())))
}

func main() {
    http.Handle("/", foo(test))
    http.ListenAndServe(":8080", nil)
}
在使用ValueCtx的时候需要注意一点，这里的key不应该设置成为普通的String或者Int类型，为了防止不同的中间件对这个key的覆盖。最好的情况是每个中间件使用一个自定义的key类型，比如这里的FooKey，而且获取Value的逻辑尽量也抽取出来作为一个函数，放在这个middleware的同包中。这样，就会有效避免不同包设置相同的key的冲突问题了。

参考
快速掌握 Golang context 包：https://deepzz.com/post/golang-context-package-notes.html

视频笔记：如何正确使用 Context - Jack Lindamood：https://blog.lab99.org/post/golang-2017-10-27-video-how-to-correctly-use-package-context.html

Go Concurrency Patterns: Context：https://blog.golang.org/context

Cancel http.Request using Context：https://medium.com/@ferencfbin/golang-cancel-http-request-using-context-1f45aeba6464

Pipeline Patterns in Go：https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d










第三部分（21~30）
第二十一题：简单密码破解
第二十二题：汽水瓶
第二十三题：删除字符串中出现次数最少的字符






第四部分（31~40）
https://blog.csdn.net/qq_25220145/article/details/78414774?depth_1-utm_source=distribute.pc_relevant.none-task&utm_source=distribute.pc_relevant.none-task
第三十一题：[中级]单词倒排
第三十二题：字符串运用-密码截取
第三十三题：整数与IP地址间的转换
第三十四题：图片整理
第三十五题：蛇形矩阵
第三十六题：字符串加密
第三十七题：统计每个月兔子的总数
第三十八题：求小球落地5次后所经历的路程和第5次反弹的高度
第三十九题：判断两个IP是否属于同一子网
第四十题：输入一行字符，分别统计出包含英文字母、空格、数字和其它字符的个数
