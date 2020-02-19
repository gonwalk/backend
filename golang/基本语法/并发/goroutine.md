# goroutine与并发

Go语言的并发通过goroutine特性完成。goroutine类似于线程，但是可以根据需要创建多个goroutine并发工作。
goroutine由Go语言的运行时（runtime）调度完成，而线程是由操作系统调度完成。

Go语言还提供channel在多个goroutine间进行通信。goroutine和channel是Go语言秉承的CSP（Communicating Sequential Process）并发模式的重要实现基础。

## 轻量级线程（goroutine）————根据需要随时创建的“线程”

Go语言中的goroutine机制：可以为使用者分配足够多的任务，系统能自动帮助使用者把任务分配到CPU上，让这些任务尽量并发运行。

goroutine的概念类似于线程，但goroutine由Go程序运行时（runtime）进行调度和管理。Go程序会智能地将goroutine中的任务合理地分配给每个CPU。
Go程序冲main包的main()函数开始，在程序启动时，Go程序就会为main()函数创建一个默认的goroutine。

### 创建goroutine方式

（1）使用普通函数创建goroutine

Go程序中使用go关键字为一个函数创建一个goroutine。一个函数可以创建多个goroutine调用，一个goroutine必定对应一个函数。

普通函数创建goroutine的格式：
```
go 函数名(参数列表)
```
使用go关键字创建goroutine时，被调用函数的返回值会被忽略。如果需要再goroutine中返回数据，可以借助通道（channel）特性，通过通道把数据从goroutine中作为返回值传出。

（2）使用匿名函数创建goroutine
注意：使用匿名函数的活闭包的方式创建goroutine时，除了将函数定义部分写在go的后面之外，还需要加上匿名函数的调用参数。
格式定义如下：
```go
go func(参数列表){      // 参数列表：函数体内的参数变量列表
    函数体              // 函数体：匿名函数内部的代码块
}(调用参数列表)         // 调用参数列表：启动goroutine时，需要向匿名函数传递的调用参数
```

在main()函数中创建一个匿名函数并为匿名函数启动goroutine。匿名函数没有参数，代码将并行执行定时打印计数的效果。




## 多并发示例

### Golang 协程顺序打印
https://blog.csdn.net/qianghaohao/article/details/97007270

（1）两个协程交替打印两个队列

A、B 两个协程分别打印 1、2、3、4 和 A，B，C，D
实现：定义 A、B 两个 channal，开 A、B 两个协程，A 协程输出[1, 2, 3, 4]、B 协程输出[A, B, C, D]，通过两个独立的 channal 控制顺序，交替输出。
```
func main() {
	A := make(chan bool, 1)
	B := make(chan bool)
	Exit := make(chan bool)
	go func() {
		s := []int{1, 2, 3, 4}
		for i := 0; i < len(s); i++ {
			if ok := <-A; ok {
				fmt.Println("A: ", s[i])
				B <- true
			}
		}
	}()
	go func() {
		defer func() {
			close(Exit)
		}()
		s := []byte{'A', 'B', 'C', 'D'}
		for i := 0; i < len(s); i++ {
			if ok := <-B; ok {
				fmt.Printf("B: %c\n", s[i])
				A <- true
			}
		}
	}()
	A <- true
	<-Exit
}
```

### A、B 两个协程顺序打印 1~20

实现：与上面基本一样，定义 A、B 两个 channal，开 A、B 两个协程，A 协程输出奇数、B 协程输出偶数，通过两个独立的 channal 控制顺序，交替输出。
```go
package main

import "fmt"

func main() {
	A := make(chan bool, 1)
	B := make(chan bool)
	Exit := make(chan bool)

	go func() {
		for i := 1; i <= 10; i++ {
			if ok := <-A; ok {
				fmt.Println("A = ", 2*i-1)
				B <- true
			}
		}
	}()
	go func() {
		defer func() {
			close(Exit)
		}()
		for i := 1; i <= 10; i++ {
			if ok := <-B; ok {
				fmt.Println("B : ", 2*i)
				A <- true
			}
		}
	}()

	A <- true
	<-Exit
}
```

### 多个协程顺序打印数字

原文链接：https://blog.csdn.net/jujueduoluo/article/details/80500450

```go
package main

import (
	"sync"
	"fmt"
	"time"
)

var (
	switchFlow chan int
	wg sync.WaitGroup
)

func routine(i int, serialNumber int) {
	time.Sleep(100 * time.Millisecond)
	loop:
	for {
		select {
		case s := <- switchFlow:
			if s == serialNumber  {
				fmt.Println("routine: ", i+1)
				break loop
			} else {
				//fmt.Println("接受到的编号是: ", s)
				switchFlow <- s
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	wg.Done()
	switchFlow <- serialNumber+1
}

func main() {
	switchFlow = make(chan int)
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go routine(i, 10+i)
	}
	//引爆点
	switchFlow <- 10

	wg.Wait()
	close(switchFlow)
	fmt.Println("程序结束")
}
```

