### 1. Raft

Raft协议的发布，对分布式行业是一大福音，虽然在核心协议上基本都是师继Paxos祖师爷(lamport)的精髓，基于多数派的协议。但是Raft一致性协议的贡献在于，定义了可易于实现的一致性协议的事实标准。把一致性协议从“阳春白雪” 变成普通学生、IT码农都可以上手一试的东东，MIT的分布式教学课程都是直接使用Raft来介绍一致性协议。

从《In Search of An Understandable Consensus Algorithm(Extend Version)》论文中，我们可以看到，与其他一致性协议的论文不同的点是，Diego 基本已经算是把一个易于工程实现算法讲得非常明白了，just do it，没有太多争议和发挥的空间，当然要实现一个工业级的靠谱的raft还是要花不少力气。

当然raft一致性协议相对易于实现主要归结为以下几个原因：

1. 模块化的拆分:把一致性协议划分为 Leader选举、MemberShip变更、日志复制、SnapShot等相对比较解耦的模块

2. 设计的简化: 比如不允许类似Paxos算法的乱序提交、使用Randomization 算法设计Leader Election算法以简化系统的状态，只有Leader、Follower、Candidate等等。

本文不打算对Basic Raft一致性协议的具体内容进行说明，而是介绍记录一些关键点，因为绝大部份内容，原文已经说明得很详实，但凡有一定英文技术，直接看raft paper就可以了，如意犹未尽，还可以把raft 作者 Diego Ongaro 200多页的博士论文刷一遍（链接在文末，可自取）。

### 2. Points

#### 2.1 Old Term LogEntry 处理

> 旧Term未提交的日志的提交依赖于新一轮的日志的提交

这个在原文 “5.4.2 Committing entries from previews terms” 有说明，但是在看的时候可能会觉得有点绕。

Raft协议约定，Candidate在使用新的Term进行选举的时候，Candidate能够被选举为Leader的条件为：

1. 得到一半以上(包括自己)节点的投票
2. 得到投票的前提是：Candidate节点的最后一个LogEntry的Epoch比投票节点大，或者在Epoch一样情况下，LogEnry的SN(serial number)必须大于等于投票者。

并且有一个安全截断机制：

1. Follower 在接收到logEntry的时候，如果发现发送者节点当前的Term大于Follower当前的Term；并且发现相同序号的(相同SN)LogEntry在Follower上存在，未Commit，并且LogEntry Term 不一致，那么Follower直接截断从[SN~文件末尾)的所有内容，然后将接收到的LogEntryAppend到截断后的文件末尾。 


在以上条件下，Raft论文列举了一个Corner Case 违反一致性协议，如图所示
 
![raft-corner-case](http://tom.nos-eastchina1.126.net/raft-corner-case.jpg)

* (a): S1 成为 Leader,Append Term2的LogEntry（黄色）到S1、S2 成功;
* (b): S1 Crash, S5使用 Term(3) 成功竞选为 Term(3)的 Leader（通过获得S3、S4、S5的投票），并且将Term为 3的 LogEntry(蓝色) Append到本地;
* (c): S5 Crash, S1 使用 Term(4) 成功竞选为Leader（通过获得S1、S2、S3的投票），将黄色的LogEntry复制到S3，得到多数派响应（S1、S2、S3)的响应，提交黄色LogEntry为Commit，并将Term为4的LogEntry(红色) Append到本地。
* (d) S5 使用新的Term(5) 竞选为Leader(得到 S2、S3、S4 的投票)，按照协议将所有所有节点上的黄色和红色的LogEntry截断覆盖为自己的Term为3 的LogEntry。

进行到这步的时候我们已经发现，黄色的LogEnry(2) 在被设置为Commit之后重新又被否定了。

所以协议又强化了一个限制；

1. **只有当前Term的LogEntry提交条件为：满足多数派响应之后(多数派<一半以上节点>Append LogEntry到日志)设置为commit；**
2. **前一轮Term未Commit的LogEntry的Commit依赖于高轮Term LogEntry的Commit**

如图所示 (c) 状态 Term2的LogEntry（黄色） 只有在 （e）状态 Term4 的LogEntry（红色）被commit才能够提交。

**Old Entry除了会提交还可能会截断**，截断的条件比较单一，Follower在接收到一个已经存在的SN(Serial Number)的未commit的LogEntry（也就是说上图从状态(c)进入到状态(d)），但是Epoch比接收到的相同SN的LogEntry小，那么截断后续未提交的LogEntry。


> 提交NO-OP LogEntry 提交系统可用性

在Leader通过竞选刚刚成为Leader的时候，有一些等待提交的LogEntry(即SN > CommitPt的LogEntry)，有可能是Commit的，也有可能是未Commit的。(PS: 因为在Raft协议中CommitPt 不是实时刷盘的)

