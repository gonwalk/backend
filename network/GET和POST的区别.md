# GET和POST的区别

GET和POST是HTTP请求的两种基本方法，两者最直观的区别是：
```
GET把参数包含在URL中，而POST通过请求体request body传递参数。
```
 
## GET和POST使用场景

你可能自己写过无数个GET和POST请求，它们直接具体的区别在哪里，什么时候该用GET和POST呢？具体可以归纳为以下几点：

```
GET参数通过URL传递，POST放在Request body中。
GET请求只能进行URL编码，而POST支持多种编码方式。
GET请求提交的数据有长度限制，大部分浏览器都是限制在2KB以内，而POST没有。
对参数的数据类型，GET只接受ASCII字符，而POST没有限制。
GET和POST底层都是TCP连接，但是GET产生一个TCP数据包；POST产生两个TCP数据包。

GET请求参数会被完整保留在浏览器历史记录里，而POST中的参数不会被保留。
GET请求会被浏览器主动cache，而POST不会，除非手动设置。
GET比POST更不安全，因为参数直接暴露在URL上，所以不能用来传递敏感信息。

GET在浏览器回退时是无害的，而POST会再次提交请求。
GET产生的URL地址可以被Bookmark，而POST不可以。
```

（本标准答案参考自w3schools）

“很遗憾，这不是我们要的回答！”，请告诉我真相。。。

## 本质区别 

如果我告诉你GET和POST本质上没有区别你信吗？


让我们扒下GET和POST的外衣，坦诚相见吧！

GET和POST是什么？GET和POST是HTTP协议中的两种发送请求的方法。

HTTP是什么？HTTP是基于TCP/IP的关于数据在万维网中如何进行通信的协议。

### GET和POST底层都是TCP连接

HTTP的底层是TCP/IP协议簇。所以GET和POST的底层也是TCP/IP，也就是说，GET/POST都是TCP连接。GET和POST能做的事情是一样一样的。你要给GET加上request body，给POST带上url参数，技术上是完全行得通的。

那么，“标准答案”里的那些区别是怎么回事？

在万维网世界中，TCP就像汽车，使用其TCP来运输数据，它很可靠，从来不会发生丢件少件的现象。但是如果路上跑的全是看起来一模一样的汽车，那这个世界看起来是一团混乱，送急件的汽车可能被前面满载货物的汽车拦堵在路上，整个交通系统一定会瘫痪。为了避免这种情况发生，交通规则HTTP诞生了。HTTP给汽车运输设定了好几个服务类别，有GET, POST, PUT, DELETE等等。

HTTP规定，当执行GET请求的时候，要给汽车贴上GET的标签（设置method为GET），而且要求把传送的数据放在车顶上（url中）以方便记录。如果是POST请求，就要在车上贴上POST的标签，并把货物放在车厢里。当然，也可以在使用GET的时候往车厢内偷偷藏点货物，但是这是很不光彩；也可以在POST的时候在车顶上也放一些数据，让人觉得傻乎乎的。HTTP只是个行为准则，而TCP才是GET和POST怎么实现的基本。

直观上，HTTP对GET和POST参数的传送渠道（url还是requrest body）提出了要求。上面“标准答案”里关于参数大小的限制又是从哪来的呢？

在万维网世界中，还有另一个重要的角色：运输公司。不同的浏览器（发起http请求）和服务器（接受http请求）就是不同的运输公司。 虽然理论上，你可以在车顶上无限的堆货物（url中无限加参数）。但是运输公司可不傻，装货和卸货也是有很大成本的，他们会限制单次运输量来控制风险，数据量太大对浏览器和服务器都是很大负担。业界不成文的规定是，（大多数）浏览器通常都会限制url长度在2K个字节，而（大多数）服务器最多处理64K大小的url。超过的部分，恕不处理。如果你用GET服务，在request body偷偷藏了数据，不同服务器的处理方式也是不同的，有些服务器会帮你卸货，读出数据，有些服务器直接忽略，所以，虽然GET可以带request body，也不能保证一定能被接收到哦。

### 本质区别是TCP连接：GET发送一次数据包、POST发送两次

好了，现在你知道，GET和POST本质上就是TCP连接，并无差别。但是由于HTTP的规定和浏览器/服务器的限制，导致它们在应用过程中体现出一些不同。



到这里真正的大BOSS才出现，这位BOSS有多神秘？当你试图在网上找“GET和POST的区别”的时候，那些你会看到的搜索结果里，从没有提到它。它究竟是什么呢？————数据包

GET和POST还有一个重大区别，简单的说：

GET产生一个TCP数据包；POST产生两个TCP数据包。

具体来说：
```
对于GET方式的请求，浏览器会把http header和data一并发送出去，服务器响应200（返回数据）；

而对于POST，浏览器先发送header，服务器响应100 continue，浏览器再发送data，服务器响应200 ok（返回数据）。
```
也就是说，GET只需要汽车跑一趟就把货送到了，而POST得跑两趟，第一趟，先去和服务器打个招呼“嗨，我等下要送一批货来，你们打开门迎接我”，然后再回头把货送过去。

