Golang中Set类型的实现方法示例详解：https://www.jb51.net/article/124157.htm

# 需求


对于Set类型的数据结构，其实本质上跟List没什么多大的区别。无非是Set不能含有重复的Item的特性，Set有初始化、Add、Clear、Remove、Contains等操作。接下来看具体的实现方式分析吧。

# 实现


仍然按照已有的编程经验来联想如何实现基本Set功能，在Java中很容易知道HashSet的底层实现是HashMap，核心的就是用一个常量来填充Map键值对中的Value选项。除此之外，重点关注Go中Map的数据结构，Key是不允许重复的，如下所示：
 
m := map[string]string{
 "1": "one",
 "2": "two",
 "1": "one",
 "3": "three",
 }
 fmt.Println(m)
 

程序会直接报错，提示重复Key值，这样就非常符合Set的特性需求了。

## 定义


前面分析出Set的Value为固定的值，用一个常量替代即可。但是笔者分析的实现源码，用的是一个空结构体来实现的，如下所示：
```go 
// 空结构体
var Exists = struct{}{}
// Set is the main interface
type Set struct {
    // struct为结构体类型的变量
    m map[interface{}]struct{}
}
``` 

为了解决上面为什么用空结构体来做常量Value，先看下面的是测试：
```go
import (
 "fmt"
 "unsafe"
)
 
// 定义非空结构体
type S struct {
  a uint16
  b uint32
}
 
func main() {
 var s S
 fmt.Println(unsafe.Sizeof(s)) // prints 8, not 6
 var s2 struct{}
 fmt.Println(unsafe.Sizeof(s2)) // prints 0
}
``` 

打印出空结构体变量的内存占用大小为0，再看看下面这个测试：

a := struct{}{}
b := struct{}{}
fmt.Println(a == b) // true
fmt.Printf("%p, %p\n", &a, &b) // 0x55a988, 0x55a988
 

很有趣，a和b竟然相等，并且a和b的地址也是一样的。现在各位应该明白了为什么会有：
 
var Exists = struct{}{}
 

这样的常量也来填充所有Map的Value了吧，Go真是精彩！！！

## 初始化


Set类型数据结构的初始化操作，在声明的同时可以选择传入或者不传入进去。声明Map切片的时候，Key可以为任意类型的数据，用空接口来实现即可。Value的话按照上面的分析，用空结构体即可：

func New(items ...interface{}) *Set {
 // 获取Set的地址
 s := &Set{}
 // 声明map类型的数据结构
 s.m = make(map[interface{}]struct{})
 s.Add(items...)
 return s
}
 

## 添加


简化操作可以添加不定个数的元素进入到Set中，用变长参数的特性来实现这个需求即可，因为Map不允许Key值相同，所以不必有排重操作。同时将Value数值指定为空结构体类型。

func (s *Set) Add(items ...interface{}) error {
 for _, item := range items {
 s.m[item] = Exists
 }
 return nil
}
 

## 包含


Contains操作其实就是查询操作，看看有没有对应的Item存在，可以利用Map的特性来实现，但是由于不需要Value的数值，所以可以用 _,ok来达到目的：

func (s *Set) Contains(item interface{}) bool {
 _, ok := s.m[item]
 return ok
}
 

## 长度和清除


获取Set长度很简单，只需要获取底层实现的Map的长度即可：

func (s *Set) Size() int {
 return len(s.m)
}
 

清除操作的话，可以通过重新初始化Set来实现，如下即为实现过程：

func (s *Set) Clear() {
 s.m = make(map[interface{}]struct{})
}
 

## 相等


判断两个Set是否相等，可以通过循环遍历来实现，即将A中的每一个元素，查询在B中是否存在，只要有一个不存在，A和B就不相等，实现方式如下所示：
func (s *Set) Equal(other *Set) bool {
 // 如果两者Size不相等，就不用比较了
 if s.Size() != other.Size() {
 return false
 }
  
 // 迭代查询遍历
 for key := range s.m {
  // 只要有一个不存在就返回false
 if !other.Contains(key) {
 return false
 }
 }
 return true
}
 

## 子集


判断A是不是B的子集，也是循环遍历的过程，具体分析在上面已经讲述过，实现方式如下所示：
 
func (s *Set) IsSubset(other *Set) bool {
 // s的size长于other，不用说了
 if s.Size() > other.Size() {
 return false
 }
 // 迭代遍历
 for key := range s.m {
 if !other.Contains(key) {
 return false
 }
 }
 return true
}
 

Ok，以上就是Go中Set的主要函数实现方式，还是很有意思的。继续加油。