所以为了防止出现非线性一致性(Non Linearizable Consistency)；即之前已经响应客户端的已经Commit的请求回退，并且为了避免出现上图中的Corner Case，往往我们需要等到下一个Term的LogEntry的Commit，才能保障提供线性一致性。

但是有可能接下来的客户端的写请求不能及时到达，那么为了保障Leader快速提供读服务，系统可首先发送一个NO-OP LogEntry 来保障快速进入正常可读状态。

#### 2.2 Current Term、VotedFor 持久化

上图其实隐含了一些需要持久化的重要信息，即 Current Term、VotedFor！！！ 为什么(b) 状态 S5 使用的Term Number 为3，而不是2？

因为竞选为Leader就必须是使用新的Term发起选举，并且得到多数派阶段的同意，同意的操作为将Current Term、VotedFor持久化。

比如（a） 状态 S1 为什么能竞选为Leader？首先S1满足成为Leader的条件，S2～S5 都可以接受 S1 成为发起Term 为2 的Leader选举。S2～S5 同意S1成为Leader的操作为：将 Current Term 设置为2、VotedFor 设置为S2 并且持久化，然后返回S1。即S1 成功成为Term 为2的Leader的前提是一个多数派已经记录 Current Term 为2 ，并且VotedFor为S2。那么(b) 状态 S5 如使用Term为2进行Leader选举，必然得不到多数派同意，因为Term 2 已经投给S1，S5只能 将Term++ 使用Term 为3 进行重新发起请求。

> Current Term、VotedFor 如何持久化?

	type CurrentTermAndVotedFor {
		Term int64 `json:"Term"`
		VotedFor int64 `json:"Votedfor"`
		Crc int32
	}

	//current state
	var currentState  CurrentTermAndVotedFor
	
	.. set value and calculate crc ...
	
	content, err := json.Marshal(currentState)
	
	//flush to disk
	f, err := os.Create("/dist/currentState.txt")
	f.Write(content)
	f.Sync()

其实很简单，只需要保存在一个单独的文件，如上为简单的go语言示例；其他简单的方式比如在设计Log File的时候，Log File Header中包含 Current Term 以及VotedFor 的位置。

	如果再深入思考一层，其实这里头有一个疑问？如何保证写了一半（写入一半然后挂了）的问题？写了Term、没写VoteFor？或者只写了Term的高32位？

