# Prometheus

1.Dev（开发）ops（运维）：DevOps经常被描述为“开发团队与运营团队之间更具协作性、更高效的关系”。由于团队间协作关系的改善，整个组织的效率因此得到提升，伴随频繁变化而来的生产环境的风险也能得到降低。

推荐文章：https://zh.wikipedia.org/wiki/DevOps

## 1.1 Prometheus是什么

Prometheus 是一套开源的系统监控报警框架。

## 1.2 Prometheus特点

（1）强大的多维度数据模型：

时间序列数据通过 metric 名和键值对来区分。
所有的 metrics 都可以设置任意的多维标签。
数据模型更随意，不需要刻意设置为以点分隔的字符串。
可以对数据模型进行聚合，切割和切片操作。
支持双精度浮点类型，标签可以设为全 unicode。

（2）灵活而强大的查询语句（PromQL）

在同一个查询语句，可以对多个 metrics 进行乘法、加法、连接、取分数位等操作。

（3）易于管理

Prometheus server 是一个单独的二进制文件，可直接在本地工作，不依赖于分布式存储。

（4)高效

平均每个采样点仅占 3.5 bytes，且一个 Prometheus server 可以处理数百万的 metrics。

（5）使用 pull 模式采集时间序列数据，这样不仅有利于本机测试而且可以避免有问题的服务器推送坏的 metrics。

（6）可以采用 push gateway 的方式把时间序列数据推送至 Prometheus server 端。

（7）可以通过服务发现或者静态配置去获取监控的 targets。

（8）有多种可视化图形界面。

（9）易于伸缩。

需要指出的是，由于数据采集可能会有丢失，所以 Prometheus 不适用对采集数据要 100% 准确的情形。但如果用于记录时间序列数据，Prometheus 具有很大的查询优势，此外，Prometheus 适用于微服务的体系架构。

参考文章：

https://www.ibm.com/developerworks/cn/cloud/library/cl-lo-prometheus-getting-started-and-practice/index.html

