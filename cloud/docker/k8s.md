# 1 K8s简介

Kubernetes(k8s)是Google开源的容器集群管理系统，它主要用于 容器编排、启动容器、自动化部署、扩展与管理容器应用、回收容器等。k8s的目标是让部署容器化的应用简单并且高效，k8s提供了应用部署、规划、更新、维护的一种机制！

简单地，可以理解为K8s是个管家，负责每个小屋子的监控、流通和控制。

## 1.1 为什么需要k8s

1.应用部署模式的演进

（1）传统的虚拟化部署模式，如下图所示：

![虚拟化部署模式](https://upload-images.jianshu.io/upload_images/6534887-f73947baaea7ccc7.png?imageMogr2/auto-orient/strip|imageView2/2/w/441/format/webp "虚拟化部署模式")

（2）容器化部署模式

![容器化部署模式](https://upload-images.jianshu.io/upload_images/6534887-53a8d68d6e0f3af0.png?imageMogr2/auto-orient/strip|imageView2/2/w/441/format/webp "容器化部署模式")

容器相比虚拟机，其优点如下：
```
容器更加轻量级，启动更快（秒级）
容器可移植性更好
```

2.管理大量的容器带来了新的挑战


## 1.2 k8s特点

容器编排调度引擎 —— k8s 的好处
```
简化应用部署
提高硬件资源利用率
健康检查和自修复
自动扩容缩容
服务发现和负载均衡
具有多租户和 Namespace隔离机制，确保每个工程执行单元测试时能做到隔离和并行。
```

# 2 原理

## 2.1 k8s集群架构

一个k8s系统，通常称为一个k8s集群（cluster）。k8s集群，同样也类似于主从结构：一个主节点和多个工作节点。这个集群主要包括两个部分：
```
一个Master，即主节点，控制和管理整个集群系统，提供集群的资源数据访问入口。

一群Node，即工作节点，里面是具体的容器，用来承载被分配Pod的运行，是Pod运行的宿主机，负责运行用户实际的应用。
```


### 2.1.1 Master节点构成

Master节点包括API Server、Scheduler、Controller manager、etcd、replication controller。
```
API Server：是整个系统的对外接口，供客户端和其它组件调用，相当于“营业厅”。

Scheduler：负责对集群内部的资源进行调度，相当于“调度室”。

Controller manager：负责管理控制器，相当于“大总管”。

ETCD：强一致性的键值对存储，k8s 集群中的所有资源对象都存储在 etcd 中。

Replication Controller是实现弹性伸缩、动态扩容和滚动升级的核心。
```
![Master节点](https://oscimg.oschina.net/oscnet/c876af8baae0cade281c52b4be3c253cbe7.jpg "Master节点构成")

### 2.1.2 Node节点构成

Node节点包括Docker、kubelet、kube-proxy、Fluentd、kube-dns（可选），还有就是Pod。
```
Pod是Kurbernetes进行创建、调度和管理的最小单位，它提供了比容器更高层次的抽象，使得部署和管理更加灵活。

kubelet：负责对Pod对应的容器的创建、启停等任务

kube-proxy：实现Kubernetes Service的通信与负载均衡机制的重要组件

容器运行时——docker：负责管理 node 节点上的所有容器和容器 IP 的分配。
```

Pod是Kubernetes最基本的操作单元。一个Pod代表着集群中运行的一个进程，它内部封装了一个或多个紧密相关的容器。除了Pod之外，K8S还有一个Service的概念，一个Service可以看作一组提供相同服务的Pod的对外访问接口。

![Node节点](https://oscimg.oschina.net/oscnet/e60f0bea1561054c44f267ed82596d1219b.jpg "Node节点构成")

## 2.2 调度策略

调度策略主要有：Scheduler调度策略、亲和调度、定向调度等。
K8s默认使用cheduler调度策略：
```
首先，用户提交自己的job到集群中；

然后，scheduler会经常检查有没有job（用来生成pod）；

接着，如果有，就会根据pod里面的要求寻找一个节点（node）。
```
其中，寻找节点的过程是这样的：
```
首先，要判断这个node是不是符合基本条件，比如说内存大小，cpu等；
其次，对符合基本条件的node，按照其他条件打分（比如负载均衡等）；
然后，选择打分最高的node。
```

# 3 问题与挑战

## 3.1 集群扩广遇到的挑战

这道题主要考察在扩广k8s 集群实现微服务容器化部署实际落地过程中遇到的挑战和踩过的坑有哪些，话题有点广，可以说的点其实挺多的，可以主要从以下几个方面来阐述的。

### 3.1.1 部署的规范流程

虽然说容器和虚拟机部署本质上没有多大区别，但还是有些许不同的。容器的可执行文件是一个镜像，而虚拟机的可执行文件往往是一个二进制文件如 jar 包或者是 war包，另外，由于容器隔离的不是特别彻底，在上文也有所阐述，针对这种情况，如何更准确获取 cgroups 给容器限定的 Memory 和 CPU 值，这给平台开发者带来相应的挑战。此外，在容器化部署时，作为用户而言，需要遵循相应的使用规范和流程，如每个 Pod 都必须设置资源限额和健康检测探针，在设置资源限额时，又不能盲目设置，需要依赖监控组件或者是开发者本身对自身应用的认知，进行相关经验值的设置。

### 3.1.2 多集群调度

对于如何管理多个 k8s 集群，如何进行跨集群调度、应用部署和资源对象管理，这对于平台本身，都是一个很大的挑战。

调度均衡问题

随着集群规模的扩大以及微服务部署的数量增加，同一个计算节点，可能会运行很多 Pod，这个时候就会出现资源争用的问题。k8s 本身调度层面有两个阶段，分别是预选阶段和优选阶段，每个阶段都有对应的调度策略和算法，关于如何均衡节点之后的调度，这需要在平台层面上对调度算法有所研究，并进行适当的调整。


## 3.2 如何解决Memory 和 CPU 隔离不彻底问题

由于 /proc 文件系统是以只读的方式挂载到容器内部，所以在容器内看到的都是宿主机的信息，包括 CPU 和 Memory，docker 是以 cgroups 来进行资源限制的，而 jdk1.9 以下版本目前无法自动识别容器的资源配额，1.9以上版本会自动识别和正常读取 cgroups 中为容器限制的资源大小。

### 3.2.1 Memory 隔离不彻底问题

Docker 通过 cgroups 完成对内存的限制，而 /proc 文件目录是以只读的形式挂载到容器中，由于默认情况下，Java 压根就看不到 cgroups 限制的内容的大小，而默认使用 /proc/meminfo 中的信息作为内存信息进行启动，默认情况下，JVM 初始堆大小为内存总量的 1/4，这种情况会导致，如果容器分配的内存小于 JVM 的内存， JVM 进程会被 linux killer 杀死。

那么目前有几种解决方式：
```
（1）升级 JDK 版本到1.9及以上，让 JVM 能自动识别 cgroups 对容器的资源限制，从而自动调整 JVM 的参数并启动 JVM 进程。

（2）对于较低版本的JDK，一定要设置 JVM 初始堆大小，并且JVM 的最大堆内存不能超过容器的最大内存值，正常理论值应该是：容器 limit-memory = JVM 最大堆内存 + 750MB。

（3）使用 lxcfs ，这是一种用户态文件系统，用来支持LXC 容器，lxcfs 通过用户态文件系统，在容器中提供下列 procfs 的文件，启动时，把宿主机对应的目录 /var/lib/lxcfu/proc/meminfo 文件挂载到 Docker 容器的 /proc/meminfo 位置后，容器中进程（JVM）读取相应文件内容时，lxcfs 的 fuse 将会从容器对应的 cgroups 中读取正确的内存限制，从而获得正确的资源约束设定。
```

### 3.2.2 CPU 隔离不彻底问题

JVM GC （垃圾回收）对于 java 程序执行性能有一定的影响，默认的 JVM 使用如下公式： ParallelGCThreads = ( ncpu <= 8 ) ? ncpu：3 + （ncpu * 5）/ 8 来计算并行 GC 的线程数，但是在容器里面，ncpu 获取的就是所在宿主机的 cpu 个数，这会导致 JVM 启动过多的 GC 线程，直接的结果就是 GC 的性能下降，java 服务的感受就是：延时增加， TPS 吞度量下降，针对这种问题，也有以下几种解决方案：
```
（1）显示传递 JVM 启动参数：“-XX: ParallelGCThreads" 告诉 JVM 应该启动多少个并行 GC 线程，缺点是需要业务感知，而且需要为不同配置的容器传递不同的 JVM 参数。

（2）在容器内使用 Hack 过的 glibc ，使 JVM 通过 sysconf 系统调用能正确获取容器内 CPU 资源核数，优点是业务无感知，并且能自动适配不同配置的容器，缺点是有一定的维护成本。
```

# 4 参考

k8s简介：https://www.jianshu.com/p/502544957c88

docker & kubernetes 面试（某互联网公司）：https://www.jianshu.com/p/2de643caefc1

docker与k8s核心概念学习笔记：https://zhuanlan.zhihu.com/p/54861341