可以看到磁盘能够保证512 Byte的写入原子性,这个在知乎[事务性(Transactional)存储需要硬件参与吗？](https://www.zhihu.com/question/39142368) 这个问答上就能找到答案。所以最简单的方法是直接写入一个tmpfile，写入完成之后，讲tmpfile mv成CurrentTermAndVotedFor文件，基本可保障更新的原子性。其他方式比如采用Append Entry的方式也可以实现。

#### 2.3 Cluser Membership 变更

在Raft的Paper中，简要说明了一种一次变更多个节点的Cluser Membership变更方式。但是没有给出更多的在Securey以及Avaliable上的更多的说明。

其实现在开源的raft实现一般都不会使用这种方式，比如Etcd raft 都是采用了更佳简洁的一次只能变更一个节点的 “single Cluser MemberShip Change” 算法。

当然single cluser MemberShip 并非Etcd 自创，其实raft 协议作者 Diego 在其博士论文中已经详细介绍了Single Cluser MemberShip Change 机制，包括Security、Avaliable方面的详细说明，并且作者也说明了在实际工程实现过程中更加推荐Single方式，首先因为简单，再则所有的集群变更方式都可以通过Single 一次一个节点的方式达到任何想要的Cluster 状态。

原文：“Raft restrict the types of change that allowed： only one server can be added or removed from the cluster at once. More complex changes in membership are implemented as a series of single-server-change”

#### 2.3.1 Safty

回到问题的第一大核心要点：**Safety**，membership 变更必须保持raft协议的约束：同一时间(同一个Term)只能存在一个有效的Leader。

> 为什么不能直接变更多个节点，直接从Old变为New有问题? for example change from 3 Node to 5 Node？

![multiLeaderInMemberChange](http://tom.nos-eastchina1.126.net/multiLeaderInMemberChange.jpg)

如上图所示，在集群状态变跟过程中，在红色箭头处出现了两个不相交的多数派（Server3、Server4、Server 5 认知到新的5 Node 集群；而1、2 Server的认知还是处在老的3 Node状态）。在网络分区情况下（比如S1、S2 作为一个分区；S3、S4、S5作为一个分区），2个分区分别可以选举产生2个新的Leader（属于configuration< C<sup><sub>old</sub></sup>>的Leader 以及 属于 new configuration < C<sup><sub>new</sub></sup> > 的 Leader ） 。

当然这就导致了Safty没法保证；核心原因是对于C<sup><sub>old</sub></sup> 和 C<sup><sub>New</sub></sup> 不存在交集，不存在一个公共的交集节点 充当仲裁者的角色。

但是如果每次只允许出现一个节点变更(增加 or 减小)，那么C<sup><sub>old</sub></sup> 和 C<sup><sub>New</sub></sup> 总会相交。 如下图所示

![multiLeaderInMemberChange](http://tom.nos-eastchina1.126.net/singleMemberChange.jpg)


> 如何实现Single membership change

以下提几个关键点：

1: 由于Single方式无论如何 C<sup><sub>old</sub></sup> 和 C<sup><sub>New</sub></sup> 都会相交，所以raft采用了直接提交一个特殊的replicated LogEntry的方式来进行 single 集群关系变更。


2: 跟普通的 LogEntry提交的不同点,configuration LogEntry 只需要持久化到server的log就生效，不用等待这个Entry提交。(PS: 原文 "The New configuration takes effect on  each server as soon as it is added to the server's log")


3: 后一轮 MemberShip Change 提交必须在前一轮 MemberShip Change Commit之后进行，以避免出现多个Leader的问题。
 
接下来我们来看通过这三点如何能够保证Safty。
 
 * configuration LogEntry 发送到了server多数派（Cnew 的多数派）

 这种情况下，多数派中肯定会产生一个Leader（当前Leader抑或是由于重启产生的新的Leader）将 configuration LogEntry 提交。（以下为与普通2.1 Old Term LogEntry不太一样，如下所示） （**PS：这里是不正确的** 有一宗case说明）
 ![](http://tom.nos-eastchina1.126.net/configureMarjorRecieved.jpg)
 

 除了以上描述偶数节点新增节点的情况，接下来的5 种case情况都在下图中举例说明
 
 (a). 加节点(奇数变偶数）
 
 (b). 减节点(偶数变奇数&Leader在Cnew)
 
 (c). 减节点(奇数变偶数&Leader在Cnew)
 
 (d). 减节点(偶数变奇数&Leader不在Cnew)
 
 (e). 减节点(奇数变偶数&Leader不在Cnew)
 
 ![CaseOfCnewReceived](http://tom.nos-eastchina1.126.net/CaseOfCnewReceived.jpg)

  
 在所有上述情况下，发现有前一轮 Configuration请求未完成，那么按照上诉3的要求，必须等待前一轮 Configuration请求 commit才处理下一轮 Configuration。所以是的“同一个Term出现2个有限leader”得以保障。
 
 (疑问：configuration LogEntry 是否可以和 普通LogEntry在同一个Baching里头批量发送和提交？)
 (疑问：为什么不需要生效，就可以启用集群关系？why，later)
 
 
* configuration LogEntry 没有发送到多数派（Cnew的多数派）

这种情况下，也就是说Cnew 并未生效。新一个节点被选择为Leader，接下来有可能出下2种情况

a. 有可能被截断取消，即未包含configuration LogEntry的节点竞选为Leader。
b. 重新被提交，即包含configuration LogEntry的节点竞选为Leader，这个接下来的流程跟上述《configuration LogEntry 发送到了server多数派》就一样了，没有必要多展开。

其实 a 截断也没有太多要展开的，回退到Cold 继续运行。

a. 在新增节点情况下，当前Leader是Cnew集群成员, 这个情况下，系统回到 Cold状态，接下来Cold情况下如有成功的提交请求，会导致configuration LogEntry不再有机会提交成功（同一个sn上已经有请求提交），拥有此configuration LogEntry 的节点想加入集群，必然要截断。（PS：当然Cold 接下来还可能会遇到configuration LogEntry 请求，这种情况下新老configuration  logEntry 请求 也存在一个 majority 交集，使得只有一个能够成功）


接下来2种case是类似的分析方法，都能够确保在一个term里头出现2个有效leader

b. 在减去节点的情况下，当前Leader节点是Cnew集群成员。

c. 在减去节点的情况下，当前Leader节点非Cnew集群成员。



>  2个 奇怪的节点



>  可用性相关

1: CacheUp

2: Removing the Current Leader

3 Disrptive serves


4 可用性论证




3. 拓展点：刚开始是如何形成一个初始的Cluster MemberShip的？





#### 2.4 线性一致性




### 3 参考文献


1. [Raft Paper](https://raft.github.io/raft.pdf)
2. [Raft 博士论文](https://web.stanford.edu/~ouster/cgi-bin/papers/OngaroPhD.pdf)
3. [事务性(Transactional)存储需要硬件参与吗？](https://www.zhihu.com/question/39142368)
