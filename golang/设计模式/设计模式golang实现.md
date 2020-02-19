Golang设计模式
https://www.jianshu.com/p/ea4d4d97b0c9

从大的方面来讲，设计模式有创建模式、结构模式、行为模式三类。

## 1.Golang单例模式实现
Golang 单例模式实现https://www.cnblogs.com/wpnine/p/10426105.html
单例模式在开发中是一种相对简单的设计模式，但它在实现上又有很多种方式

熟悉java的同学知道在java中实现单例常见的有懒汉式、饿汉式、双重检查、内部静态类、枚举单例等（[传送门](https://www.cnblogs.com/garryfu/p/7976546.html)）

而由于语言的特性，golang目前常见的有以下四种方式（懒汉式、饿汉式、双重检查、sync.Once）

### 1.1 懒汉式----非线程安全

非线程安全，即在多线程下可能会创建多次对象
```
 // 使用结构体代替类
type Tool struct {
    values int
}

// 建立私有变量
var instance *Tool

// 获取单例对象的方法，引用传递返回
func GetInstance() *Tool {
    if instance == nil {
        instance = new(Tool)
    }

    return instance
}
```
 

### 1.2 懒汉式----线程安全

在上面非线程安全的懒汉模式基础上，利用sync.Mutex进行加锁，保证线程安全，但由于每次调用该方法都进行了加锁操作，在性能上相对不高效

```
// 使用结构体代替类
type Tool struct {
    values int
}

// 建立私有变量
var instance *Tool

// 锁对象
var lock sync.Mutex

// 加锁保证线程安全
func GetInstance() *Tool {
    lock.Lock()
    defer lock.Unlock()
    if instance == nil {
        instance = new(Tool)
    }

    return instance
}
```
 

### 1.3 饿汉式

直接创建好对象，这样不需要判断为空，同时也是线程安全。唯一的缺点是在导入包的同时会创建该对象，并持续占有在内存中。

var instance Tool

func GetInstance() *Tool {
    return &instance
}
 

### 1.4 双重检查

在懒汉式（线程安全）的基础上再进行忧化，减少加锁的操作。保证线程安全同时不影响性能

// 锁对象
var lock sync.Mutex
var instance Tool

// 第一次判断不加锁，第二次加锁保证线程安全，一旦对象建立后，获取对象就不用加锁了

func GetInstance() *Tool {
    if instance == nil {
        lock.Lock()

        if instance == nil {
            instance = new(Tool)
        }

        lock.Unlock()
    }

    return instance
}

 

### 1.5 sync.Once

通过sync.Once 来确保创建对象的方法只执行一次


var once sync.Once
var instance Tool

func GetInstance() *Tool {
    once.Do(func() {
        instance = new(Tool)

    })
    return instance
}


sync.Once内部本质上也是双重检查的方式，但在写法上会比自己写双重检查更简洁，以下是Once的源码
```
func (o *Once) Do(f func()) {
　　　//判断是否执行过该方法，如果执行过则不执行
    if atomic.LoadUint32(&o.done) == 1 {
        return
    }
    // Slow-path.
    o.m.Lock()
    defer o.m.Unlock()
　　//进行加锁，再做一次判断，如果没有执行，则进行标志已经扫行并调用该方法
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```


# 2.观察者模式

设计模式有创建模式、结构模式、行为模式三类，观察者模式属于行为模式，下面是它的定义：

## 2.1 观察者模式定义
观察者模式(Observer): 定义对象间的一种一对多的依赖关系，以便当一个对象的状态发生改变时，所有依赖于它的对象都得到通知并自动更新。

## 2.2 观察者模式组成
观察者模式包含如下角色：
目标（Subject）: 目标知道它的观察者。可以有任意多个观察者观察同一个目标。 提供注册和删除观察者对象的接口。

具体目标（ConcreteSubject）:  将有关状态存入各ConcreteObserver具体观察者对象。

观察者(Observer): 为那些在目标发生改变时需获得通知的对象定义一个更新接口。当它的状态发生改变时, 向它的各个观察者发出通知。

具体观察者(ConcreteObserver): 维护一个指向ConcreteSubject具体对象的引用。存储有关状态，这些状态应与目标的状态保持一致。实现Observer的更新接口以使自身状态与目标的状态保持一致。

## 2.3 观察者模式效果

Observer模式允许独立地改变目标和观察者。可以单独复用目标对象而无需同时复用其观察者, 反之亦然。也可以在不改动目标和其他的观察者的前提下增加观察者。

下面是关于观察者模式的其他方面的一些优点:
```
1 )观察者模式可以实现表示层和数据逻辑层的分离，并定义了稳定的消息更新传递机制，抽象了更新接口，使得可以有各种各样不同的表示层作为具体观察者角色。

2 )在观察目标和观察者之间建立一个抽象的耦合：一个目标所知道的仅仅是它有一系列观察者, 每个都符合抽象的Observer类的简单接口。目标不知道任何一个观察者属于哪一个具体的类。这样目标和观察者之间的耦合是抽象的和最小的。因为目标和观察者不是紧密耦合的, 它们可以属于一个系统中的不同抽象层次。一个处于较低层次的目标对象可与一个处于较高层次的观察者通信并通知它, 这样就保持了系统层次的完整。如果目标和观察者混在一块, 那么得到的对象要么横贯两个层次 (违反了层次性), 要么必须放在这两层的某一层中(这可能会损害层次抽象)。

3) 支持广播通信：不像通常的请求, 目标发送的通知不需指定它的接收者。通知被自动广播给所有已向该目标对象登记的有关对象。目标对象并不关心到底有多少对象对自己感兴趣；它唯一的责任就是通知它的各观察者。这就使得可以在任何时刻自由地增加和删除观察者。处理还是忽略一个通知取决于观察者。

4) 观察者模式符合“开闭原则”的要求。
```
## 2.4 观察者模式代码示例

```go
package main
           
import (
    "container/list"
)
           
type Subject interface {        // 目标接口
    Attach(Observer) //注册Registe观察者 
    Detach(Observer) //释放Free观察者 
    Notify()         //通知Notify所有注册的观察者 
}

type Observer interface {       // 观察者接口
    Update(Subject) //观察者对目标进行更新状态 
}
           
//implements Subject
type ConcreteSubject struct {   // 具体目标结构
    observers *list.List        // 观察者对象列表
    value     int
}
           
func NewConcreteSubject() *ConcreteSubject {
    s := new(ConcreteSubject)
    s.observers = list.New()
    return s
}
           
func (s *ConcreteSubject) Attach(observe Observer) {     //注册观察者 
    s.observers.PushBack(observe)
}
           
func (s *ConcreteSubject) Detach(observer Observer) {   //释放观察者 
    for ob := s.observers.Front(); ob != nil; ob = ob.Next() {
        if ob.Value.(*Observer) == &observer {
            s.observers.Remove(ob)
            break
        }
    }
}
           
func (s *ConcreteSubject) Notify() { //通知所有观察者 
    for ob := s.observers.Front(); ob != nil; ob = ob.Next() {
        ob.Value.(Observer).Update(s)
    }
}
           
func (s *ConcreteSubject) setValue(value int) {
    s.value = value
    s.Notify()
}
           
func (s *ConcreteSubject) getValue() int {
    return s.value
}
           
/**
 * 具体观察者 implements Observer
 *
 */
type ConcreteObserver1 struct {
}
           
func (c *ConcreteObserver1) Update(subject Subject) {
    println("ConcreteObserver1  value is ", subject.(*ConcreteSubject).getValue())
}
           
/**
 * 具体观察者 implements Observer
 *
 */
type ConcreteObserver2 struct {
}
           
func (c *ConcreteObserver2) Update(subject Subject) {
    println("ConcreteObserver2 value is ", subject.(*ConcreteSubject).getValue())
}
           
func main() {
           
    subject := NewConcreteSubject()
    observer1 := new(ConcreteObserver1)
    observer2 := new(ConcreteObserver2)
    subject.Attach(observer1)
    subject.Attach(observer2)
    subject.setValue(5)          
}
```

运行结果：
ConcreteObserver1 vaue is  5
ConcreteObserver2 vaue is  5

# 3.装饰器模式

Golang装饰器设计模式（九）：https://blog.csdn.net/weixin_40165163/article/details/90740155
读代码之golang装饰器模式：https://studygolang.com/articles/24735?fr=sidebar

Go 装饰器模式教程：https://www.codercto.com/a/66476.html

装饰器设计模式允许向一个现有的对象添加新的功能，同时又不改变其结构。这种类型的设计模式属于结构型模式，它是作为现有的类的一个包装。

这种模式创建了一个装饰类，用来包装原有的类，并在保持类方法签名完整性的前提下，提供额外的功能。装饰器本质上允许包装现有功能并在开始或结尾处添加自己的自定义功能。

意图：动态地给一个对象添加一些额外的职责。就增加功能来说，装饰器模式相比生成子类的方式更为灵活。

主要解决：一般的，为了扩展一个类经常使用继承的方式实现，由于继承为类引入了静态特征，并且随着扩展功能的增加，子类会很庞大。

何时使用：在不想增加很多子类的情况下扩展类。

如何解决：将具体功能职责划分，同时继承装饰者模式。

关键代码：
```
Component 类充当抽象角色，不应该具体实现。
修饰类引用和继承 Component 类，具体扩展类重写父类方法。
优点：装饰结构和被装饰结构可以独立发展，不会相互耦合，装饰模式是继承的一个替代模式，装饰模式可以动态扩展一个实现类的功能。
```
缺点：多层装饰比较复杂

使用场景：
```
扩展一个类的功能
动态增加功能，动态撤销
```

实现：

下面的程序中，创建一个Shape接口和实现了 Shape 接口的结构体。然后创建一个实现了 Shape 接口的抽象装饰ShapeDecorator结构体，并把 Shape对象作为它的实例变量。

RedShapeDecorator 是实现了 ShapeDecorator 的实体。

DecoratorPatternDemo，演示类使RedShapeDecorator 来装饰 Shape 对象。

```go
package DecoratorPattern
 
import "fmt"
 
type Shape interface {
	Draw1()
}
 
type Rectangle struct {
}
 
func (r *Rectangle) Draw1() {
	fmt.Println("Shape: Rectangle")
}
 
type Circle struct {
}
 
func (c *Circle) Draw1() {
	fmt.Println("Shape: Circle")
}
 
type ShapeDecorator struct {
	decoratedShape Shape
}
 
func (s *ShapeDecorator) ShapeDecorator(decoratedShape Shape) {
	s.decoratedShape = decoratedShape
}
 
func (s *ShapeDecorator) Draw1() {
	s.decoratedShape.Draw1()
}
 
type RedShapeDecorator struct {
	shapeDecorator ShapeDecorator
}
 
func (s *RedShapeDecorator) RedShapeDecorator(decoratedShape Shape) {
	s.shapeDecorator.ShapeDecorator(decoratedShape)
}
 
func (s *RedShapeDecorator) Draw1() {
	s.shapeDecorator.Draw1()
	s.setRedBorder(s.shapeDecorator.decoratedShape)
}
 
func (s *RedShapeDecorator) setRedBorder(decoratedShape Shape) {
	fmt.Println("Border Color: Red")
}
```

原文链接：https://blog.csdn.net/weixin_40165163/article/details/90740155

# 4.工厂模式
Golang设计模式实现1-工厂模式https://www.cnblogs.com/ximen/p/9361284.html
工厂模式
工厂模式（Factory Pattern）是 Java 中最常用的设计模式之一。这种类型的设计模式属于创建型模式，它提供了一种创建对象的最佳方式。

在工厂模式中，我们在创建对象时不会对客户端暴露创建逻辑，并且是通过使用一个共同的接口来指向新创建的对象。

介绍
意图：定义一个创建对象的接口，让其子类自己决定实例化哪一个工厂类，工厂模式使其创建过程延迟到子类进行。

主要解决：主要解决接口选择的问题。

何时使用：我们明确地计划不同条件下创建不同实例时。

如何解决：让其子类实现工厂接口，返回的也是一个抽象的产品。

关键代码：创建过程在其子类执行。

应用实例： 1、您需要一辆汽车，可以直接从工厂里面提货，而不用去管这辆汽车是怎么做出来的，以及这个汽车里面的具体实现。 2、Hibernate 换数据库只需换方言和驱动就可以。

优点： 1、一个调用者想创建一个对象，只要知道其名称就可以了。 2、扩展性高，如果想增加一个产品，只要扩展一个工厂类就可以。 3、屏蔽产品的具体实现，调用者只关心产品的接口。

缺点：每次增加一个产品时，都需要增加一个具体类和对象实现工厂，使得系统中类的个数成倍增加，在一定程度上增加了系统的复杂度，同时也增加了系统具体类的依赖。这并不是什么好事。

使用场景： 1、日志记录器：记录可能记录到本地硬盘、系统事件、远程服务器等，用户可以选择记录日志到什么地方。 2、数据库访问，当用户不知道最后系统采用哪一类数据库，以及数据库可能有变化时。 3、设计一个连接服务器的框架，需要三个协议，"POP3"、"IMAP"、"HTTP"，可以把这三个作为产品类，共同实现一个接口。

注意事项：作为一种创建类模式，在任何需要生成复杂对象的地方，都可以使用工厂方法模式。有一点需要注意的地方就是复杂对象适合使用工厂模式，而简单对象，特别是只需要通过 new 就可以完成创建的对象，无需使用工厂模式。如果使用工厂模式，就需要引入一个工厂类，会增加系统的复杂度。

实现
我们将创建一个 Shape 接口和实现 Shape 接口的实体类。下一步是定义工厂类 ShapeFactory。

FactoryPatternDemo，我们的演示类使用 ShapeFactory 来获取 Shape 对象。它将向 ShapeFactory 传递信息（CIRCLE / RECTANGLE / SQUARE），以便获取它所需对象的类型。

工厂模式的 UML 图

步骤 1
创建一个接口

type Shape interface {
    Draw()
}
　

步骤 2
创建实现接口的实体类。


type Rectangle struct {
}
 
func (this Rectangle) Draw() {
    fmt.Println("Inside Rectangle::draw() method.")
}
 
 
 
type Square struct {
}
 
func (this Square) Draw() {
    fmt.Println("Inside Square ::draw() method.")
}
 
 
 
type Circle struct {
}
 
func (this Circle) Draw() {
    fmt.Println("Inside Circle  ::draw() method.")
}
　　

步骤 3
先给我一千亿，不、先创建一个工厂，哈哈，生成基于给定信息的实体类的对象。

type ShapeFactory struct {
}
 
//使用 getShape 方法获取形状类型的对象
func (this ShapeFactory) getShape(shapeType string) Shape {
 
    if shapeType == "" {
        return nil
    }
    if shapeType == "CIRCLE" {
        return Circle{}
    } else if shapeType == "RECTANGLE" {
        return Rectangle{}
    } else if shapeType == "SQUARE" {
        return Square{}
    }
    return nil
}
　　

步骤 4
使用该工厂，通过传递类型信息来获取实体类的对象。

func main() {
    factory := ShapeFactory{}
    factory.getShape("CIRCLE").Draw()
    factory.getShape("RECTANGLE").Draw()
    factory.getShape("SQUARE").Draw()
}
　　

步骤 5
执行程序，输出结果：


Inside Circle  ::draw() method.
Inside Rectangle::draw() method.