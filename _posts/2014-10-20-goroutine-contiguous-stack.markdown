---
title: 'goroutine contiguous stack'
layout: post
tags:
    - golang
    - performance
---

在我们学习关于golang goroutine的文章时，或多或少很多类似的断言：相比于linux pthread，在 golang中我们可以很轻松的创建100k＋的goroutine，而不用担心其带来的开销，其中一个原因是goroutine初始stack非常小，在当前release的1.3 版本中，一个goroutine初试创建只需要4K 的stack，而linux pthead 则需要2M或者更多的stack空间，那到底是不是这样的？

如下在一个进程中创建100个线程，主进程和线程sleep的方式简单测试下pthread 线程初试创建占用的内存资源。

    void* func(void* arg){ 
        while (true){ 
             usleep(1000000); 
        } 
    }
    int main(int argc, char* argv[]) {
        pthread_t tid; 
        int n = 10000; 
        while (n != 0) { 
             if (pthread_create(&tid, NULL, func, &n) != 0 ) {
                 printf("create fail"); 
             } 
             n = n - 1; 
        } 
        while (1) { usleep(10000000); } 
     } 
    

通过ps看到的资源使用情况如下

    USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND 
    12768    21828  4.9  1.0 81976776 83988 pts/3  Sl+  08:08   0:01 ./a.out
    

可以看到创建10000个线程之后，该程序实际占用的内存（RSS）为83988KB，算下来每个线程所占用的内存空间才8K左右，远远不是很多文章中所说的1M或者时2M or more。

## pthread 线程堆栈

通过strace分析phtread_create 得到如下结果

    mmap(NULL, 8392704, 
        PROT_READ|PROT_WRITE,MAP_PRIVATE|MAP_ANONYMOUS|MAP_STACK, 
        -1, 0) = 0x7fa274760000
    brk(0)                                  = 0x1fd2000
    brk(0x1ff3000)                          = 0x1ff3000
    mprotect(0x7fa274760000, 4096, PROT_NONE) = 0
    clone(child_stack=0x7fa274f5ffd0,
         flags=CLONE_VM|CLONE_FS|CLONE_FILES|CLONE_SIGHAND|
               CLONE_THREAD|CLONE_SYSVSEM|CLONE_SETTLS|
               CLONE_PARENT_SETTID|CLONE_CHILD_CLEARTID,  
         parent_tidptr=0x7fa274f609d0, tls=0x7fa274f60700,    
         child_tidptr=0x7fa274f609d0) = 31316
    

我们可以看到，在调用pthread_create的时候，首先使用mmap分配了8392704 Byte（8196kB）堆栈空间，但是在创建线程的时候，如果不指定堆栈大小，理应使用系统定义的默认最大空间， 通过ulimit -s 可以看到值为8192kB

    hzsunjianliang@inspur1:~/github/golang$ ulimit -s
    8192
    

mmap多映射了4k（一页），可以看到在mmap之后又调用mprotect将堆栈尾部空间的权限设置为PROT_NONE即不可读写和执行，所以基本可以判断多mmap的1页内存空间主要是用于内存溢出情况下的检测。最后将mmap返回的堆栈作为clone的一个参数创建一个线程。 通过mmap分配的内存并不会直接映射为实际使用的物理内存空间，只有当实际使用的时候，在发现当前虚拟地址空间没有分配实际无力内存的情况下，会触发操作系统缺页中断从而再分配实际物理内存。

## golang  contiguous stack 实现

goroutine作为golang的独立调度单元，每个goroutine能够独立运行的重要元素为其独立的栈空间，golang 1.2的实现类似于pthread，分配固定大小的空间，由于是固定的所以即不能太大也不能太小。而Go1.3 引入了 contiguous stack，可以在goroutine初试创建时分配非常小的栈空间（1.3为4k，后续1.5roadmap中说到会减到2k），在使用过程中自动进行增长和收缩。这使得我们可以在golang中创建很多很多的goroutine而不用担心内存耗尽。这激发我们编写各种各样的并发模型而不用太担心其可能对内存照成很大的开销。

## 实现原理

golang在每次执行函数调用的时候，首先，其runtime会检测当前的栈空间是否足够使用，如果不够使用，会触发类似“缺页中断”，Go 的runtime会保存此事函数的上下文环境，然后malloc一块内存，将旧堆栈的内存copy到新的堆栈，并做一些合理的调整。当函数返回的时候，函数会在新的堆栈中继续运行，仿佛整个过程啥事都没发生过。所以理论上来说goroutine可以使用“无限大的堆栈空间”

## 实现细节

Go的运行库中，每个goroutine对应一个结构体G（类似于linux操作系统的中进程控制块），此结构中保存有stackbase 和stackguard用于定义其使用的栈信息，每次函数调用时候都会检测当前函数需要使用的栈空间是否够用，如果不够用就进行扩张。

接下来我们分析golang的汇编代码进行分析

    package main 
    import  "fmt"  
    func main(){ 
        a := 1 
        strb := "hello " 
        a = a + 1 
        strb += "world" 
        fmt.Print(a, strb) 
        main()  
    }
    

go tool 6g -S continuousStack.go | head -8

    "".main t=1 size=352 value=0 args=0 locals=0xb8
    000000 00000 (continuousStack.go:5) TEXT         "".main+0(SB),$184-0
    000000 00000 (continuousStack.go:5) MOVQ    (TLS),CX
    0x0009 00009 (continuousStack.go:5) LEAQ    -56(SP),AX
    0x000e 00014 (continuousStack.go:5) CMPQ    AX,(CX)
    0x0011 00017 (continuousStack.go:5) JHI ,26
    0x0013 00019 (continuousStack.go:5) CALL,runtime.morestack00_noctxt(SB)
    0x0018 00024 (continuousStack.go:5) JMP ,0
    0x001a 00026 (continuousStack.go:5) SUBQ    $184,SP
    

从上面可以看到，在进入main函数之后，首先从TLS中取得第一个字段，也就是g－>stackguard字段，然后将当前SP值减去函数预计将要使用的局部堆栈空间56byte，如果得到的值小于stackguard则表示当前栈空间不够使用，需要调用runtime.morestack分配更大的堆栈空间。

more：\[连续栈\]\[1\]

## 参考资料

1.  https://github.com/tiancaiamao/go-internals/blob/master/ebook/03.5.md
2.  https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub
3.  http://stackoverflow.com/questions/6270945/linux-stack-sizes
4.  http://www.unix.com/unix-for-dummies-questions-and-answers/174134-kernel-stack-vs-user-mode-stack.html
