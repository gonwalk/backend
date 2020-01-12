本文最常见的Kafka面试题，同时也是对Apache Kafka初学者必备知识点的一个整理与介绍。

参考：
14个最常见的Kafka面试题及答案
http://www.toutiao.com/i6456660580726997517/



# 1 Kafka基本概念

Apache Kafka是由Apache开发的一种发布订阅消息系统，它是一个分布式的、分区的和重复的日志服务。

## 1.1 Kafka概念与优势

### 1.1.1 什么是kafka

       Kafka是分布式发布-订阅消息系统，它最初是由LinkedIn公司开发的，之后成为Apache项目的一部分，Kafka是一个分布式，可分区的、冗余备份的、持久性的日志服务，它主要用于处理流式数据。

 

### 1.1.2 为什么要使用 kafka? 

     缓冲和削峰：上游数据时有突发流量，下游可能扛不住，或者下游没有足够多的机器来保证冗余，kafka在中间可以起到一个缓冲的作用，把消息暂存在kafka中，下游服务就可以按照自己的节奏进行慢慢处理。

     解耦和扩展性：项目开始的时候，并不能确定具体需求。消息队列可以作为一个接口层，解耦重要的业务流程。只需要遵守约定，针对数据编程即可获取扩展能力。

     冗余：可以采用一对多的方式，一个生产者发布消息，可以被多个订阅topic的服务消费到，供多个毫无关联的业务使用。

     健壮性：消息队列可以堆积请求，所以消费端业务即使短时间死掉，也不会影响主要业务的正常进行。


### 1.1.3 Kafka相对传统技术有什么优势?

Apache Kafka与传统的消息传递技术相比优势在于：
```
快速:单一的Kafka代理可以处理成千上万的客户端，每秒处理数兆字节的读写操作。

可伸缩:在一组机器上对数据进行分区和简化，以支持更大的数据

持久:消息是持久性的，并在集群中进行复制，以防止数据丢失。

容错保证:它提供了容错保证和持久性。
```

### 1.1.4 Kafka的特性
```
高吞吐量、低延迟：kafka每秒可以处理几十万条消息，它的延迟最低只有几毫秒，每个topic可以分多个partition, consumer group 对partition进行consume操作。

可扩展性：kafka集群支持热扩展。

持久性、可靠性：消息被持久化到本地磁盘，并且支持数据备份防止数据丢失。

容错性：允许集群中节点失败（若副本数量为n,则最多允许n-1个节点失败）。

高并发：支持数千个客户端同时读写。
```

链接：https://www.jianshu.com/p/18e1a9ae6f9c

### 1.1.5 kafka使用场景

1.分布式消息中间件；
2.跟踪网站活动；
3.日志聚合。
日志聚合系统通常从服务器收集物理日志文件，并将其抽象为一个中心系统（可能是文件服务器或HDFS）进行处理。
kafka从这些日志文件中提取信息，并将其抽象为一个更加清晰的消息流。这样可以实现更低的延迟处理，且易于支持多个数据源及分布式数据的消耗。


## 1.2 消息队列
### 1.2.1 什么是传统的消息传递方法?

传统的消息传递方法包括两种：
```
排队：在队列中，一组用户可以从服务器中读取消息，每条消息都发送给其中一个人。

发布-订阅：在这个模型中，消息被广播给所有的用户。
```

### 1.2.2 为什么要使用消息队列?
      异步通信：很多时候，用户不想也不需要立即处理消息。消息队列提供了异步处理机制，允许用户把一个消息放入队列，但并不立即处理它。想向队列中放入多少消息就放多少，然后在需要的时候再去处理它们。

