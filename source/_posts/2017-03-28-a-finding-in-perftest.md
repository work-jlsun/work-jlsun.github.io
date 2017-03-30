title: 数据存储中Zipf分布
date: 2017-03-28 08:07:30
tags: 存储
---

最近团队在做对存储系统做一些性能测试，期间遇到了不少问题，测试过程得出的数据也没有很好数据支持，所以尝试了非常多的方法来对性能问题进行定位。

小家伙 @王欢明 还是挺厉害的，使用了非常多的工具进行性能问题定位，包括iosnoop对IO请求进行跟踪、iostat进行磁盘状态记录、go-pprof从runtime层面收集性能profile数据、使用go-torch对profile生成直观的火焰图、使用trace2heatmap对延迟数据生产热力图 等等。

纵然是花了很多时间和精力去测试分析，但是某些测试结果具有一定误导性，包含多变量的系统中从外部整体去测试其实很难发现真正原因，走了一些弯路。

所以最后小伙伴通过一定手段对系统中的一些不确定性的环节进行简化确定，真相才慢慢浮出水面。

以下整理分享。

### 性能说明

存储系统中采用多副本技术保障数据可靠性，数据一致性复制协议采用类PacificA协议,leader并发将数据发送到follower(三副本)，所有candidate(leader & follower)完成磁盘写入之后，由leader回复client写入成功。

如下为两个测试数据说明图。

![putdata](http://tompublic.nos-eastchina1.126.net/putdata.jpg)

![fdatasyn](http://tompublic.nos-eastchina1.126.net/fdatasync.jpg)

从上两测试图我们可以看出:

* 写入的整体性能趋势基本是跟磁盘fdatasync的分布呈现相同的分布和趋势。 
* 客户端Put的平均响应时间是fdatasync的2倍左右

*疑问*：leader接收到client写发送过来的数据是并发发送到follower，在低压力情况下，理论上网络上的开销相比磁盘的sync开销基本可以忽略不计，为什么三副本写入结果达到了fdatasync的两倍左右？

后来我们对fdatasync做了mock，响应时间就设置为固定50ms，测试结果发现三副本写入结果基本就是50ms左右。非常符合预期的结果，结果说明了相比于fdatasync程序本身造成的开销基本可以忽略不计，那么现实情况下问题出在哪里？

### zipf分布

从单节点fdatasync的响应时间分布看,是一个典型的[zipf分布](https://zh.wikipedia.org/zh-cn/齊夫定律),大部分请求响应时间较小，而小部分请求响应时间特别大。

所以使用程序拟合了这样类似的分布，并且通过模拟的方式验证了了一个结果

* 从单纯三副本fdatasync来说，并发写三副本的平均响应值差不多为单次fdatasync的2倍左右
* 多数派协议中，比如三副本的系统中，并发写入F／2 +1 = 2 个副本成功情况下 的响应时间比单词fdatasync的平均响应时间好要小。 
	
*拟合程序*
	
	package main
	
	import (
	    "fmt"
	    "math/rand"
	    "sort"
	    "time"
	)
	
	func main() {
	    r := rand.New(rand.NewSource(int64(time.Now().Second())))
	    zipf := rand.NewZipf(r, 2.7, 25, 300)
	
	    data := make([]int, 0)
	
	    N := 20
	    for i := 0; i != N; i++ {
	        item  :=  int(zipf.Uint64() + 30)
	        data = append(data, item)
	    }
	    sort.Ints(data)
	    fmt.Printf("%+v", data)
	}

程序拟合出的曲线如下：

![](http://tompublic.nos-eastchina1.126.net/zipf-simu.png)

三副本提交和两副本quorum情况下的平均提交时间与分布曲线。

![](http://tompublic.nos-eastchina1.126.net/zipf_avg_simu.png)

相关程序详见

### 结论

任何看似奇怪的问题后面有可能隐藏着不为人知的更深层次的原因，执着的专研分析精神,所谓。
