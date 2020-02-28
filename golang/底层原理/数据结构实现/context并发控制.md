# 1. 前言

在上篇Golang高效实践之并发实践channel篇中我给大家介绍了Golang并发模型，详细的介绍了channel的用法，和用select管理channel。比如说我们可以用channel来控制几个goroutine的同步和退出时机，但是我们需要close channel通知其他接受者，当通知和通信的内容混在一起时往往比较复杂，需要把握好channel的读写时机，以及不能往已经关闭的channel中再写入数据。如果有没有一种更好的上下文控制机制呢？答案就是文章今天要介绍的context，context正是close channel的一种封装，通常用来控制上下文的同步。

# 2. Context介绍

Context包定义了Context类型，Context类型携带着deadline生命周期，和取消信号，并且可以携带用户自定义的参数值。通常用Context来控制上下文，Context通过参数一层层传递，或者传递context的派生，一旦Context被取消，所有由该Context派生的Context也会取消。WithCancel，WithDeadline，和WithTimeout函数可以从一个Context中派生另外一个Context和一个cancel函数。调用cancel函数可以取消由context派生出来的Context。cancel函数会释放context拥有的资源，所以当context不用时要尽快调用cancel。

Context应该作为函数的第一个参数，通常使用ctx命名，例如：
```go
func DoSomething(ctx context.Context, arg Arg) error {

// … use ctx …

}
```

不要传递nil context，即便接受的函数允许我们这样做也不要传递nil context。如果你不确定用哪个context的话可以传递context.TODO。

同一个context可以在不同的goroutine中访问，context是线程安全的。

## 2.1 Context结构定义

```go
type Context interface {
        // Deadline returns the time when work done on behalf of this context
        // should be canceled. Deadline returns ok==false when no deadline is
        // set. Successive calls to Deadline return the same results.
        Deadline() (deadline time.Time, ok bool)

        // Done returns a channel that's closed when work done on behalf of this
        // context should be canceled. Done may return nil if this context can
        // never be canceled. Successive calls to Done return the same value.
        //
        // WithCancel arranges for Done to be closed when cancel is called;
        // WithDeadline arranges for Done to be closed when the deadline
        // expires; WithTimeout arranges for Done to be closed when the timeout
        // elapses.
        //
        // Done is provided for use in select statements:
        //
        //  // Stream generates values with DoSomething and sends them to out
        //  // until DoSomething returns an error or ctx.Done is closed.
        //  func Stream(ctx context.Context, out chan<- Value) error {
        //      for {
        //          v, err := DoSomething(ctx)
        //          if err != nil {
        //              return err
        //          }
        //          select {
        //          case <-ctx.Done():
        //              return ctx.Err()
        //          case out <- v:
        //          }
        //      }
        //  }
        //
        // See https://blog.golang.org/pipelines for more examples of how to use
        // a Done channel for cancelation.
        Done() <-chan struct{}

        // If Done is not yet closed, Err returns nil.
        // If Done is closed, Err returns a non-nil error explaining why:
        // Canceled if the context was canceled
        // or DeadlineExceeded if the context's deadline passed.
        // After Err returns a non-nil error, successive calls to Err return the same error.
        Err() error

        // Value returns the value associated with this context for key, or nil
        // if no value is associated with key. Successive calls to Value with
        // the same key returns the same result.
        //
        // Use context values only for request-scoped data that transits
        // processes and API boundaries, not for passing optional parameters to
        // functions.
        //
        // A key identifies a specific value in a Context. Functions that wish
        // to store values in Context typically allocate a key in a global
        // variable then use that key as the argument to context.WithValue and
        // Context.Value. A key can be any type that supports equality;
        // packages should define keys as an unexported type to avoid
        // collisions.
        //
        // Packages that define a Context key should provide type-safe accessors
        // for the values stored using that key:
        //
        //     // Package user defines a User type that's stored in Contexts.
        //     package user
        //
        //     import "context"
        //
        //     // User is the type of value stored in the Contexts.
        //     type User struct {...}
        //
        //     // key is an unexported type for keys defined in this package.
        //     // This prevents collisions with keys defined in other packages.
        //     type key int
        //
        //     // userKey is the key for user.User values in Contexts. It is
        //     // unexported; clients use user.NewContext and user.FromContext
        //     // instead of using this key directly.
        //     var userKey key
        //
        //     // NewContext returns a new Context that carries value u.
        //     func NewContext(ctx context.Context, u *User) context.Context {
        //         return context.WithValue(ctx, userKey, u)
        //     }
        //
        //     // FromContext returns the User value stored in ctx, if any.
        //     func FromContext(ctx context.Context) (*User, bool) {
        //         u, ok := ctx.Value(userKey).(*User)
        //         return u, ok
        //     }
        Value(key interface{}) interface{}
}
```

## 2.2 WithCancel函数

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```
WithCancel函数返回parent的一份拷贝和一个新的Done channel。当concel 函数被调用的时候或者parent的Done channel被关闭时（cancel被调用），context的Done channel将会被关闭。取消context将会释放context相关的资源，所以当context完成时代码应该尽快调用cancel方法。例如：
```go
package main