## 1.3 术语
Kafka 中的术语
![kafka应用](https://upload-images.jianshu.io/upload_images/189732-96c861447d01a82d.png?imageMogr2/auto-orient/strip|imageView2/2/w/952/format/webp)
```
Broker：中间的kafka cluster，存储消息，是由多个server组成的集群。任何正在运行中的Kafka示例都称为Broker。
Topic：Topic其实就是一个传统意义上的消息队列，可以看做kafka给消息提供的分类方式。broker用来存储不同topic的消息数据。
Partition：即分区。一个Topic将由多个分区组成，每个分区将存在独立的持久化文件，任何一个Consumer在分区上的消费一定是顺序的；当一个Consumer同时在多个分区上消费时，Kafka不能保证总体上的强顺序性（对于强顺序性的一个实现是Exclusive Consumer，即独占消费，一个队列同时只能被一个Consumer消费，并且从该消费开始消费某个消息到其确认才算消费完成，在此期间任何Consumer不能再消费）。（通常有几个节点就初始化几个Partition，Partition指的是一个topic有几个，分布在四个节点，把四个节点的数据均分在不同的partition实现均衡）
producer：往broker中某个topic里面生产数据。
consumer：从broker中某个topic获取数据。
Consumer Group：即消费组。一个消费组是由一个或者多个Consumer组成的，对于同一个Topic，不同的消费组都将能消费到全量的消息，而同一个消费组中的Consumer将竞争每个消息（在多个Consumer消费同一个Topic时，Topic的任何一个分区将同时只能被一个Consumer消费）。
```

kafka工作原理介绍 https://www.jianshu.com/p/0272d5e4ffad



## 1.4 Kafka与ActiveMQ、RabbitMQ之间的比较

Kafka的topic只能精确匹配。Kafka的producer和consumer集群需要依赖zookeeper保证高可用性。从数据吞吐量来说，ActiveMQ最低，Kafka最高。从存储可靠性来说，RabbitMQ最高，Kafka最低。从支持的协议来说，ActiveMQ最多。

Kafka是partition组成的，一个partition可以有多个消费者，但是多个消费者必须在不同的组里。如果要做到高吞吐量，可以通过多个线程在不同的offset开始读取message，然后处理对应的message即可。

![Kafka与ActiveMQ、RabbitMQ之间的比较](https://upload-images.jianshu.io/upload_images/16132650-9feb3a4f3fa2d7fc.png?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp)

链接：https://www.jianshu.com/p/18e1a9ae6f9c



# 2 kafka原理

## 2.1 消息队列Message Queue 的通讯模式
```
点对点通讯：点对点方式是最为传统和常见的通讯方式，它支持一对一、一对多、多对多、多对一等多种配置方式，支持树状、网状等多种拓扑结构。

多点广播：MQ 适用于不同类型的应用。其中重要的，也是正在发展中的是"多点广播"应用，即能够将消息发送到多个目标站点 (Destination List)。可以使用一条 MQ 指令将单一消息发送到多个目标站点，并确保为每一站点可靠地提供信息。MQ 不仅提供了多点广播的功能，而且还拥有智能消息分发功能，在将一条消息发送到同一系统上的多个用户时，MQ 将消息的一个复制版本和该系统上接收者的名单发送到目标 MQ 系统。目标 MQ 系统在本地复制这些消息，并将它们发送到名单上的队列，从而尽可能减少网络的传输量。

发布/订阅 (Publish/Subscribe) 模式：发布/订阅功能使消息的分发可以突破目的队列地理指向的限制，使消息按照特定的主题甚至内容进行分发，用户或应用程序可以根据主题或内容接收到所需要的消息。发布/订阅功能使得发送者和接收者之间的耦合关系变得更为松散，发送者不必关心接收者的目的地址，而接收者也不必关心消息的发送地址，而只是根据消息的主题进行消息的收发。

群集 (Cluster)：为了简化点对点通讯模式中的系统配置，MQ 提供 Cluster(群集) 的解决方案。群集类似于一个域 (Domain)，群集内部的队列管理器之间通讯时，不需要两两之间建立消息通道，而是采用群集 (Cluster) 通道与其它成员通讯，从而大大简化了系统配置。此外，群集中的队列管理器之间能够自动进行负载均衡，当某一队列管理器出现故障时，其它队列管理器可以接管它的工作，从而大大提高系统的高可靠性。
```

## 2.2 Kafka的Leader的选举机制
### 2.2.1 Kafka的Leader是什么
首先Kafka会将接收到的消息分区（partition），每个主题（topic）的消息有不同的分区。这样一方面消息的存储就不会受到单一服务器存储空间大小的限制，另一方面消息的处理也可以在多个服务器上并行。
其次为了保证高可用，每个分区都会有一定数量的副本（replica）。这样如果有部分服务器不可用，副本所在的服务器就会接替上来，保证应用的持续性。


但是，为了保证较高的处理效率，消息的读写都是在固定的一个副本上完成。这个副本就是所谓的Leader，而其他副本则是Follower。而Follower则会定期地到Leader上同步数据。
### 2.2.2 Leader选举机制
如果某个分区所在的服务器出了问题，不可用，kafka会从该分区的其他的副本中选择一个作为新的Leader。之后所有的读写就会转移到这个新的Leader上。现在的问题是应当选择哪个作为新的Leader。显然，只有那些跟Leader保持同步的Follower才应该被选作新的Leader。
Kafka会在Zookeeper上针对每个Topic维护一个称为ISR（in-sync replica，已同步的副本）的集合，该集合中是一些分区的副本。只有当这些副本都跟Leader中的副本同步了之后，kafka才会认为消息已提交，并反馈给消息的生产者。如果这个集合有增减，kafka会更新zookeeper上的记录。

如果某个分区的Leader不可用，Kafka就会从ISR集合中选择一个副本作为新Leader。
显然通过ISR，kafka需要的冗余度较低，可以容忍的失败数比较高。假设某个topic有f+1个副本，kafka可以容忍f个服务器不可用。

### 2.2.3 为什么不用少数服从多数的方法
少数服从多数是一种比较常见的一致性算法和Leader选举法。它的含义是只有超过半数的副本同步了，系统才会认为数据已同步；选择Leader时也是从超过半数的同步的副本中选择。这种算法需要较高的冗余度。譬如只允许一台机器失败，需要有三个副本；而如果只容忍两台机器失败，则需要五个副本。而kafka的ISR集合方法，分别只需要两个和三个副本。
### 2.2.4 如果所有的ISR副本都失败了怎么办
此时有两种方法可选，一种是等待ISR集合中的副本复活，一种是选择任何一个立即可用的副本，而这个副本不一定是在ISR集合中。这两种方法各有利弊，实际生产中按需选择。
如果要等待ISR副本复活，虽然可以保证一致性，但可能需要很长时间。而如果选择立即可用的副本，则很可能该副本并不一致。
## 2.3 kafka集群partition分布原理分析

![kafka集群中的partition分布](https://img-blog.csdn.net/20180720193030886?watermark/2/text/aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L1Bhbl9ZVA==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70)
在Kafka集群中，每个Broker都有均等分配Partition的Leader机会。
上述图Broker Partition中，箭头指向为副本，以Partition-0为例:broker1中parition-0为Leader，Broker2中Partition-0为副本。
上述图种每个Broker(按照BrokerId有序)依次分配主Partition,下一个Broker为副本，如此循环迭代分配，多副本都遵循此规则。

副本分配算法如下：
将所有n个Broker和待分配的i个Partition排序.
将第i个Partition分配到第(i mod n)个Broker上.
将第i个Partition的第j个副本分配到第((i + j) mod n)个Broker上.

## 2.4 Kafka中的Zookeeper是什么?
Zookeeper是一个开放源码的、高性能的协调服务，它用于Kafka的分布式应用。

Zookeeper主要用于在集群中不同节点之间进行通信。

在Kafka中，Zookeeper被用于提交偏移量，如果节点在所有情况下都失败了，它也可以从之前提交的偏移量中获取。

除此之外，它还执行其他活动，如: leader检测、分布式同步、配置管理、识别新节点何时离开或连接集群、节点实时状态等等。

### 2.4.1 可以在没有Zookeeper的情况下使用Kafka吗?

不能。不可能越过Zookeeper，直接联系Kafka broker。一旦Zookeeper停止工作，它就不能服务客户端请求。

### 2.4.2 Zookeeper在kafka的作用
无论是kafka集群，还是producer和consumer都依赖于zookeeper来保证系统可用性，保证集群保存一些meta信息。
Kafka使用zookeeper作为其分布式协调框架，很好的将消息生产、消息存储、消息消费的过程结合在一起。
同时借助zookeeper，kafka能够将生产者、消费者和broker在内的所以组件在无状态的情况下，建立起生产者和消费者的订阅关系，并实现生产者与消费者的负载均衡。

## 2.5 kafka消费信息

### 2.5.1 Kafka的用户如何消费信息?

在Kafka中传递消息是通过使用sendfile API完成的。它支持将字节从套接口转移到磁盘，通过内核空间保存副本，并在内核用户之间调用内核。


4、在Kafka中broker的意义是什么?

在Kafka集群中，broker术语用于引用服务器。

5、Kafka服务器能接收到的最大信息是多少?

Kafka服务器默认可以接收到的最大消息的大小是1000000字节=1MB。



8、解释如何提高远程用户的吞吐量?

如果用户位于与broker不同的数据中心，则可能需要调优套接口缓冲区大小，以对长网络延迟进行摊销。

9、解释一下，在数据制作过程中，你如何能从Kafka得到准确的信息?

在数据中，为了精确地获得Kafka的消息，你必须遵循两件事: 在数据消耗期间避免重复，在数据生产过程中避免重复。

这里有两种方法，可以在数据生成时准确地获得一个语义:

每个分区使用一个单独的写入器，每当你发现一个网络错误，检查该分区中的最后一条消息，以查看您的最后一次写入是否成功

在消息中包含一个主键(UUID或其他)，并在用户中进行反复制

10、解释如何减少ISR中的扰动?broker什么时候离开ISR?

ISR是一组与leaders完全同步的消息副本，也就是说ISR中包含了所有提交的消息。ISR应该总是包含所有的副本，直到出现真正的故障。如果一个副本从leader中脱离出来，将会从ISR中删除。

11、Kafka为什么需要复制?

Kafka的信息复制确保了任何已发布的消息不会丢失，并且可以在机器错误、程序错误或更常见些的软件升级中使用。

12、如果副本在ISR中停留了很长时间表明什么?

如果一个副本在ISR中保留了很长一段时间，那么它就表明，跟踪器无法像在leader收集数据那样快速地获取数据。

13、请说明如果首选的副本不在ISR中会发生什么?

如果首选的副本不在ISR中，控制器将无法将leadership转移到首选的副本。

14、有可能在生产后发生消息偏移吗?

在大多数队列系统中，作为生产者的类无法做到这一点，它的作用是触发并忘记消息。broker将完成剩下的工作，比如使用id进行适当的元数据处理、偏移量等。

作为消息的用户，你可以从Kafka broker中获得补偿。如果你注视SimpleConsumer类，你会注意到它会获取包括偏移量作为列表的MultiFetchResponse对象。此外，当你对Kafka消息进行迭代时，你会拥有包括偏移量和消息发送的MessageAndOffset对象。

转自

14个最常见的Kafka面试题及答案
http://www.toutiao.com/i6456660580726997517/


参考：
kafka面试题(附答案)：https://blog.csdn.net/qq_23160237/article/details/88376561


 

3.kafka中的broker 是干什么的?

      broker 是消息的代理，Producers往Brokers里面的指定Topic中写消息，Consumers从Brokers里面拉取指定Topic的消息，然后进行业务处理，broker在中间起到一个代理保存消息的中转站。

 

4.kafka中的 zookeeper 起到什么作用，可以不用zookeeper么?

       zookeeper 是一个分布式的协调组件，早期版本的kafka用zk做meta信息存储，consumer的消费状态，group的管理以及 offset的值。考虑到zk本身的一些因素以及整个架构较大概率存在单点问题，新版本中逐渐弱化了zookeeper的作用。新的consumer使用了kafka内部的group coordination协议，也减少了对zookeeper的依赖，但是broker依然依赖于ZK，zookeeper 在kafka中还用来选举controller 和 检测broker是否存活等等。

 

5.kafka follower如何与leader同步数据?

        Kafka的复制机制既不是完全的同步复制，也不是单纯的异步复制。完全同步复制要求All Alive Follower都复制完，这条消息才会被认为commit，这种复制方式极大的影响了吞吐率。而异步复制方式下，Follower异步的从Leader复制数据，数据只要被Leader写入log就被认为已经commit，这种情况下，如果leader挂掉，会丢失数据，kafka使用ISR的方式很好的均衡了确保数据不丢失以及吞吐率。Follower可以批量的从Leader复制数据，而且Leader充分利用磁盘顺序读以及send file(zero copy)机制，这样极大的提高复制性能，内部批量写磁盘，大幅减少了Follower与Leader的消息量差。

 

6.什么情况下一个 broker 会从 isr中踢出去

         leader会维护一个与其基本保持同步的Replica列表，该列表称为ISR(in-sync Replica)，每个Partition都会有一个ISR，而且是由leader动态维护 ，如果一个follower比一个leader落后太多，或者超过一定时间未发起数据复制请求，则leader将其重ISR中移除 ，详细参考 kafka的高可用机制

 

7.kafka 为什么那么快?

Cache Filesystem Cache PageCache缓存

顺序写 由于现代的操作系统提供了预读和写技术，磁盘的顺序写大多数情况下比随机写内存还要快。

Zero-copy 零拷贝技术减少拷贝次数

Batching of Messages 批量量处理。合并小的请求，然后以流的方式进行交互，直顶网络上限。

Pull 拉模式 使用拉模式进行消息的获取消费，与消费端处理能力相符。

 

8.kafka producer如何优化写入速度？

增加线程

提高 batch.size

增加更多 producer 实例

增加 partition 数

设置 acks=-1 时，如果延迟增大：可以增大 num.replica.fetchers（follower 同步数据的线程数）来调解；

跨数据中心的传输：增加 socket 缓冲区设置以及 OS tcp 缓冲区设置。

优化方面的参考 kafka最佳实践

 

9.kafka producer 写数据，ack  为 0， 1， -1 的时候代表啥， 设置 -1 的时候，什么情况下，leader 会认为一条消息 commit了?

1（默认）  数据发送到Kafka后，经过leader成功接收消息的的确认，就算是发送成功了。在这种情况下，如果leader宕机了，则会丢失数据。
0 生产者将数据发送出去就不管了，不去等待任何返回。这种情况下数据传输效率最高，但是数据可靠性确是最低的。
-1 producer需要等待ISR中的所有follower都确认接收到数据后才算一次发送完成，可靠性最高。当ISR中所有Replica都向Leader发送ACK时，leader才commit，这时候producer才能认为一个请求中的消息都commit了。
 

10.kafka  unclean 配置代表啥，会对 spark streaming 消费有什么影响?

     unclean.leader.election.enable 为true的话，意味着非ISR集合的broker 也可以参与选举，这样有可能就会丢数据，spark streaming在消费过程中拿到的 end offset 会突然变小，导致 spark streaming job挂掉。如果unclean.leader.election.enable参数设置为true，就有可能发生数据丢失和数据不一致的情况，Kafka的可靠性就会降低；而如果unclean.leader.election.enable参数设置为false，Kafka的可用性就会降低。

 

11.如果leader crash时，ISR为空怎么办?

      kafka在Broker端提供了一个配置参数：unclean.leader.election,这个参数有两个值：
true（默认）：允许不同步副本成为leader，由于不同步副本的消息较为滞后，此时成为leader，可能会出现消息不一致的情况。
false：不允许不同步副本成为leader，此时如果发生ISR列表为空，会一直等待旧leader恢复，降低了可用性

 

12.kafka的message格式是什么样的?

      一个Kafka的Message由一个固定长度的header和一个变长的消息体body组成。header部分由一个字节的magic(文件格式)和四个字节的CRC32(用于判断body消息体是否正常)构成。当magic的值为1的时候，会在magic和crc32之间多一个字节的数据：attributes(保存一些相关属性， 比如是否压缩、压缩格式等等);如果magic的值为0，那么不存在attributes属性。body是由N个字节构成的一个消息体，包含了具体的key/value消息。

 

13.kafka中consumer group 是什么概念?

       同样是逻辑上的概念，是Kafka实现单播和广播两种消息模型的手段。同一个topic的数据，会广播给不同的group；同一个group中的worker，只有一个worker能拿到这个数据。换句话说，对于同一个topic，每个group都可以拿到同样的所有数据，但是数据进入group后只能被其中的一个worker消费。group内的worker可以使用多线程或多进程来实现，也可以将进程分散在多台机器上，worker的数量通常不超过partition的数量，且二者最好保持整数倍关系，因为Kafka在设计时假定了一个partition只能被一个worker消费（同一group内）。


原文链接：https://blog.csdn.net/qq_23160237/article/details/88376561


Kafka面试题与答案全套整理【转】https://blog.csdn.net/weixin_42139816/article/details/93327793

Kafka的用途有哪些？使用场景如何？
总结下来就几个字:异步处理、日常系统解耦、削峰、提速、广播
如果再说具体一点例如:消息,网站活动追踪,监测指标,日志聚合,流处理,事件采集,提交日志等

Kafka中的ISR、AR又代表什么？ISR的伸缩又指什么
ISR:In-Sync Replicas 副本同步队列
AR:Assigned Replicas 所有副本

ISR是由leader维护，follower从leader同步数据有一些延迟（包括延迟时间replica.lag.time.max.ms和延迟条数replica.lag.max.messages两个维度, 当前最新的版本0.10.x中只支持replica.lag.time.max.ms这个维度），任意一个超过阈值都会把follower剔除出ISR, 存入OSR（Outof-Sync Replicas）列表，新加入的follower也会先存放在OSR中。AR=ISR+OSR。

Kafka中的HW、LEO、LSO、LW等分别代表什么？
HW:High Watermark 高水位，取一个partition对应的ISR中最小的LEO作为HW，consumer最多只能消费到HW所在的位置上一条信息。
LEO:LogEndOffset 当前日志文件中下一条待写信息的offset

HW/LEO这两个都是指最后一条的下一条的位置而不是指最后一条的位置。

LSO:Last Stable Offset 对未完成的事务而言，LSO 的值等于事务中第一条消息的位置(firstUnstableOffset)，对已完成的事务而言，它的值同 HW 相同

LW:Low Watermark 低水位, 代表 AR 集合中最小的 logStartOffset 值

Kafka中是怎么体现消息顺序性的？
kafka每个partition中的消息在写入时都是有序的，消费时，每个partition只能被每一个group中的一个消费者消费，保证了消费时也是有序的。
整个topic不保证有序。如果为了保证topic整个有序，那么将partition调整为1.

Kafka中的分区器、序列化器、拦截器是否了解？它们之间的处理顺序是什么？
拦截器->序列化器->分区器

Kafka生产者客户端的整体结构是什么样子的？

Kafka生产者客户端中使用了几个线程来处理？分别是什么？
2个，主线程和Sender线程。主线程负责创建消息，然后通过分区器、序列化器、拦截器作用之后缓存到累加器RecordAccumulator中。Sender线程负责将RecordAccumulator中消息发送到kafka中.

Kafka的旧版Scala的消费者客户端的设计有什么缺陷？

“消费组中的消费者个数如果超过topic的分区，那么就会有消费者消费不到数据”这句话是否正确？如果不正确，那么有没有什么hack的手段？
不正确，通过自定义分区分配策略，可以将一个consumer指定消费所有partition。

消费者提交消费位移时提交的是当前消费到的最新消息的offset还是offset+1?
offset+1

有哪些情形会造成重复消费？
消费者消费后没有commit offset(程序崩溃/强行kill/消费耗时/自动提交偏移情况下unscrible)

那些情景下会造成消息漏消费？
消费者没有处理完消息 提交offset(自动提交偏移 未处理情况下程序异常结束)

KafkaConsumer是非线程安全的，那么怎么样实现多线程消费？
1.在每个线程中新建一个KafkaConsumer

2.单线程创建KafkaConsumer，多个处理线程处理消息（难点在于是否要考虑消息顺序性，offset的提交方式）

简述消费者与消费组之间的关系
消费者从属与消费组，消费偏移以消费组为单位。每个消费组可以独立消费主题的所有数据，同一消费组内消费者共同消费主题数据，每个分区只能被同一消费组内一个消费者消费。

当你使用kafka-topics.sh创建（删除）了一个topic之后，Kafka背后会执行什么逻辑？
创建:在zk上/brokers/topics/下节点 kafkabroker会监听节点变化创建主题
删除:调用脚本删除topic会在zk上将topic设置待删除标志，kafka后台有定时的线程会扫描所有需要删除的topic进行删除

topic的分区数可不可以增加？如果可以怎么增加？如果不可以，那又是为什么？
可以

topic的分区数可不可以减少？如果可以怎么减少？如果不可以，那又是为什么？
不可以

创建topic时如何选择合适的分区数？
根据集群的机器数量和需要的吞吐量来决定适合的分区数

Kafka目前有那些内部topic，它们都有什么特征？各自的作用又是什么？
__consumer_offsets 以下划线开头，保存消费组的偏移

优先副本是什么？它有什么特殊的作用？
优先副本 会是默认的leader副本 发生leader变化时重选举会优先选择优先副本作为leader

Kafka有哪几处地方有分区分配的概念？简述大致的过程及原理
创建主题时
如果不手动指定分配方式 有两种分配方式

消费组内分配

简述Kafka的日志目录结构
每个partition一个文件夹，包含四类文件.index .log .timeindex leader-epoch-checkpoint
.index .log .timeindex 三个文件成对出现 前缀为上一个segment的最后一个消息的偏移 log文件中保存了所有的消息 index文件中保存了稀疏的相对偏移的索引 timeindex保存的则是时间索引
leader-epoch-checkpoint中保存了每一任leader开始写入消息时的offset 会定时更新
follower被选为leader时会根据这个确定哪些消息可用

Kafka中有那些索引文件？
如上

如果我指定了一个offset，Kafka怎么查找到对应的消息？
1.通过文件名前缀数字x找到该绝对offset 对应消息所在文件

2.offset-x为在文件中的相对偏移

3.通过index文件中记录的索引找到最近的消息的位置

4.从最近位置开始逐条寻找

如果我指定了一个timestamp，Kafka怎么查找到对应的消息？
原理同上 但是时间的因为消息体中不带有时间戳 所以不精确

聊一聊你对Kafka的Log Retention的理解
kafka留存策略包括 删除和压缩两种
删除: 根据时间和大小两个方式进行删除 大小是整个partition日志文件的大小
超过的会从老到新依次删除 时间指日志文件中的最大时间戳而非文件的最后修改时间
压缩: 相同key的value只保存一个 压缩过的是clean 未压缩的dirty 压缩之后的偏移量不连续 未压缩时连续

