package main

import (
	"fmt"
	"regexp"
)

/*
后端面试问题：

1、写一个判断版本号格式的方法，版本号分3部分，样例为：a1.1.1,第一部分为字母加数字，
后面两部分为纯数字（支持扩展性，支持规定数字的范围）
2、我们有一个业务是打卡团分数的排名（数据量大概在50万，分数来源可以不用管，但是实时在更新）
 a)排行榜中，按名次从小到大排列
 b)分数相同的团的名次是一样的，例如第一名有3个，第二名有10个（分数相同的团排序先后顺序可以无所谓）
 c)能快速查找某个团的名次
 设计一个满足上面需求的后台设计及接口
*/

func main() {
	a1 := "a1.1.1"
	isValid1 := judgeVersion(a1)
	fmt.Println(isValid1)
	a2 := "1.23.56"
	isValid2 := judgeVersion(a2)
	fmt.Println(isValid2)

	a3 := "c1.12.19"
	
}

func judgeVersion(str string) bool {
	matchStr := `^([a-zA-Z][0-9]+).([0-9])+.([0-9])`
	var isValidVersion = regexp.MustCompile(matchStr)
	return isValidVersion.MatchString(str)
}


func judgeVersionPatch(str string , start1, end1, start2, end2 int) bool {
	matchStr := `^([a-zA-Z][0-9]+).` + `([` + start1 + `-` +  end1 +`])+.([`0-9`])`
	var isValidVersion = regexp.MustCompile(matchStr)
	return isValidVersion.MatchString(str)
}