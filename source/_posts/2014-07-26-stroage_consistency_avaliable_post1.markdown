---
title: ' 存储系统一致性与可用性'
layout: post
tags:
    - storage
    - consistency
    - avaliable
---

### 1. 关于可用性

提供系统可用性的关键是在相关组件失效情况下，系统能多快恢复并继续正确得提供服务。

#### 1.1 如何恢复服务

How to recover a single node from power failure。

* wait for reboot

Data is durable, but service is unavailable temporarily。

* Use multiple nodes to provide service

Another node takes over to provide service, How to make sure nodes respond in the same way?

对于无状态的服务器而言，没有关系。对于有状态的服务器而言,如何去保证 take over的服务器提供正确的服务?

### 1.2 两种基本方法

* 多副本
* 纠删码
纠删码基本原理是把一份数据分割成多份，计算出若干冗余块，比如 切割成26份数据再使用纠删算法计算4份冗余块，可以使得任意丢4块数据还能恢复数据进行服务，核心还是通过数据冗余来实现一定的高可用。以下内容只针对简单多副本展开讨论，本文不对纠删码进行详述。

对于简单多副本要保证多个副本的状态是一致的，才能保证另一个副本对应服务器能提供正确的服务。如何实现多副本之间的一致性?以下讨论相关理论及实现.

## 2. 理论

实现多副本一致的一个基本理论为 replicated stat machine，即副本状态机：副本所处的初始状态相同；每个副本执行的操作顺序相同，并且每一个操作都相同； 那么最终所有副本的状态是一致。

![OpaaOPb](/media/files/2014/07/OPaOPb.png)

那如何保证在多个 client 并发操作下的保证操作的顺序性？以下分析 primary-backup协议，Qurom，paoxs 协议等协议




### 2.1 Master Slave

Master / Slave 相对来讲是较简单和自然而然的方式，Master 决定操作的顺序，Slave 节点执行序列化之后的操作。

![Primary.png](/media/files/2014/07/Primary.png)


Master Slave 协议较为简单， 但是其可用性不是很高， 在其中某一个节点宕机的情况下，系统无法继续运行。

如下图所示，在 state0 状态 Master 与 Slave 的 LSN(log sequence number)为 10，处于一致的状态，在 state1 时 Slave 发生宕机，Master 继续接受写请求，到 state2 时 Master 也生宕机，在 state3 时 Slave 恢复，但是 Master 与 Slave 的状态不一致，所以 Slave 即不能提供读也不能提供写， 此时即发生阻塞。 若想继续提供服务， 必须进行人工干预。 Master/slaver
协议是一种阻塞式的协议。 （数据库中的 “两阶段提交协议” 就是一种典型的阻塞式的协议）

![3.png](/media/files/2014/07/3.png)


### 2.2  Quroum/(NRW)

Quorum 机制是一种简单有效的副本管理机制。NWR 协议是其中的典型，NWR 为Amazon 公司设计的分布式 kv 系统 dynamo 中使用的分布式副本协议。

其中 N 为副本的数目， W 是写成功需要写的份数， R 为读成功需要读的份数， 保证 R+W>N，即能读到正确的数据。

假设 N=3，W=1，R=3，即如果副本数目为 3，写 1 份即算成功，那么至少读 3 份才能读到写入的数据。 并且为了防止写丢失， 还必须用 “last write wins” 的方式解决写冲突问题。

NRW 可能出现数据丢失情况

如下图 所示， 例如用户先调用 OP1 写数据 obj1 的 A 副本成功返回， 然后用户调用 OP2写数据 obj1 的副本 C 成功后返回，此时使用”last write wins”的进行合并使得最终副本 A 和C 上的数据完全一致。但是如果 A 与 B 上的时钟不一致，比如 A 的时钟比 C 快，那么在数据整合的时候会出现先写的数据副本 A 上的操作覆盖后写的数据 C， 导致后写的 OP2 丢失，这显然不是用户想看到的信息。所以系统必须排除这种状况才能保证数据的最终一致性， 否则就是弱一致性。

![qurom.png](/media/files/2014/07/qurom.png)


时钟同步：在分布式系统中本身就是很难事件的东西。

结论:NRW 其实也只是看上去很美的东西。但是 dynamo 为了使得系统能够有更好的可扩展性及可用性， 放宽对一致性的要求， 不支持强一致性， 而支持最终一致性甚至是弱一致性。 （弱一致性估计是很难使用户接受）


### 2.3 paxos协议

paxos 协议中每个节点都是对等的，每个节点都可以接收客户端的写请求，并尝试完成客户端的写请求，其核心是基于一个多数派的抢占机制式协议。

![basicPaxos.png](/media/files/2014/07/basicPaxos.png)


### 2.4 关于CAP

各种分布式协议所面临的问题就是分布式文件系统的三大难题:副本一致性，系统可用性，分区容忍性（网络异常的容忍能力）。

CAP理论明确提出了不要妄图设计一种对 CAP 三大属性都完全拥有的系统，因为这种系统在理论上就已经被证明不存在。设计系统的时候需要在 C、A、P 这三方面有所折中。

Primary/backup：MySQL：具有完全的 C，很糟糕的 A，很糟糕的 P， （任何一个节点宕机都会导致服务中断）
Quroum 协议：Dynamo/cassandra，有一定的 C，有较好的 A，也有较好的 P，是一种较为平衡的分布式协议。
Paxos 协议：Spanner , Chubbuy, Zookeeper, megastore 具有完全的 C，较好的 A， 较好的 P。
发现很多开源系统在高可用、性能、可扩展性等很多方面没有一个完善的方案。希望设计一个 k/v 存储系统，系统支持强一致性，并且在 C，A，P 这三方面都较为均衡，系统具有一定可扩展性，详见下文。




