package main

import (
	"fmt"
	"sync"
)

var (
	ch = make(chan int, 1)
	wg sync.WaitGroup
)

// 100个并发按顺序打印1~10000
func orderGoroutine(k int) {
	fmt.Println(k)
	for i := k; i <= 10000; {
		//没到该 goroutine,就放回去
		temp := <-ch
		if temp != k {
			ch <- temp
			//fmt.Println(k)
		} else {
			fmt.Println(i)
			i += 100
			if k+1 == 101 {
				ch <- 1			// ch接收1~100的数字，即100个管道
			} else {
				ch <- (k + 1)
			}
		}
	}
	wg.Done()
}

// 100个并发随机打印1~10000个数字
// func unorderGoroutine(k int) {
// 	// fmt.Println(k)
// 	for i := k; i <= 10000; {
// 		//没到该 goroutine,就放回去
		
// 		// fmt.Println(i)
// 		// // i += 100
// 		// if k+1 == 101 {
// 		// 	ch <- 1			// ch接收1~100的数字，即100个管道
// 		// } else {
// 		// 	ch <- (k + 1)
// 		// }
		
// 	}
// 	wg.Done()
// }

func main() {

	wg.Add(100)						// 开启100个协程
	ch <- 1
	for i := 1; i <= 100; i++ {
		go orderGoroutine(i)
	}
	wg.Wait()
}
