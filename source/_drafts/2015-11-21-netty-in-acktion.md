title: netty_in_acktion
date: 2015-11-21 14:42:50
categories: 技术
tags:
---

# Netty In Action

Netty是当前Java编程中使用得最多的non-block IO 框架，在内部技术选型过程中，对非block方式的IO框架做过一定的调研。

在内部使用tomcat组赛IO出现过一些问题

* 后端个别服务堵住导致线程永远，所有服务开始


之前的解决方案

nginx ＋ tomcat ；nginx抗链接＋ IO

* Disk IO 问题
* 流量放大问题

等一些在系统压力大的情况下极容易导致服务不稳定。解决问题的方案
* 业务拆分

调研解决方案，netty、or go


# 相关参考

netty in action 

