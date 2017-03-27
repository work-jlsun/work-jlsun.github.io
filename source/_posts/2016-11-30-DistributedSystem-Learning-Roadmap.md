title: Distributed System Learning Plan
date: 2016-11-30 17:44:25
tags:
---

[10]: http://research.microsoft.com/users/lamport/pubs/lamport-paxos.pdf
[11]: http://research.microsoft.com/users/lamport/pubs/paxos-simple.pdf
[12]: http://research.google.com/archive/paxos_made_live.html
[13]: http://scholar.google.com/scholar?q=Paxos+Made+Practical
[14]: https://www.youtube.com/watch?v=JEpsBg0AO6o
[15]: https://www.youtube.com/watch?v=YbZ3zDzDnrw
[16]: http://thesecretlivesofdata.com/raft/
[17]: https://www.microsoft.com/en-us/research/wp-content/uploads/2008/02/tr-2008-25.pdf
为了提高团队成员在分布式系统方面的专业实力，制定了一份学习计划，希望通过执行本计划使得成员能够从0开始更佳全面的了解分布式系统的基本理论，做到知其然知其所以然，更深入得理解分布式系统设计的关键点，从而更好的指导工程实践。

### 1. Introduction


学习目标：学习如何阅读一篇论文，了解分布式系统的基本概念

 1. scalability
 2. availability
 3. performance
 4. latency
 5. fault tolerance
 6. and so on

参考文献：

