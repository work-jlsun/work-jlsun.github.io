title: HTTP之X-Forwarded-For
date: 2016-03-11 07:08:10
tags:
- HTTP
- Nginx
---

### X-Forwarded-For?

HTTP协议是允许代理模式的，代理服务器的功能很多，比如

* HTTP请求的负载均衡
* HTTP资源的缓存服务
* HTTP请求的权限控制
* 等等

这些功能在当前都有非常多的实际应用。比如CDN厂商的资源分发缓存服务、Web应用服务器的前端反向代理、公司内部统一HTTP请求的出口控制等等。

代理服务器可以拿到HTTP请求并对请求进行修改，这导致到达源站的请求不一定是原始的client发出的请求；如果client请求中间经过了某些不可信的代理服务器，这会导致最终收到的请求很大程度上有一定不可信。

撇开代理服务器对请求HTTP请求的修改，另外一点值得注意的是，在代理模式下，会导致client IP信息丢失。因为client建立的端到端TCP连接是client与HTTP代理服务器之间的，而不是client到origin源站服务器之间的；如何才能保留这个信息，同时trace中间经过的代理服务器？HTTP设计了一标准的协议头部，即X-Forwarded-For。

X-Forwarded-For是HTTP超文本协议中的一个标准头部，详见[RFC](https://tools.ietf.org/html/rfc7239)。HTTP协议X-Forwarded-For头部的目的是为了保留client信息并trace从client到源站服务器之间所有途经代理的IP信息。其基本格式为`X-ForWarded-For:clientip, proxy1, proxy2`。


![](/media/files/2016/03/x-forwarded-for.jpg)

如上图所示，客户端将请求发送到代理1，代理1按照协议要求，在头部中添加客户端的Ip信息`X-Forwarded-For:client-ip1`并把请求转发到代理2，代理2同样也按照协议要求，在头部中添加代理1的IP信息`X-Forwarded-For:client-ip1, proxy-ip1`,所以最后源站收到请求为蓝色框中的请求。



### 一点实践 

为了提升海外客户的用户体验(提高请求的成功率和响应速度)，最最简单的方式为租用一些海外专线代理，然后将海外请求解析到专线代理，由代理转发国内。

![](/media/files/2016/03/overseaproxy.jpg)

但是我们有个小奇葩的需求，需要得知用户所处的不可控网络中的最后一个外网IP,所以在web service中有一个逻辑是取X-Forwarded-For头部中的最后一个外网Ip。**在没有海外代理之前一切正常运行，但是在部署海外代理之后，由于海外代理与源站之间是通过外网IP进行转发到，这就导致了X-Forwarded-For头部最后一个外网IP变成海外代理的外网IP**，最直观的当然是修改web service 逻辑去除海外代理的IP；但是使用nignx本身的配置逻辑可以更加方便简洁得解决，如下：

```
location /{
	set $XForwadFor $proxy_add_x_forwarded_for;
	if  ($remote_addr = "海外代理IP") {
		##set the X-Forward-for header
		set $XForwadFor $http_x_forwarded_for;
	}
	proxy_set_header X-Forwarded-For $XForwadFor;
	proxy_pass http://webservice;
}
```

	