# 1. map

golang中的map是一种数据类型，将键与值绑定到一起，底层是用哈希表实现的，可以快速的通过键找到对应的值。

类型表示：map[keyType][valueType] key一定要是可比较的类型（可以理解为支持==的操作），value可以是任意类型。

初始化：map只能使用make来初始化，声明的时候默认为一个为nil的map，此时进行取值，返回的是对应类型的零值（不存在也是返回零值）。添加元素无任何意义，还会导致运行时错误。向未初始化的map赋值引起 panic: assign to entry in nil map。

map的声明和初始化的方式：
```go
m1 := make(map[string]string)       // 先用make声明，再赋值初始化
m1["phone"] = "123"
m1["add"] = "beijing"
m1["age"] = "12"

// 声明与初始化一起进行
m := map[byte][]string{
    '2':"abc",
    '3':"def",
    '4':"ghi",
    '5':"jkl",
    '6':"mno",
    '7':"pqrs",
    '8':"tuv",
    '9':"wxyz",
}
```

清空map：对于一个有一定数据的集合 exp，清空的办法就是再次初始化: exp = make(map[string]int)，如果后期不再使用该map，则可以直接：exp= nil 即可，但是如果还需要重复使用，则必须进行make初始化，否则无法为nil的map添加任何内容。

属性：与切片一样，map 是引用类型。当一个 map 赋值给一个新的变量，它们都指向同一个内部数据结构。因此改变其中一个也会反映到另一个。作为形参或返回参数的时候，传递的是地址的拷贝，扩容时也不会改变这个地址，因此改变其中一个值，另一个指向该变量的字典的值也会相应地发生改变。
```go
func main() {
    exp := map[string]int{
        "steve": 20,
        "jamie": 80,
    }
    fmt.Println("Ori exp", exp)
    newexp:= exp
    newexp["steve"] = 18
    fmt.Println("exp changed", exp)
}

//Ori age map[steve:20 jamie:80]
//age changed map[steve:18 jamie:80]
```

## 1.1 map底层数据结构

go中的数据结构-字典map：https://www.cnblogs.com/33debug/p/11851585.html

Go中的map在可以在 $GOROOT/src/runtime/map.go找到它的实现。哈希表的数据结构中一些关键的域如下所示：

### map底部三层结构

map的底层主要是由三个结构构成:
```
hmap --- map的最外层的数据结构，包括了map的各种基础信息，如：大小count、buckets数组指针等。
mapextra --- 记录map的额外信息，hmap结构体里的extra指针指向的结构，例如overflow bucket。
bmap --- 代表bucket，每一个bucket最多放8个kv键值对，最后由一个overflow字段指向下一个bmap，注意key、value、overflow字段都不显示定义，而是通过maptype计算偏移获取的。
```

其中hmap.extra.nextOverflow指向的是预分配的overflow bucket，预分配的用完了那么值就变成nil。

#### hmap结构

```go
type hmap struct {
    count        int  //元素个数
    flags        uint8   
    B            uint8 //扩容常量
    noverflow    uint16 //溢出 bucket 个数
    hash0        uint32 //hash 种子
    buckets      unsafe.Pointer //bucket 数组指针
    oldbuckets   unsafe.Pointer //扩容时旧的buckets 数组指针
    nevacuate    uintptr  //扩容搬迁进度
    extra        *mapextra //记录溢出相关
}

type bmap struct {
    tophash        [bucketCnt]uint8  
    // Followed by bucketCnt keys 
    //and then bucketan Cnt values  
    // Followed by overflow pointer.
}
```
说明：每个map的底层都是hmap结构体，它是由若干个描述hmap结构体的元素、数组指针、extra等组成，buckets数组指针指向由若干个bucket组成的数组，其每个bucket桶里存放的是key-value数据(通常是8个)和overflow字段（指向下一个bmap），每个key插入时会根据hash算法归到同一个bucket中，当一个bucket中的元素超过8个（key-value键值对）的时候，hmap会使用extra中的overflow来扩展存储key。

