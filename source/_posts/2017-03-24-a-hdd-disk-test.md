title: 硬盘性能简测
date: 2017-03-24 18:48:55
tags:
---

各种存储系统，数据库、文件系统，在性能上无不都在与磁盘做斗争。希望能够尽量发挥系统有限的资源，提供最大化的性能。其中涉及到的技术包括

* page cache
* write buffer
* raid卡缓存
* WAL log
* bloom filter
* b+ tree storage
* skip list storage
* block cache
* memtable sstable
* 数据压缩,比如Snappy
* 等等

存储硬件在性能上除了近些年的SSD之外，HDD性能一直没有太大的飞跃式变化。当前磁盘(HDD)的性能究竟在什么水平，团队成员 @冉攀峰 同学 针对我们近期采购的磁盘做了详细的测试工作，主要是写入方面的性能测试（对于应用来说更多可控的空间在写入层的层面、读取性能的优化跟应用的workload相关），测试不同情况下磁盘在吞吐、延迟、iops上的性能水平。



#### 1 性能结果

测试工具：fio，详见:https://github.com/axboe/fio

测试硬盘：主流8T硬盘，详细这里不说明

#### 2 性能结果

以下主要为写入情况下的一些性能指标：

![](http://tompublic.nos-eastchina1.126.net/seqw_lat_iosize.png)

![](http://tompublic.nos-eastchina1.126.net/seqw_tps_iosize.png)

从上图可以得出如下几点：

1. 顺序写入情况下(1 thread),在2M左右基本是磁盘寻道延迟占主导因素(30ms左右)
2. 就写入优化而言，注意控制并发，因为写入相对可控，需要充分利用磁盘顺序IO的特性。
3. 磁盘的吞吐极限在100MB左右，随机iops大概在100左右。

在设计时候可以做如下考虑

* 对于通用型存储系统设计来说，需要在latency和带宽利用上有一个折中，即不能使得延迟太高，影响响应时间，也不能使得写入size太小，不然整体带宽出不来。在系统设计的时候选择合适的并发和io size。比如1MB 写入大小，并发控制在4 thread 到8 thread。

* 对于小文件，最好做group commit，比如延迟10ms、20ms以合并吸收更多的小写，这样对延迟影响较小，但是对带宽利用上有很大的收益，从侧面来说又提高了系统的iops能力。 

