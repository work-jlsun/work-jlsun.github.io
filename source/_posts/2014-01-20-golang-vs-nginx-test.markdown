---
title: 'Nginx And Go Http 并发性能'
date: 2014-01-20 15:27:57
layout: post
tags:
    - golang
    - nginx
    - performance
---



## 1 测试硬件

```
Intel(R) Xeon(R) CPU           X3440  @ 2.53GHz
cpu cache size	: 8192 KB
DRAM：8G
```

## 2 测试软件

```
2.6.32-5-amd64 #1 SMP
nginx：ngx_openresty-1.4.3.4 
go :go version go1.3.3 linux/amd64
ab：ApacheBench, Version 2.3 
```

## 3 测试配置
---------

### 3.1  一些内核配置

```
/proc/sys/fs/file-max                    3145728
/proc/sys/fs/nr_open                     1048576
/proc/sys/net/core/netdev_max_backlog    1000
/proc/sys/net/core/rmem_max              131071
/proc/sys/net/core/wmem_max              131071
/proc/sys/net/core/somaxconn             128
/proc/sys/net/ipv4/ip_forward            0
/proc/sys/net/ipv4/ip_local_port_range   8192	65535
/proc/sys/net/ipv4/tcp_fin_timeout       60
/proc/sys/net/ipv4/tcp_keepalive_time    7200
/proc/sys/net/ipv4/tcp_max_syn_backlog   2048
/proc/sys/net/ipv4/tcp_max_tw_buckets    1048576
/proc/sys/net/ipv4/tcp_no_metrics_save   0
/proc/sys/net/ipv4/tcp_syn_retries       5
/proc/sys/net/ipv4/tcp_synack_retries    5
/proc/sys/net/ipv4/tcp_tw_recycle        0
/proc/sys/net/ipv4/tcp_tw_reuse          0
/proc/sys/vm/min_free_kbytes             11489
/proc/sys/vm/overcommit_memory           0
```

### 3.2 Nginx

```
worker_processes  8;
events {
    worker_connections  2046;
    use epoll;
}
http {
    include       mime.types;
    default_type  application/octet-stream;
    location / {
            root   html;
            index  index.html index.htm;
    }
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
            root   html;
    }
    location /ab-test {
        	proxy_http_version 1.1;
        	proxy_set_header Connection ""; 
        	content_by_lua 'ngx.print("aaa---here omit other char a, total 512--- aaaaaa")';    
    }
}
```

### 3.3 golang 测试代码

[go512server.go][3]


## 4 测试方法

### 4.1 测试工具及命令

使用 ab测试不同并发场景下nginx和golang http 服务的性能，测试数据大小512Byte。

测试命令示例：ab -n 1000000 -c 5000  -k  "http://127.0.0.1:8081/512b"

(ps: 所有测试结果，都是3次之后取平均值)


### 4.2  测试结果


* 短连接场景

并发请求量 | 100 | 200 | 500 | 1000 | 2000 | 5000
----|-----|----|----|------|----
nginx(tps) | 12741.62   | 12598.08     |  11917.15   |  12016.63 | 11640.36   |   6047.29
go（tps）  |  11310.32   | 11208.87     | 10731.40    |  10757.3  | 10750.26   |     10869.80

ps： 端连接情况下 并发5000 情况下， nginx情况不知道是为什么（nginx进程cpu利用看起来不是很均衡）


* 长连接场景

并发请求量 | 100 | 200 | 500 | 1000 | 2000 | 5000
----|-----|----|----|------|----
nginx（tps）    |     61249.81   |     60672.71  |   59548.39  |   55287.55    |           58375.65    |     60662.44
go（tps）    |    55257.64       |     53288.23 |   49006.64  |  46362.55  |             48042.18      |    47855.02


* golang + nginx （golang as proxy）

并发请求量 | 100 | 200 | 500 | 1000 | 2000 | 5000
----|-----|----|----|------|----
go（tps）  |   31535.37   | 29081.96  |  30250.24  |   28921.48  |  26631.12  |    25333.64

1: [golang  proxy 代码][3]

2: golang作为proxy的时候性能基本为非proxy的一半左右，这个是可以理解的，一个请求的响应时间就是nginx + go两层的响应时间。

3: 使用golang自带的httpclient连接后端的nginx


* nginx + golang （nginx as proxy）

并发请求量 | 100 | 200 | 500 | 1000 | 2000 | 5000
----|-----|----|----|------|----
TPS | 43336.19    | 41722.05   | 37984.94  |  34033.42  |  29489.74  |         25693.03


* golang proxy简单稳定性测试

```
5000 并发测试 30 分钟，tps 达到24751，基本没有因为随着时间的增长而对性能造成很大的影响，资源使用也比较稳定

1000并发测试 1小时 ， tsp达到 27651.10，基本没有因为随着时间的增长而对性能造成很大的影响，资源使用比较稳定
```

## 5 测试结论


基于以上4中场景，上百组测试下，得出一下简单结论

* golang表现还是较为出色，相比于标杆Nginx性能差20%左右
* golang在高并发压力测试下稳定性还是不错，可以接受的
  
当前我们基于简单测试环境下的测试验证golang，现实环境远远比测试环境复杂，后续我们会在NosMedia 开发测试上线过程中不断总结经验。

## 6 参考资料

* [golang test][1] 
* [nginx-lua vs golang][2]


[1]: https://gist.github.com/hgfischer/7965620
[2]: http://blog.lifeibo.com/blog/2013/01/28/ngx-lua-and-go.html
[3]: https://github.com/work-jlsun/golang/blob/develop/go512server.go
[4]: https://github.com/work-jlsun/golang/blob/develop/goproxytest.go


## 7 坑
1. goang http 长连接问题
2.  Connection reset  by peer (104)