* [How To Read An Engineering Research Paper](https://cseweb.ucsd.edu/~wgg/CSE210/howtoread.html)
* [Distributed systems at a high level](http://book.mixu.net/distsys/intro.html)

 
 
### 2. Up and down the level of abstraction

学习目标：了解分布式系统的基本问题和理论。

1. Meaning Of Abstraction
2. System Model
3. 分布式系统中的运行单元
4. 什么是网络分区 
5. 同步 vs. 异步
6. 一致性问题
7. CAP 、FLP 、ACID 理论
8. 一致性模型

参考文献：

* [Up and down the level of abstraction](http://book.mixu.net/distsys/abstractions.html)
* [Brewer’s Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web.pdf](http://tom.nos-eastchina1.126.net/Brewer%E2%80%99s%20Conjecture%20and%20the%20Feasibility%20of%20Consistent%2C%20Available%2C%20Partition-Tolerant%20Web.pdf)
* [Impossibility of Distributed Consensus with One Faulty Process.pdf](http://tom.nos-eastchina1.126.net/Impossibility%20of%20Distributed%20Consensus%20with%20One%20Faulty%20Process.pdf)
* [FLP Impossibility的证明](http://danielw.cn/FLP-proof)
* [分布式系统工程实践－>2.3/2.4](http://tom.nos-eastchina1.126.net/%E5%88%86%E5%B8%83%E5%BC%8F%E7%B3%BB%E7%BB%9F%E5%B7%A5%E7%A8%8B.pdf)

### 3. Time And Order

学习目标：

1. Total and partial order
2. 系统中的时间time
3. "Global-clock"、"Local-clock"、"no-clock"
4. Vector clock  in  Detail

参考文献：

* [Time And Order](http://book.mixu.net/distsys/time.html)
* [Time, Clocks, and the
Ordering of Events in
a Distributed System ](http://tom.nos-eastchina1.126.net/p558-lamport.pdf)

### 4. Replica And Consensus Problem

学习目标：了解和学习副本一致性,比如Primary/Backup、2PC 以及 分布式共识算法。

1. Replication: Syn & Async
2. Primary/Backup & 2PC 
3. Partition tolerant consensus algorithms
4. Algorithms Examples:Raft

参考文献：

* [Distributed systems->4、5](http://book.mixu.net/distsys/index.html)
* [Raft lecture (Raft user study)][15]
* [Raft Understandable Distributed Consensus][16]

扩展阅读：Partition-tolerant consensus algorithms: PacificA、Paxos

* [Paxos lecture (Raft user study)][14]
* [The Part-Time Parliament][10] 
* [Paxos Made Simple][11] 
* [Paxos Made Live - An Engineering Perspective][12] 
* [Paxos Made Practical][13] 


### 5.  Data Distribution(Replica PlaceMent)

学习目标：分布式存储系统扩展新很重要的一方面是数据(副本)的划分和放置，这里需要学习基本的划分方式。

1. Range :字典序拆分、List
2. Hash：Consistent hash、DHT、CRUSH

参考文献：

* [分布式系统设计白皮书->可扩展性](http://doc.hz.netease.com/pages/worddav/preview.action?fileName=%E5%8F%AF%E6%89%A9%E5%B1%95%E6%80%A7.docx&pageId=36447308)
* [分布式系统工程->2.1 数据分布方式](http://tom.nos-eastchina1.126.net/%E5%88%86%E5%B8%83%E5%BC%8F%E7%B3%BB%E7%BB%9F%E5%8E%9F%E7%90%86%E4%BB%8B%E7%BB%8D.pdf)
* [CRUSH](http://tom.nos-eastchina1.126.net/weil-crush-sc06.pdf)


### 6. Distributed System Example

分布式存储泛指存储存储和管理数据的系统， 与无状态的应用服务器不同， 如何处理各种故障以保证数据一致，数据不丢， 数据持续可用， 是分布式存储系统的核心问题，也是极具挑战的问题。 如下总结了分布式存储领域的经典学习论文。（参考：[distributed-storage-papers](http://www.bitstech.net/2013/12/27/distributed-storage-papers/)）

0. **Ceph: Reliable, Scalable, and High-Performance Distributed Storage. Sage A. Weil.** 功能强大的开源海量存储系统， 支持文件系统、块设备、以及S3接口。 主要技术特色： CRUSH数据对象定位算法， 基于动态子树的文件系统元数据管理。

1. **The Google File System.** Sanjay Ghemawat, Howard Gobioff, and Shun-Tak Leung。 基于普通服务器构建超大规模文件系统的典型案例，主要面向大文件和批处理系统， 设计简单而实用。 GFS是google的重要基础设施， 大数据的基石， 也是Hadoop HDFS的参考对象。 主要技术特点包括： 假设硬件故障是常态（容错能力强）， 64MB大块， 单Master设计，Lease/链式复制， 支持追加写不支持随机写。

2. **Bigtable: A Distributed Storage System for Structured Data.**Fay Chang, Jeffrey Dean, Sanjay Ghemawat, et. 支持PB数据量级的多维非关系型大表， 在google内部应用广泛，大数据的奠基作品之一 ， Hbase就是参考BigTable设计。 Bigtable的主要技术特点包括： 基于GFS实现数据高可靠， 使用非原地更新技术（LSM树）实现数据修改， 通过range分区并实现自动伸缩等。

3. **Spanner: Google’s Globally-Distributed Database.** James C. Corbett, Jeffrey Dean, et. 第一个用于线上产品的大规模、高可用， 跨数据中心且支持事务的分布式数据库。 主要技术特点包括， 基于GPS和原子钟的全球同步时间机制TrueTime， Paxo， 多版本事务等。

4. **PacificA: Replication in Log-Based Distributed Storage Systems.** Wei Lin, Mao Yang, et. 面向log-based存储的强一致的主从复制协议， 具有较强实用性。 这篇文章系统地讲述了主从复制系统应该考虑的问题， 能加深对主从强一致复制的理解程度。 技术特点： 支持强一致主从复制协议， 允许多种存储实现， 分布式的故障检测/Lease/集群成员管理方法。

5. **Object Storage on CRAQ, High-throughput chain replication for read-mostly workloads.** Jeff Terrace and Michael J. Freedman. 支持强一直的链式复制方法， 支持从多个副本读取数据。

6. **Finding a needle in Haystack: Facebook’s photo storage.** Doug Beaver, Sanjeev Kumar, Harry C. Li, Jason Sobel, Peter Vajgel. Facebook分布式Blob存储， 主要用于存储图片。 主要技术特色： 小文件合并成大文件， 小文件元数据放在内存因此读写只需一次IO。

7. **Windows Azure Storage: A Highly Available Cloud Storage Service with Strong Consistency.** Brad Calder, Ju Wang, Aaron Ogus, Niranjan Nilakantan, et. 微软的分布式存储平台， 除了支持类S3对象存储，还支持表格、队列等数据模型。 主要技术特点： 采用Stream/Partition两层设计（类似BigTable）；写错（写满）就封存Extent， 使得副本字节一致， 简化了选主和恢复操作； 将S3对象存储、表格、队列、块设备等融入到统一的底层存储架构中。

8. **The Chubby lock service for loosely-coupled distributed systems.** Mike Burrows. Google设计的高可用、可靠的分布式锁服务， 可用于实现选主、分布式锁等功能， 是ZooKeeper的原型。 主要技术特点： 将paxo协议封装成文件系统接口， 高可用、高可靠，但是不保证有很强性能。

9. **Paxos Made Live – An Engineering Perspective.** Tushar Chandra, Robert Griesemer，Joshua Redstone. 从工程实现角度说明了Paxos在chubby系统的应用， 是理解Paxo协议及其应用场景的必备论文。 主要技术特点： paxo协议， replicated log， multi-paxos。

10. **Dynamo: Amazon’s Highly Available Key-Value Store.** Giuseppe DeCandia, Deniz Hastorun, Madan Jampani, et. Amazon设计的高可用的kv系统， 主要技术特点：综和运用一致性哈希，vector clock， 最终一致性构建一个高可用的kv系统， 可应用于amazon购物车场景。



### 7. Others

* [分布式系统工程](http://tom.nos-eastchina1.126.net/%E5%88%86%E5%B8%83%E5%BC%8F%E7%B3%BB%E7%BB%9F%E5%B7%A5%E7%A8%8B.pdf)
* [分布式系统原理介绍](http://tom.nos-eastchina1.126.net/%E5%88%86%E5%B8%83%E5%BC%8F%E7%B3%BB%E7%BB%9F%E5%8E%9F%E7%90%86%E4%BB%8B%E7%BB%8D.pdf)
* [存储系统一致性与可用性](https://work-jlsun.github.io/2014/07/26/stroage_consistency_avaliable_post1.html)
* [分布式系统设计白皮书](http://doc.hz.netease.com/pages/viewpage.action?pageId=35319682)(PS：网易内部资料)




