title: Split Data With Salt
date: 2016-12-28 07:21:12
categories: code
tags:
    - code
    - 架构
---

前一阵子团队来一新人，分享了其所从事物联网相关领域的事件经验。其主要核心逻辑是收集处理来自海量的物联网设备上报的一些信息。比如大楼电路、温度、湿度等一些周期性的上报信息，记录这些信息的方式为使用HBase。其中一点是使用Salt 来避免数据分布的不均匀特性。



	|  EntityID | TimeStamp  | Temp | Current| humidity|
	| --- | --- |--- |--- |---|
	| 1 | time1 | 10| 5|20| 
	| 1 | time2 | 10| 5|20| 
	| 2 | time1 | 10| 5|20| 
	| 2 | time2 | 10| 5|20| 
	| 2 | time3 | 10| 5|20| 
	| 2 | time4 | 10| 5|20| 


一般来说Key的设计是EntityID + TimeStamp；但是现实中EntityKey划分相对比较集中，比如大企业一般来说会划分一大段连续的EntityKey，其产生的数据就非常集中，很容易导致查询热点(多个Entity之间的联合需求基本很少，一般都是根据EntityID来查询数据)。

这个时候 可以**给数据加点盐**，根据EntityID换算一个固定的Salt(**盐**)作为Key 的一部分，比如Hash(EntityKey)。

	|  Salt| EntityKey | TimeStamp  | Temp | Current| humidity|
	| --- | --- |--- |--- |---|---|
	|100100| 1 | time1 | 10| 5|20| 
	|100100| 1 | time2 | 10| 5|20| 
	|100120| 2 | time1 | 10| 5|20| 
	|100120| 2 | time2 | 10| 5|20| 
	|100120| 2 | time3 | 10| 5|20| 
	|100120| 2 | time4 | 10| 5|20| 


INALL：其实加盐技术是存储方面split数据还算是蛮通用的一个引用层次的优化，做后台开发有时候很想通过通用的手段解决一些棘手的问题，其实很多时候，换个角度，从应用场景出发，反而能够事半而功倍。