import (
    "context"
    "fmt"
)

func main() {
    // gen generates integers in a separate goroutine and
    // sends them to the returned channel.
    // The callers of gen need to cancel the context once
    // they are done consuming generated integers not to leak
    // the internal goroutine started by gen.
    gen := func(ctx context.Context) <-chan int {
        dst := make(chan int)
        n := 1
        go func() {
            for {
                select {
                case <-ctx.Done():
                    return // returning not to leak the goroutine
                case dst <- n:
                    n++
                }
            }
        }()
        return dst
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // cancel when we are finished consuming integers

    for n := range gen(ctx) {
        fmt.Println(n)
        if n == 5 {
            break
        }
    }
}
```

## 2.3 WithDeadline函数

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

WithDeadline函数返回parent context调整deadline之后的拷贝，如果parent的deadline比要调整的d更早，那么派生出来的context的deadline就等于parent的deadline。当deadline过期或者cancel函数被调用时，又或者parent的cancel函数被调用时，context的Done channel将会被触发。例如：
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    d := time.Now().Add(50 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), d)

    // Even though ctx will be expired, it is good practice to call its
    // cancelation function in any case. Failure to do so may keep the
    // context and its parent alive longer than necessary.
    defer cancel()

    select {
    case <-time.After(1 * time.Second):
        fmt.Println("overslept")
    case <-ctx.Done():
        fmt.Println(ctx.Err())
    }

}
```
Err方法会返回context退出的原因，这里是context deadline exceeded。

## 2.4 WithTimeout函数

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
WithTimeout相当于调用WithDeadline(parent, time.Now().Add(timeout)),例如：
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Pass a context with a timeout to tell a blocking function that it
    // should abandon its work after the timeout elapses.
    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
    defer cancel()

    select {
    case <-time.After(1 * time.Second):
        fmt.Println("overslept")
    case <-ctx.Done():
        fmt.Println(ctx.Err()) // prints "context deadline exceeded"
    }

}
```

## 2.5 Background函数

func Background() Context
Backgroud函数返回一个非nil的空context。该context不会cancel，没有值，没有deadline。通常在main函数中调用，初始化或者测试，作为顶级的context。

## 2.6WithValue函数

func WithValue(parent Context, key, val interface{}) Context
WithValue函数返回parent的拷贝，并且key对应的值是value。例如：
```go
package main

import (
    "context"
    "fmt"
)

func main() {
    type favContextKey string

    f := func(ctx context.Context, k favContextKey) {
        if v := ctx.Value(k); v != nil {
            fmt.Println("found value:", v)
            return
        }
        fmt.Println("key not found:", k)
    }

    k := favContextKey("language")
    ctx := context.WithValue(context.Background(), k, "Go")

    f(ctx, k)
    f(ctx, favContextKey("color"))

}
```

# 3. webserver实战

有了上面的理论知识后，我将给大家讲解一个webserver的编码，其中就用到context的超时特性，以及上下文同步等。代码放在github上面，是从google search仓库中fork出来并做了一些改动。该项目的代码用到go module来组织代码，如果对go module不熟悉的同学可以参考我的这篇博客。

server.go文件是main包，里面包含一个http server：

func main() {
    http.HandleFunc("/search", handleSearch)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
例如通过/search?q=golang&timeout=1s访问8080端口将会调用handle函数handleSearch来处理，handleSearch会解析出来要查询的关键字golang，并且指定的超时时间是1s。该timeout参数会用于生成带有timeout属性的context，该context会贯穿整个请求的上下文，当超时时间触发时会终止search。

```go
func handleSearch(w http.ResponseWriter, req *http.Request) {
    // ctx is the Context for this handler. Calling cancel closes the
    // ctx.Done channel, which is the cancellation signal for requests
    // started by this handler.
    var (
        ctx    context.Context
        cancel context.CancelFunc
    )
    timeout, err := time.ParseDuration(req.FormValue("timeout"))
    if err == nil {
        // The request has a timeout, so create a context that is
        // canceled automatically when the timeout expires.
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
    } else {
        ctx, cancel = context.WithCancel(context.Background())
    }
    defer cancel() // Cancel ctx as soon as handleSearch returns.
}
```

并且使用WithValue函数传递客户端的IP：
```go
const userIPKey key = 0

// NewContext returns a new Context carrying userIP.
func NewContext(ctx context.Context, userIP net.IP) context.Context {
    return context.WithValue(ctx, userIPKey, userIP)
}
```

google包里面的Search函数实际的动作是将请求的参数传递给https://developers.google.com/custom-search，并且带上context的超时属性，当context超时的时候将会直接返回，不会等待https://developers.google.com/custom-search的返回。


# 参考

Go 并发控制context实现原理剖析：https://my.oschina.net/renhc/blog/2249581

golang通过context控制并发的应用场景实现：https://www.jb51.net/article/178003.htm

Golang 高效实践之并发实践context篇：https://www.cnblogs.com/makelu/p/11215530.html