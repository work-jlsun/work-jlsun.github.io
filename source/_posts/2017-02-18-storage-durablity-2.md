title: 分布式存储系统可靠性-系统估算示例
date: 2017-02-18 08:30:02
tags:
---


#### 1 估算示例

上文[分布式存储系统可靠性-如何估算](https://work-jlsun.github.io/2017/01/24/storage-durablity.html)中，我们提供了一些基本的估算的方法。接下来我们提供一个具体的估算的示例子。
	
系统示例: N = 7200块磁盘的存储系统中，R=3副本，单盘容量 DiskSize = 8T,磁盘年平均故障率 AFR = 4%，系统采用大多数存储系统采用的小文件合并成大文件分片的方式存储，分片大小为PartSize = 10GB，系统整体空间利用率Percent =70%(系统总体承载数据量大约在13PB)，分片采用随机放置方式,随机放置情况下 S = Min(C(R,N), (N\*DiskSize*Percent/R)/T) 。



**以下为计算恢复时间 1 小时情况下的概率的过程**

```
	S =  Min( C(7200,3), 7200 * 8 * 1024 * 70% / 3 / 10) = 1376256
	
	Pa(T,K) = P(1,K) = ( X / C(7200,K) ) * P( N(1) = K)
	
	λ = 4% * 7200 /365/24 = 0.03 (按照4%的坏盘率计算999块盘的系统中平均每小时的坏盘数量)
	
	P( N(1) = K) = (λ**K) *  (e ** -λn)/K!
	
	Pb(T) = Σ P(1,K) ; K∈[3, 7200] 
	
	Pc = 1 - (1-Pb(T))**(365*24/1)
```


可以看到其中 选K个节点情况下，命中copyset的概率为 X/ C(K, 7200)。 这个X的计算我们暂时还没有确定，我们看下这个如何计算

K = 3 的情况下，即Pa（1，3）= S/C(7200,3) \* P( N(1) = 3) 比较简单，但是在K∈[4, 7200]点情况下，情况就比较复杂了，貌似没有一种比较好的统一的组合概率计算这个概率。如果非要计算基本可以使用 *蒙特卡罗* 方法case by case 进行计算，该方法详见附录。如果我们在分析概率的基础上进行进一步简化估算。

通过计算P(N(1),K)  可以计算得到1小时内坏K个的概率

K | P(N(1), K)| P(1,K)
---|---|---
3 | 5.73e-6 | 1.26e-10
4 | 4.71e-8 | ?
5 | 3.09e-10| ?
6 | 1.69e-12| ?
7 | 7.97e-15| ?
8 | 3.27e-17| ?
9 | 1.19e-19| ?

从上面我们可以看到，在K >= 6情况下，假设选中copyset的概率为1，对结果的影响也是比k=3小2个数量级以上。所以基本只需要统计考虑K = 4，K=5 情况下的即可。K=4，5 情况下的选中CopySet的概率 基本为 C(S,1) * C(N-3, K-3) / C(N,K)。

K | P(N(1), K)| P(1,K)
---|---| ---
3 | 5.73e-6 | 1.26e-10
4 | 4.71e-8 | 4.17e-12
5 | 3.09e-10| 6.85e-14

 
P(T) ~=  Σ P(1,K) ; K∈[3, 5] = 1.31e-10 (即t=1 小时内 丢数据的概率为1.31e-1)

P = P = 1 - (1-P(T))\*\*(365*24/T) = 1.1e-06， 即6个9的可靠性。


#### 2 附录
####  2.1 估算代码

```
	#!/usr/bin/python
	# -*- coding: utf-8 -*-
	import decimal
	import math
	import time
	
	# 随机分布情况下系统的copyset组合数
	def RandomCopySets(N, DiskSize, RepNum, Percent, PartSize):
	    setNum = (N * DiskSize * Percent / (RepNum *  PartSize))
	    MaxCopySetNum = C(N, RepNum)
	
	    copysetNum = 0
	    if setNum > MaxCopySetNum:
	        copysetNum = MaxCopySetNum
	    else:
	        copysetNum =  setNum
	
	    return int(copysetNum)
	
	# N 个磁盘存储系统中T时间同时损坏K块盘的概率,年故障率ARF
	def  KdiskFailRate(N, T, ARF, K):
	    # λ 每小时的换盘数量
	    lambda1 = decimal.Decimal(str(N*AFR/24/365))
	    return poisson(lambda1, T, K)
	
	# 副本数R的N 个磁盘存储系统中T时间内造成数据丢失的概率, 只统计R -> 2R-1个副本情况下的丢失数据概率(大于R个情况下，在一遍情况下对结果影响比较小)
	def LossDataInT(S, N, RepNum, T, ARF):
	    loosRate = decimal.Decimal(str(0.0))
	    for k in range(RepNum, RepNum*2):
	        kdrate = KdiskFailRate(N,T,ARF,k)
	
	        singlerate = S * C(N-3, k-3)/C(N,k)
	
	        kdlossrate = kdrate * singlerate
	
	        print "k = " + str(k)  + ", " +str(kdrate) + ", " + str(kdlossrate)
	
	        loosRate += kdlossrate
	    print loosRate
	    return loosRate
	
	# define loseRate in one Year
	def LoseRate(S, N, RepNum, T, AFR):
	    return 1 - (1 - LossDataInT(S, N, RepNum, T, AFR))**(365*24/T)
	
	#组合运算
	def C(n, m):
	  return factorial(n) / (factorial(m)*factorial(n-m))
	
	#泊松分布
	def poisson(lam, t, R):
	    e=decimal.Decimal(str(math.e))
	    return ((lam * t) ** R) * (e**(-lam*t)) / factorial(R)
	
	#t时间内损坏R块磁盘的概率
	def probability(t, R):
	  return poisson(t, R)
	
	#级数
	def factorial(n):
	  S = decimal.Decimal("1")
	  for i in range(1, n+1):
	    N = decimal.Decimal(str(i))
	    S = S*N
	  return S
	
	
	
	# case 1
	N = 7200
	DiskSize = 8*1024
	Percent = 0.7
	PartSize = 10
	RepNum = 3
	T = 1
	AFR = 0.04
	
	S  =  RandomCopySets(N, DiskSize, RepNum, Percent, PartSize)
	
	print LoseRate(S, N, RepNum, T, AFR)
```

####  2.2 蒙特卡罗

[蒙特卡罗方法（Monte Carlo Method）](https://en.wikipedia.org/wiki/Monte_Carlo_method)。是一种计算方法。原理是通过大量随机样本，去了解一个系统，进而得到所要计算的值。
它非常强大和灵活，又相当简单易懂，很容易实现。对于许多问题来说，它往往是最简单的计算方法，有时甚至是唯一可行的方法。 简单介绍可参见[蒙特卡罗方法入门](http://www.ruanyifeng.com/blog/2015/07/monte-carlo-method.html)。

在解这里的 X / C(7200, K) 的计算基本可以转换为不直接求X，而使用足够多的计算机随机试验来直接计算概率。 首先构建一个含S个随机的copyset
组合，然后从7200个节点中随机选择K个节点，记录每一次是否命中copyset，实验次数越多，概率越准确。

