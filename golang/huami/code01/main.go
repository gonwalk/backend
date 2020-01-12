package main

import (
	"fmt"
	"math/rand"
	"os"
)

/*
随机生成数字[1-200] 保存在文件中，每一行只能显示一个数字。
且随机生成的数字如果是偶数，则保持不变,奇数则写入该数的平方
*/

func main(){
	url := "D:/file/random.txt"
	file, err := os.Open(url)
	if err != nil {
		fmt.Println("打开文件失败:", err)
	}
	var intList []int
	for i := 0 ; i < 200; i++{
		curNum := rand.Intn(200)
		if curNum % 2 == 0{
			intList = append(intList, curNum)
		}else{
			intList = append(intList, curNum * curNum)
		}
	}
	writeNum(intList, url)
	defer file.Close()
	for _, num := range intList {
		fmt.Println(num)
	}
}

func writeNum(numList []int, path string){
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	var buffer string
	for _, num := range numList{
		buffer = fmt.Sprintf("%d", num)
		_, err := f.WriteString(buffer)
		if err != nil {
			fmt.Println("err=", err)
		}
		f.WriteString("\n")
	}
}