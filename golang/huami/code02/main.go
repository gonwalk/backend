package main

import (
	"fmt"
	"net/url"
)
/*
.实现一个方法，获取 Url 上的 Query 参数

例如： https://www.baidu.com/index.php?tn=98012088_3_dg&ch=1

GetQuery("tn")  => 98012088_3_dg

GetQuery()  => ["tn"=> "98012088_3_dg","ch"=>"1"]

*/

func main(){
	// var url string
	res := GetQuery()
	fmt.Println(res)
	//uri := "https://www.baidu.com/index.php?tn=98012088_3_dg&ch=1"
	//m := json.Unmarshal(uri)
	//fmt.Println(m)
	//
	//str := "a,b,c,d,e,f,g,h,i"
	//l := strings.Parse(str, ",")
	//fmt.Println(l)
}

func GetQuery()(map[string]string){
	m := make(map[string]string)
	uri := "https://www.baidu.com/index.php?tn=98012088_3_dg&ch=1"
	val, err := url.ParseRequestURI(uri)
	if err != nil {
		return m
	}
	
	for k, v := range val.Query(){
		fmt.Println(k, v)
		// for _, r := range v {
		// 	// urlList = append(urlList, r)
		// 	m[k] = r
		// }
		
		m[k] = v[0]
		
	}
	fmt.Println(m)

	return m
}