因为POST需要两步，时间上消耗的要多一点，看起来GET比POST更有效。因此Yahoo团队有推荐用GET替换POST来优化网站性能。但这是一个坑！跳入需谨慎。为什么？
```
1. GET与POST都有自己的语义，不能随便混用。

2. 据研究，在网络环境好的情况下，发一次包的时间和发两次包的时间差别基本可以无视。而在网络环境差的情况下，两次包的TCP在验证数据包完整性上，有非常大的优点。

3. 并不是所有浏览器都会在POST中发送两次包，Firefox就只发送一次。 
```
现在，当面试官再问你“GET与POST的区别”的时候，你的内心是不是这样的？


## 差异分析

### 一、原理区别

GET不会修改信息，POST可能会修改服务器上的资源。

一般我们在浏览器输入一个网址访问网站都是GET请求；在FORM表单中，可以通过设置Method指定提交方式为GET或者POST提交方式，默认为GET提交方式。

HTTP定义了与服务器交互的不同方法，最基本的五种为：GET，POST，PUT，DELETE，HEAD，其中GET和HEAD被称为安全方法，因为使用GET和HEAD的HTTP请求不会产生什么动作。不会产生动作意味着GET和HEAD的HTTP请求不会在服务器上产生任何修改的结果。但是安全方法并不是什么动作都不产生，这里的安全方法仅仅指不会修改信息。

根据HTTP规范，POST是一种可能会修改服务器上的资源的请求方法。比如CSDN的博客，用户提交一篇文章或者一个读者提交评论是通过POST请求来实现的，因为再提交文章或者评论提交后资源（即某个页面）不同了，或者说资源被修改了，这些便是“不安全方法”。

### 二、表现形式区别

搞清楚了两者的原理区别后，我们来看一下在实际应用中的区别。

首先，我们先看一下HTTP请求的格式：
```
<method> <request-URL> <version>
<headers>

<entity-body>
```

在HTTP请求中，其第一行必须是一个请求行，包括请求方法、请求URL、报文所用HTTP版本信息。紧接着是一个headers头部信息，可以有零个或一个首部，用来说明服务器要使用的附加信息。在首部之后就是一个空行，最后就是报文实体的主体部分，包含一个由任意数据组成的数据块。但是并不是所有的报文都包含实体的主体部分。

GET请求实例：
```
GET http://weibo.com/signup/signup.php?inviteCode=2388493434
Host: weibo.com
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
```

POST请求实例：
```
POST /inventory-check.cgi HTTP/1.1
Host: www.joes-hardware.com
Content-Type: text/plain
Content-length: 18

item=bandsaw 2647
```

接下来看看两种请求方式的区别：

1、数据的编码及存放位置

GET请求：请求的数据会附加在URL之后，以?分割URL和传输数据，多个参数用&连接。URL的编码方式采用的是ASCII编码，而不是unicode，即要求所有的非ASCII字符在都编码之后再传输。

POST请求：POST请求会把请求的数据放置在HTTP请求包的包体中。在上面的POST请求实例中，item=bandsaw就是实际的传输数据。

因此，GET请求的数据会暴露在地址栏中（不安全），而POST请求则不会。

2、传输数据的大小

在HTTP规范中，没有对URL的长度和传输的数据大小进行限制。但是在实际开发过程中，对于GET，特定的浏览器和服务器对URL的长度一般有2KB数据大小的限制。因此，在使用GET请求时，传输数据会受到URL长度的限制。

对于POST，由于不是URL传值，理论上是不会受限制的，但是实际上各个服务器会规定对POST提交数据大小进行限制，Apache、IIS都有各自的配置进行修改。

3、安全性

POST的安全性比GET的高。这里的安全是指真正的安全，而不同于GET安全方法中的安全，上面提到的安全方法中的安全仅仅是不修改服务器的数据。比如，在进行登录操作时，通过GET请求，用户名和密码都会暴露在URL上，因为登录页面有可能被浏览器缓存，其他人查看浏览器的历史记录时（通过书签bookmarks）就可以看到，此时的用户名和密码就很容易被泄露。除此之外，GET请求提交的数据还可能会造成Cross-site request frogery攻击

4、HTTP中的GET，POST，SOAP协议都是在HTTP上运行的。

### 三、HTTP响应

HTTP响应报文的格式
```
<version> <status> <reason-phrase>
<headers>

<entity-body>
```

status，状态码描述了请求过程中发生的情况

reason-phrase 是数字状态码的可读版本

常见的状态码以及含义如下：
```
200 OK 服务器成功处理请求

301/302 Moved Permanently（重定向）请求的URL已移走。响应报文中应该包含一个Location URL，说明资源现在所处的位置

304 Not Modified（未修改） 客户的缓存资源是最新的，要客户端使用缓存内容

404 Not Found 未找到资源

501 Internal Server Error 服务器遇到错误，使其无法对请求提供服务
```

HTTP响应示例：
```
HTTP/1.1 200 OK

Content-type: text/plain
Content-length: 12

Hello World!
```

## 参考资料

Get与Post的区别？（面试官最想听到的答案）：
https://blog.csdn.net/ever_siyan/article/details/87935455

HTTP请求中POST与GET的区别：https://blog.csdn.net/yipiankongbai/article/details/24025633