![go中map的数据结构](https://img2018.cnblogs.com/blog/1069650/201911/1069650-20191114174802127-585623786.png)

图中len 就是当前map的元素个数，也就是len()返回的值，也是结构体中hmap.count的值。bucket array是指数组指针，指向bucket数组；hash seed是哈希种子；overflow指向下一个bucket。

#### bmap结构

bmap的详细结构如下
![bmap结构](https://img2018.cnblogs.com/blog/1069650/201911/1069650-20191114203834314-657349672.png)


在map中出现哈希冲突时，首先以bmap为最小粒度挂载，一个bmap累积8个kv之后，就会申请一个新的bmap（overflow bucket）挂在这个bmap的后面形成链表，优先用预分配的overflow bucket，如果预分配的用完了，那么就malloc一个挂上去。这样减少对象数量，减轻管理内存的负担，利于gc。注意golang的map不会shrink，内存只会越用越多，overflow bucket中的key全删了也不会释放。

bmap中所有key存在一块，所有value存在一块，这样做方便内存对齐。当key大于128字节时，bucket的key字段存储的会是指针，指向key的实际内容；value也是一样。

hash值的高8位存储在bucket中的tophash字段。每个桶最多放8个kv对，所以tophash类型是数组[8]uint8。把高八位存储起来，这样不用完整比较key就能过滤掉不符合的key，加快查询速度。实际上当hash值的高八位小于常量minTopHash时，会加上minTopHash，区间[0, minTophash)的值用于特殊标记。查找key时，计算hash值，用hash值的高八位在tophash中查找，有tophash相等的，再去比较key值是否相同。



## 1.2 map的底层实现及哈希冲突

### map的底层实现

golang中的map采用了HashTable的实现，通过数组+链表实现的。一个哈希表会有一定数量的桶，哈希表将键值对均匀存储到这些桶（使用一个指向数组的指针表示）中。哈希表在存储键值对时，会先用哈希函数把键值转换为哈希值，哈希表先用哈希值的低几位去定位到一个哈希桶，然后再去这个哈希桶中查找这个键。由于键值对总是被捆绑在一起存在，一旦找到了键，就找到了值。go的字典中，每一个键值对都是它的哈希值代表的，字典不会独立存储任何键的键值，但会独立存储他们的哈希值。

### 哈希冲突解决方式
Hash算法解决冲突的四种方法：
https://www.cnblogs.com/lyfstorm/p/11044468.html

Hash算法解决冲突的方法一般有以下几种常用的解决方法：

#### 1.开放定址法

所谓的开放定址法就是一旦发生了冲突，就去寻找下一个空的散列地址，只要散列表足够大，空的散列地址总能找到，并将记录存入 
公式为：fi(key) = (f(key)+di) MOD m (di=1,2,3,……,m-1) 
※ 用开放定址法解决冲突的做法是：当冲突发生时，使用某种探测技术在散列表中形成一个探测序列。沿此序列逐个单元地查找，直到找到给定的关键字，或者 
碰到一个开放的地址（即该地址单元为空）为止（若要插入，在探查到开放的地址，则可将待插入的新结点存人该地址单元）。查找时探测到开放的地址则表明表 
中无待查的关键字，即查找失败。 
比如说，我们的关键字集合为{12,67,56,16,25,37,22,29,15,47,48,34},表长为12。 我们用散列函数f(key) = key mod l2 
当计算前S个数{12,67,56,16,25}时，都是没有冲突的散列地址，直接存入： 


计算key = 37时，发现f(37) = 1，此时就与25所在的位置冲突。 
于是我们应用上面的公式f(37) = (f(37)+1) mod 12 = 2。于是将37存入下标为2的位置


#### 2.再哈希法

再哈希法又叫双哈希法，有多个不同的Hash函数，当发生冲突时，使用第二个，第三个，….等多个哈希函数去计算地址，直到无冲突。虽然不易发生聚集，但是增加了计算时间。

#### 3.链地址法

链地址法的基本思想是：每个哈希表节点都有一个next指针，多个哈希表节点可以用next指针构成一个单向链表，被分配到同一个索引(数组的索引，map的结构使用数组+链表的方式)上的多个节点可以用这个单向链表连接起来，如： 
键值对k2, v2与键值对k1, v1通过计算后的索引值都为2，这时就会产生冲突，但是可以通过next指针将k2, k1所在的节点连接起来，这样就解决了哈希的冲突问题。 

#### 4.建立公共溢出区

这种方法的基本思想是：将哈希表分为基本表和溢出表两部分，凡是和基本表发生冲突的元素，一律填入溢出表。

## 1.3 map的键类型不能是哪些类型

不能是函数类型、字典类型、切片类型。
原因：go语言规范规定，键类型的值之间必须可以施加==操作符和!=操作符，即必须支持判等操作，由于函数类型，字典类型，切片类型不支持判等操作，所以map不支持这些类型。
```
如果键类型是接口类型，那么键值的实际类型也不能是上述3种，否则会引发panic。
如果键的类型是数组类型，那么还要确保该类型的元素类型不是函数类型、字典类型或切片类型。
如果键的类型是结构体类型，那么还要保证其中字段类型的合法性。 不论不合法的类型被埋藏得有多深，比如map[[1][2][3][]string]int，Go 语言编译器都会把它揪出来。
```

## 1.4 map对键类型有何要求

求哈希和判等操作比较快的类型适合做键类型。
不建议使用高级数据类型作为类型的原因不仅仅因为求哈希以及判等速度比较慢，而且它们的值存在变数。

## 1.5 在值为nil的字典上进行读写操作时会发生什么？

除了添加键值对，我们在一个值为nil的字典上做任何操作都不会引起错误。而一旦试图在一个值为nil的字典中添加键值对时，会抛出一个panic。


## 1.6 map是并发安全的吗

map类型的值不是并发安全的，即使只是添加或删除操作，也是不安全的，根本原因在于字典值内部有时候会根据需要进行存储方面的调整

# 2. sync.map

通过以上几个方面，我们了解了map底层实现，并且知道map不是并发安全的，那么golang有没有提供并发安全的map呢？答案是yes，那就是sync.map。Go官方在2017年发布的Go1.9中，正式加入了并发安全的字典类型sync.map，该map有以下几个特点：
```
（1）并发安全，且虽然用到了锁，但是显著减少了锁的争用。
sync.map出现之前，如果想要实现并发安全的map，只能自行构建，使用sync.Mutex或sync.RWMutex，再加上原生的map就可以轻松做到。sync.map也用到了锁，但是在尽可能的避免使用锁，因为使用锁意味着要把一些并行化的东西串行化，会降低程序性能，因此能用原子操作就不要用锁，但是原子操作局限性比较大，只能对一些基本的类型提供支持，在sync.map中将两者做了比较完美的结合。

（2）存取删操作的算法复杂度与map一样，都是O(1)。

（3）不会做类型检查。
sync.map只是go语言标准库中的一员，而不是语言层面的东西，也正是因为这一点，go语言的编译器不会对其中的键和值进行特殊的类型检查
```

下面我们通过分析sync.map中的源码，来看上面第1点和第2点是如何做到的

结构体
结构体包含Map、readOnly、entry。
其中Map中包含了两个map:
(1)一个优先读map：read，read不是只允许只读，它是可以有写操作的，但是只允许对已存在的值进行写操作，不允许增加新元素，删除也只是做标记，并不是真正的删除。即键的集合不能被改变，所以键值是不全的
(2)一个需要加锁进行操作的读写map：dirty，其中的键值对集合总是完全的，而且不包含已被逻辑删除的键值对
这两个map是实现并发安全的关键所在
```
type Map struct{
       m  Mutex         //互斥锁，用于对dirty进行加锁
       read atomic.Value  //优先读map，支持原子操作，并不是只读map，可以有写操作
       dirty map[interface{}]*entry  //当前最新map，需要加锁操作，允许读写，
       misses int     //记录read读取不到数据，需要加锁读取的次数，当misses等于dirty的长度时，会将dirty复制到read
}
```

readOnly是存储在read中的元素
```
type readOnly struct{
       m   map[interface{}]*entry//  该map中存储的entry指针，与dirty中存储的entry指针一样，因此虽然引入的两个map，但是底层存的是指针
       amended bool    //amended为false时，代表dirty中还没有数据，如果dirty中存在一些在read中没有的数据，则该值为true，可以看作修改的标记
}
```
entry存放着真正的实体数据
```
type entry struct{
       p unsafe.Pointer  
//其中p有如下取值：
//（1）如果p==nil，read中元素已经被删除， m.dirty == nil
//（2）如果p==expunged：表示read中的元素已经被删除，这种情况出现在将read复制到dirty中，即复制的过程会先将nil标记为expunged，然后不将其复制到dirty
//（3）否则，元素是合法的，被记录到m.read.m[key]中，如果m.dirty!=nil，那么值也在m.dirty[key]中
// 一个元素可以用nil通过原子替换的方式进行删除，当下一次创建m.drity，会自动用expunged替换nil，不会将其复制到dirty中
}
```

主要方法如下：
```
Load（key interface{}）（value，ok bool）{}
读操作：给定key，Load返回key对应的value，如果没有，返回nil
Store(key,value interface{}){}
存储键值对（key，value）
LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
如果key存在，则返回已存在的value，否则，存储(key，value)，loaded=true代表key原来不存在，存储了新值，loaded=false代表没有加载，读到的是之前存储的值
Delete(key interface{}) {}
删除key
Range(f func(key, value interface{}) bool){}
遍历map中的元素并调用f函数
```
下面我们一次看下这几个方法实现的源码

Load方法：
```go
//首先总结下过成：每次Load都先从read读取，当read中不存在且amended为true，就从dirty读取数据 。无论
//dirty map中是否存在该元素，都会执行missLocked函数，该函数将misses+1，
//当m.misses < len(m.dirty)时，便会将dirty复制到read，此时再将dirty置为nil,misses=0。
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	read, _ := m.read.Load().(readOnly)    
	e, ok := read.m[key]             //先从read字典中读
	if !ok && read.amended {//如果read中没有，且read.amended==true(该值表示值在dirty中，不在read中) 则需要加锁，然后去dirty中查找
		m.mu.Lock()
		//此处有一个double check操作，防止在加锁过程中，dirty转换成read，从而读不到数据
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 无论查找元素是否存在，都要记录miss值，以dirty升级为read
			m.missLocked()
		}
		m.mu.Unlock()
	}
    //  元素不存在直接返回
	if !ok {
		return nil, false
	}
	return e.load()
}
//记录miss值，当miss值大于等于dirty长度时，dirty升级为read，dirty置为nil，miss置为0
func (m *Map) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnly{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}
func (e *entry) load() (value interface{}, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expunged {
		return nil, false
	}
	return *(*interface{})(p), true
}
```

总结：通过Load方法可知，在进行读操作时，首先去read中去查找，这样避免了读写冲突，只有在read中取不到值时，才会加锁，去dirty中进行查找，并且会进行动态调整，即miss次数多了之后，会将dirty升级为read，非常适合热点数据的查询

Store方法：写方法
```go
func (m *Map) Store(key, value interface{}) {
    //如果read中存在这个键，且这个键没有被标记删除，直接在read中写入
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}
//read中不存在，需要加锁，到dirty中写入
	m.mu.Lock()
//doublecheck，同Load
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
// 之前读read读不到，而doublecheck时发现read中可以读到非nil的值，意味着什么？read不允许增加删除值，所以只能说明加锁之前有dirty升级为read的操作
//如果读到的值为expunged，说明生成dirty时，复制read中的元素，对于nil的元素，搞成了expunged，所以意味着dirty不为nil，且这个元素中没有该元素
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
        // 更新read中的值
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {// 此时，read中没有该元素，而dirty中有该元素，需要修改dirty中的值为最新值
		e.storeLocked(&value)
	} else {
		if !read.amended {
			// read.amended==false,说明dirty map为空，需要将read map 复制一份到dirty map
			m.dirtyLocked()
                       // 设置read.amended==true，说明dirty map有数据
			m.read.Store(readOnly{m: read.m, amended: true})
		}
               // 设置元素进入dirty map，此时dirty map拥有read map和最新设置的元素
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}
func (e *entry) tryStore(i *interface{}) bool {
    // 获取对应Key的元素，判断是否标识为删除
    p := atomic.LoadPointer(&e.p)
    if p == expunged {
        return false
    }
    for {
        // cas尝试写入新元素值
        if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
            return true
        }
        // 判断是否标识为删除
        p = atomic.LoadPointer(&e.p)
        if p == expunged {
            return false
        }
    }
}
func (e *entry) unexpungeLocked() (wasExpunged bool) {
    return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}
func (m *Map) dirtyLocked() {
    if m.dirty != nil {
        return
    }
    read, _ := m.read.Load().(readOnly)
    m.dirty = make(map[interface{}]*entry, len(read.m))
    for k, e := range read.m {
        // 如果标记为nil或者expunged，则不复制到dirty map
        if !e.tryExpungeLocked() {
            m.dirty[k] = e
        }
    }
}
```

总结：在写操作时，会先去read中看是否存在该key，如果存在，只需要原子操作更新该值即可，如果不存在，则需要加锁到dirty中去操作

LoadOrStore方法
```go
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
    // 不加锁的情况下读取read map
    // 第一次检测
    read, _ := m.read.Load().(readOnly)
    if e, ok := read.m[key]; ok {
        // 如果元素存在（是否标识为删除由tryLoadOrStore执行处理），尝试获取该元素已存在的值或者将元素写入
        actual, loaded, ok := e.tryLoadOrStore(value)
        if ok {
            return actual, loaded
        }
    }

    m.mu.Lock()
    // 第二次检测，参看Store方法
    read, _ = m.read.Load().(readOnly)
    if e, ok := read.m[key]; ok {
        if e.unexpungeLocked() {
            m.dirty[key] = e
        }
        actual, loaded, _ = e.tryLoadOrStore(value)
    } else if e, ok := m.dirty[key]; ok {
        actual, loaded, _ = e.tryLoadOrStore(value)
        m.missLocked()
    } else {
        if !read.amended {
            m.dirtyLocked()
            m.read.Store(readOnly{m: read.m, amended: true})
        }
        m.dirty[key] = newEntry(value)
        actual, loaded = value, false
    }
    m.mu.Unlock()

    return actual, loaded
}
```

//如果没有删除元素，tryLoadOrStore将自动加载或存储一个值。如果删除元素，tryLoadOrStore保持条目不变并返回ok= false。
```go
func (e *entry) tryLoadOrStore(i interface{}) (actual interface{}, loaded, ok bool) {
    p := atomic.LoadPointer(&e.p)
    // 元素标识删除，直接返回
    if p == expunged {
        return nil, false, false
    }
    // 存在该元素真实值，则直接返回原来的元素值
    if p != nil {
        return *(*interface{})(p), true, true
    }

    // 如果p为nil(此处的nil，并是不是指元素的值为nil，而是atomic.LoadPointer(&e.p)为nil，元素的nil在unsafe.Pointer是有值的)，则更新该元素值
    ic := i
    for {
        if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
            return i, false, true
        }
        p = atomic.LoadPointer(&e.p)
        if p == expunged {
            return nil, false, false
        }
        if p != nil {
            return *(*interface{})(p), true, true
        }
    }
}
```

Delete方法
```
func (m *Map) Delete(key interface{}) {
    // 第一次检测
    read, _ := m.read.Load().(readOnly)
    e, ok := read.m[key]
    if !ok && read.amended {
        m.mu.Lock()
        // 第二次检测
        read, _ = m.read.Load().(readOnly)
        e, ok = read.m[key]
        if !ok && read.amended {
            // 不论dirty map是否存在该元素，都会执行删除
            delete(m.dirty, key)
        }
        m.mu.Unlock()
    }
    if ok {
        // 如果在read中，则将其标记为删除（nil）
        e.delete()
    }
}
func (e *entry) delete() (hadValue bool) {
    for {
        p := atomic.LoadPointer(&e.p)
        if p == nil || p == expunged {
            return false
        }
        if atomic.CompareAndSwapPointer(&e.p, p, nil) {
            return true
        }
    }
}
```

总结：删除操作会首先去read中查找，如果有，会直接标记为nil，否则，会加锁到dirty中进行真正的删除操作

Range方法
```go
func (m *Map) Range(f func(key, value interface{}) bool) {
    // first check
    read, _ := m.read.Load().(readOnly)
    // read.amended=true,说明dirty map包含所有有效的元素（含新加，不含被删除的），使用dirty map
    if read.amended {
        // double check
        m.mu.Lock()
        read, _ = m.read.Load().(readOnly)
        if read.amended {
            // 使用dirty map并且升级为read map
            read = readOnly{m: m.dirty}
            m.read.Store(read)
            m.dirty = nil
            m.misses = 0
        }
        m.mu.Unlock()
    }
    //使用read map作为读
    for k, e := range read.m {
        v, ok := e.load()
        // 被删除的不计入
        if !ok {
            continue
        }
        // 函数返回false，终止
        if !f(k, v) {
            break
        }
    }
}
```

总结：由代码可知，Range方法最终读的是read中的数据，所以不一定准确

应用场景
由以上方法可得知，无论是读操作，还是更新操作，亦或者删除操作，都会先从read进行操作，因为read的读取更新不需要锁，是原子操作，这样既做到了并发安全，又做到了尽量减少锁的争用，虽然采用的是空间换时间的策略，通过两个冗余的map，实现了这一点，但是底层存的都是指针类型，所以对于空间占用，也是做到了最大程度的优化。
但是同时也可以得知，当存在大量写操作时，会导致read中读不到数据，依然会频繁加锁，同时dirty升级为read，整体性能就会很低，所以sync.Map更加适合大量读、少量写的场景。

如何保证键类型和值类型的正确性
在讲到map时，我们提到了map的键类型的不能是哪些类型：即函数类型、切片类型、字典类型，同样的，对sync.map的键类型也是一样的，不能为这3种类型。
而sync.map中涉及到的键和值类型都为interface{}类型，也就是空接口，所以必须依赖我们自己来保证键类型和值类型的正确性，那么问题就来了：如何保证并发安全字典中键和值类型的正确性？

方案一：让sync.Map只存储某个特定类型的键
指定键类型只能是int、string、或者某类结构体，一旦完全确定键类型之后，就可以通过使用类型断言表达式对键的类型做检查，而且，如果你要是把并发安全字典封装到一个结构体类型里面，会更加方便，完全可以让GO语言编译器帮助做类型检查
```go
// IntStrMap 代表键类型为int、值类型为string的并发安全字典。
type IntStrMap struct {
	m sync.Map
}

func (iMap *IntStrMap) Delete(key int) {
	iMap.m.Delete(key)
}

func (iMap *IntStrMap) Load(key int) (value string, ok bool) {
	v, ok := iMap.m.Load(key)
	if v != nil {
		value = v.(string)
	}
	return
}

func (iMap *IntStrMap) LoadOrStore(key int, value string) (actual string, loaded bool) {
	a, loaded := iMap.m.LoadOrStore(key, value)
	actual = a.(string)
	return
}

func (iMap *IntStrMap) Range(f func(key int, value string) bool) {
	f1 := func(key, value interface{}) bool {
		return f(key.(int), value.(string))
	}
	iMap.m.Range(f1)
}

func (iMap *IntStrMap) Store(key int, value string) {
	iMap.m.Store(key, value)
}
```

以上代码中的方法除了确定了类型外，其他都是一样的，这些方法在接受键和值的时候，就不用再做类型检查了，取值的时候，也不用担心类型不正确，因为在当初存入的时候，就已经由编译器保证了。

方案一的实现很简单，但是缺点也是显而易见的，非常不灵活，不能灵活改变键和值的类型，需求多了之后，会产生很多雷同的代码，因此我们来看方案二

方案二：封装的结构体类型的所有方法，与sync.Map类型完全一致，此时需要类型检查
```go
// ConcurrentMap 代表可自定义键类型和值类型的并发安全字典。
type ConcurrentMap struct {
	m         sync.Map
	keyType   reflect.Type  //键类型
	valueType reflect.Type  //值类型
}

func NewConcurrentMap(keyType, valueType reflect.Type) (*ConcurrentMap, error) {
	if keyType == nil {
		return nil, errors.New("nil key type")
	}
	if !keyType.Comparable() {
		return nil, fmt.Errorf("incomparable key type: %s", keyType)
	}
	if valueType == nil {
		return nil, errors.New("nil value type")
	}
	cMap := &ConcurrentMap{
		keyType:   keyType,
		valueType: valueType,
	}
	return cMap, nil
}

func (cMap *ConcurrentMap) Delete(key interface{}) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	cMap.m.Delete(key)
}

func (cMap *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	return cMap.m.Load(key)
}

func (cMap *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	actual, loaded = cMap.m.LoadOrStore(key, value)
	return
}

func (cMap *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	cMap.m.Range(f)
}

func (cMap *ConcurrentMap) Store(key, value interface{}) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	cMap.m.Store(key, value)
}
```

在我们初始化结构体的时候，键类型和值类型就需要完全确定了，并且在这种情况下，必须先要保证键的类型是可比较的，以上代码中的panic完全可以根据自己的需求，替换为error。

可以看到第二种方案完全弥补了第一种方案的缺陷，但是使用了反射，这会降低程序的性能，因此选择哪种方案，还需要根据自身的需求来决定


原文链接：https://blog.csdn.net/xixisuli/article/details/89190856