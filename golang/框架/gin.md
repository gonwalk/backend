
# Gin的路由

Gin 框架的路由结构浅析：https://segmentfault.com/a/1190000016655709

Gin 是 go 语言的一款轻量级框架，风格简单朴素，支持中间件，动态路由等功能。gin项目github地址
路由是web框架的核心功能。在没有读过 gin 的代码之前，在我眼里的路由实现是这样的：根据路由里的 / 把路由切分成多个字符串数组，然后按照相同的前子数组把路由构造成树的结构；寻址时，先把请求的 url 按照 / 切分，然后遍历树进行寻址。

比如：定义了两个路由 /user/get，/user/delete，则会构造出拥有三个节点的路由树，根节点是 user，两个子节点分别是 get delete。

上述是一种实现路由树的方式，且比较直观，容易理解。对 url 进行切分、比较，时间复杂度是 O(2n)。

Gin的路由实现使用了类似前缀树的数据结构，只需遍历一遍字符串即可，时间复杂度为O(n)。

当然，对于一次 http 请求来说，这点路由寻址优化可以忽略不计。

Engine

Gin 的 Engine 结构体内嵌了 RouterGroup 结构体，定义了 GET，POST 等路由注册方法。

Engine 中的 trees 字段定义了路由逻辑。trees 是 methodTrees 类型（其实就是 []methodTree），trees 是一个数组，不同请求方法的路由在不同的树（methodTree）中。

最后，methodTree 中的 root 字段（*node类型）是路由树的根节点。树的构造与寻址都是在 *node的方法中完成的。

UML 结构图
engine结构图

trees 是个数组，数组里会有不同请求方法的路由树。

tree结构

node

node 结构体定义如下

type node struct {
    path      string           // 当前节点相对路径（与祖先节点的 path 拼接可得到完整路径）
    indices   string           // 所以孩子节点的path[0]组成的字符串
    children  []*node          // 孩子节点
    handlers  HandlersChain    // 当前节点的处理函数（包括中间件）
    priority  uint32           // 当前节点及子孙节点的实际路由数量
    nType     nodeType         // 节点类型
    maxParams uint8            // 子孙节点的最大参数数量
    wildChild bool             // 孩子节点是否有通配符（wildcard）
}
path 和 indices

关于 path 和 indices，其实是使用了前缀树的逻辑。

举个栗子：
如果我们有两个路由，分别是 /index，/inter，则根节点为 {path: "/in", indices: "dt"...}，两个子节点为{path: "dex", indices: ""}，{path: "ter", indices: ""}

handlers

handlers里存储了该节点对应路由下的所有处理函数，处理业务逻辑时是这样的：

func (c *Context) Next() {
    c.index++
    for s := int8(len(c.handlers)); c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}
一般来说，除了最后一个函数，前面的函数被称为中间件。

如果某个节点的 handlers为空，则说明该节点对应的路由不存在。比如上面定义的根节点对应的路由 /in 是不存在的，它的 handlers就是[]。

nType

Gin 中定义了四种节点类型：

const (
    static nodeType = iota // 普通节点，默认
    root       // 根节点
    param      // 参数路由，比如 /user/:id
    catchAll   // 匹配所有内容的路由，比如 /article/*key
)
param 与 catchAll 使用的区别就是 : 与 * 的区别。* 会把路由后面的所有内容赋值给参数 key；但 : 可以多次使用。
比如：/user/:id/:no 是合法的，但 /user/*id/:no 是非法的，因为 * 后面所有内容会赋值给参数 id。

wildChild

如果孩子节点是通配符（*或者:），则该字段为 true。

一个路由树的例子

定义路由如下：

r.GET("/", func(context *gin.Context) {})
r.GET("/index", func(context *gin.Context) {})
r.GET("/inter", func(context *gin.Context) {})
r.GET("/go", func(context *gin.Context) {})
r.GET("/game/:id/:k", func(context *gin.Context) {})
得到的路由树结构图为：
![路由树结构](https://segmentfault.com/img/bVbh22h?w=1640&h=1256)