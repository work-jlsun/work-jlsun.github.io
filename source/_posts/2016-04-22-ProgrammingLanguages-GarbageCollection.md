title: Garbage Collection
date: 2016-04-22 17:12:41
tags:
- GC
- Pragraming Languages
---


前一整子研究了下golang的GC机制，然后顺藤摸瓜找了些基本的GC算法的资料，学习了下GC算法方面的内容。

推荐 The University of Texas at Austin 计算机系[编程语言](http://www.cs.utexas.edu/~shmat/courses/cs345/cs345_home.html)一课中的 [《Garbage Collection》](http://www.cs.utexas.edu/~shmat/courses/cs345/11garbage.ppt) 一课的内容，叫兽就是教授，总结的真好。看看能不能翻译、整理、吸收下。

#### 1 内存的分布

主要的几种内存的分布区域

* Static area

固定的大小，固定的内容，在编译的时候已经确定分配。

* Runtime stack

可变的大小，可变的内容(动态纪录)；典型的场景为函数的调用链(局部变量)和返回。

* Heap 

固定的大小，变化的内容；动态分配的objects和data structures; 比如malloc in C，new in Java。


#### 2 Cells And LiveNess

Cell: 在heap中的数据单元；

Cells 由存在于各种存储区域中的指针所指向

* registers、stack中的指针
* global/static memory
* 其它的heap cells

Roots(根对象): registers、stack locations、global/static variables

一个cell是存活的当其仅当其指针存在于Root中或者由其它存活的cell间接指向。


#### 3 Garbage

**Garbage** - 程序运行过程中产生的无法被访问到的heap memory 或者内存释放了但是在程序里头还有引用的heap memory。

* An allocated block of heap memory does not have a reference to it (cell is no longer “live”)
  
* Another kind of memory error: a reference exists to a block of memory that is no longer allocated
  

**Garbage Collection** -  自动管理动态申请的内存空间，即自动回收不再使用的head memory以便后续继续提供给程序使用。

```
class Node {
	int value;
	node next;
}
node p,q;

p = new node();
q = new node();

q = p;
delete p;
```

![Example Of Garbage](/media/files/2016/04/garbage1.jpg)


##### 3.1 Why Garbage Collection

为什么现代程序设计语言基本都选择具备垃圾回收机制？

- 如今的应用程序非常可以随意自由的进行内存分配
	- 1GB laptop,1-4GB desktops, 8~512GB servers
	- 64-bit address spaces 

- 但是却没有很好的管理
	- Memorey leaks, dangling references, double free, misaligned addresses, null pointer deference, heap fragmentation
	- Poor use of refererce locality, resulting in highcache miss rates and/or excessive demand paging

- 显式的内存管理容易打破上层编程逻辑的抽象


##### 3.2 The Perfect Garbage Collector

完美的垃圾回收应该具备如下特点

基本要求

* No visible impact on program execution
* Works with any program and its data structures
	* For Eaxmple，handles cyclic data structures

进阶要求

* Collects garbage(and obly garbage) cells quickly
	* Incremental; can meet real-time constraints
* Has excellent spatial locality of reference
	* No excessive paging, no negative cache effects
* Manages the heap efficiently
	* Always satisfies an allocation request and does not fragment
	 

##### 3.3 Summary of GC Techiques

基本的GC技术及其核心特点如下

* Reference counting
	*  直接跟踪所有的live cell
	*  不论是否在heap上分配内存，GC机制都会发生
	*  不会发现所有的垃圾	 	
* Tracing
	* GC takes place and identyfies live cells when a request for memory fails 	
	* Mark-sweep
	* Copy  collection
* Modern techniques: generational GC 
	
以下详细介绍这几种常见GC算法和其优缺点。


#### 4 Reference Counting

* 简单的实时统计每一个cell的引用计数
* 存储count 和 引用计数增减会造成 一定的space和time的开销
	* Reference counts是实时维护的，所以没有 "stop-and-collect" 即停止程序进行垃圾回收的过程
	* 一种天生的增量垃圾回收方式
* 典型应用场景：C++ “Smart pointer”

![ReferenceCount](/media/files/2016/04/referenceCount.jpg)

** 优点 **

* Incremental overhead(增量的开销)
	* cell管理插入在程序的正常程序执行逻辑过程中
	* 适合于对交互性或者实时计算型负载
* 实现比较简单
* 与手动内存管理可以共存
* Spatial locality of reference is good
	* Access pattern to virtual memory pages no worse than the program, so no excessive paging
* 可以快速进行内存复用
	* If RC == 0, put back onto the free list 	
	
** 缺点 **

* 空间开销
	*  一个word用于计数、1 个字节的用于间接指针记录
* 时间开销
	* 更新一个指针指向另外一个cell需要如下步骤
		1. 确保指针并不是自引用
		2. 在old cell上执行 count--;如果引用为0，则删除old cell
		3. 更新指针指向new cell的地址
		4. 对new cell 执行  count++

* 一次 miss的count(increment/decrement)操作会导致dangling pointer/memory leak

* 循环的数据结构可能会导致内存泄漏

![cycleReference.jpg](/media/files/2016/04/cycleReference.jpg)


** Reference Count 典型实现 “Smart Poniter” in C++ **


![std::auto_ptr in ANSI C++](/media/files/2016/04/c++autopointer.jpg)

优点分析如下:

* sizeof(RefObj<T>) =  8 bytes (8 bytes per reference-counted object)
* sizeof(Ref<T>) = 4 bytes
	* Fits in a register
	* Easily passed by value as an argument or result of a function 
	* Takes no more space than regular pointer but much "safer"

源码实现

```
tmplate<class T> class RefObj {
	private:
		T* obj;
		int cnt;
	
	public:
		RefObj(T* t): obj(t), cnt(0){}
		~RefObj() { delete obj;}
		
		int inc() { return ++cnt;}
		int dec() { return --cnt;}
		
		operator T*() { return obj;}
		operator T&() { return *obj;}
		T& operator *() { return *obj;}
}

tmplate <class T> class Ref {
	private:
		RefObj<T> *ref;
		Ref<T>* operator&{} {}
	public:
		Ref():ref(0){}
		Ref(T* p): ref(new RefObj<T>(p)){ref->inc();}
		Ref(const Ref<T>& r): ref(r.ref){ref->inc();}
		~Ref(){if (ref->dec() == 0  delete ref;)}
		
		Ref<T>& operator=(const Ref<T>& that) {
			if (this != &that) {
				if (ref->dec() == 0) delete ref;
				ref = tha.ref;
				ref->inc();
			}
			return *this;
		}
		
		T* operator->() {return *ref;}
		T& operator*() {return *ref;}
};
```

智能指针的使用

```
// auto_ptr::get example
#include <iostream>
#include <memory>
#include <string>

using namespace std;

std::auto_ptr<string> proc() {
    std::auto_ptr<string> s (new string("Hello, world")); // ref count set to 1
    int x = s->length();  // s.operator->() returns string object ptr
    cout << "s length = " << x << endl;
    return s;
} // ref count goes to 2 on copy out, then 1 when s is auto-destructed

int main()
{
    std::auto_ptr<string> a = proc();  // ref count is 1 again
} // ref count goes to zero and string is destructed, along with Ref and RefObj objects

```
[C++ auto_ptr的接口及使用方法](http://www.cplusplus.com/reference/memory/auto_ptr/)


#### 5 Mark-Sweep Garbage Collection

* 每个cell有一个 mark bit
* Garbage 在heap 使用完之前都是unreachable、undetected的、GC在heap内存用尽之后触发执行，并且程序执行开始挂起(stop-the-world)

两个阶段

* Marking phase("标记")：从所有root开始，标记所有live cell。
* Sweep phase("清理")：将所有 "未标记"的cell 汇总到freelist循环使用；clear所有mark bit的标记位，等待下次gc重新标记。

![mark-swift](/media/files/2016/04/mark-swift)

** mark-sweep 优缺点分析 **

优点

* 能够正确的处理 cycle
* 几乎没有空间开销
	* 1 bit用于标记cell(may limit max values that can be stored in a cell, e.g.., for integer cells)

缺点

* 常规的执行必须被挂起(stop the world)
* may touch all virtual memory pages
	* May lead to excessive paging if the working-set size is small and the heap is not all in physical memory.
* heap 可能会碎片化
	* Cache misses, page thrashing; more complex allocation. 


#### 6 Copying Collector

* 将heap划分为 "from-space"以及"to-space"
* from-space 中的cells是被trace的，在回收的时候将live cell ("scavenged") 拷贝到 to-space
	* 为了保证数据结构之间在回收后是正确指向的，必须更新所有原来在from-space中的cell对应的指针。 	 
	* PS：This is why references in Java and other languages are not pointers, but indirect abstractions for pointers

* 当to-space 使用完，两边space角色切换 	

![cheneys](/media/files/2016/04/cheneys.jpg)


cheneys 算法pesucode 如下 
```
class Object [
	// remains null for normal objects
	// non-null for forwarded objects
	
	Object* _forwardee;
	
	
	public:
		void forward_to(address new_addr);
		Object* forwardee();
		bool is_forwarded();
		size_t size();
		Iterator<Object**> object_fields();
};

void Object::forward_to(address new_addr){
	_forwardee = new_adrr;
}

Object* forwardee(){
	return _forwardee;
}

bool Object::is_forwarded() {
	return _forwardee != nullptr;
}


class Semispace {
	address _bottom;
	address _top;
	address _end;
	
	public:
		Semispace(address bottom, address end);
	
		address bottom() 	{	return _bottom;	}
		address top() 		{	return _top;	}
		address end()			{	return _end;	}
	
		bool contains(address obj);
		address allocate(size_t size);
		void reset();
};


Semispace::Semispace(address bottom, address end){
	_bottom 	= bottom;
	_top	 	= bottom;
	_end		= end;	
}

bool Semispace::contains(address obj){
	return _bottom <= obj && obj < _top; 
}


address Semispace::allocate(size_t size) {
	if (_top + size <= end) {
		address obj = _top;
		_top += size;
		return obj;
	}else {
		return nullptr;
	}
}

void Semispace::reset() {
	_top = _bottom;
}


class Heap {
	Semispace* _from_space;
	Semispace* _to_space;
	
	void swap_spaces();
	Obkect * evacuate(Object* obj);
	
	publice:
		Heap(address bottom, address end);
		
		address allocate(size_t size);
		void collect();
		void process_reference(Object** slot);
}

Heap::Heap(address bottom, address end) {
	size_t space_size = (end - bottom) /2;
	address boudary = bottom + space_size;
	_from_space = new Semispace(bottom, boundary);
	_to_space =  new Semispace(boundary, end);
}

void Heap::swap_spaces() {
	Smeispace *temp = _from_space;
	_from_space = to_space;
	_to_space = temp;
	
	//After sawpping, the to-space is assumed to be empty.
	//Reset its allocation pointer.
	_to_space->reset();
}

address Head::allocate(size_t size) {
	return _from_space->allocate();
}


Object* Head::evacuate(Object* obj){
	size_t size = obj->size();
	
	// allocate space in to_space and copy object to there
	address new_addr = _to_space->allocate(size);
	copy(new_addr, obj, size);
	
	// set forwarding pointer in old object;
	Object* new_obje = (Object*) new_addr;
	obj->forward_to(new_obj);
	
	return new_obj;
}


void Heap::collect() {
	// the from-space contains objects, and the to-space is empty now.
	
	address scanned = _to_space->bottom();
	
	// scavenge objects directly referenced by the root set
	foreach(Object** slot in ROOTS){
		precess_refernce(slot);
	}
	
	// breadth-first scanning of object graph(层级遍历所有对象)
	while (scaneed < _to_space->top()) {
		Object *parent_obj = (Object*)scanned;
		
		foreach(Object** slot in parent_obj->object_fields()){
			process_reference(slot);
			//note: _to_space->top() moves if any object 
			//is newly copied into to-space.
		}
		scanned += parent_obj->size();
	}
	
	//Now all live objects will have been evacated into the to_space
	// and we don't need the data in the from-space anymore
	swap_spaces();
}



void Head::process_reference(Object ** slot){
	Object *obj  = *slot;
	if (obj != nullptr && _from_space->contains(obj)){
		Object * new_obj = obj->is_forwarded()?
			obj->forwardee():evacuate(obj);
	}
	
	*slot = new_obj;
}

//参见：https://gist.github.com/rednaxelafx/8412637/forks
```

优点

* 相对较小的cell 分配开销
	* 由于分配都是连续的，out-of-space  检测只需要一次addr比较
	* 能够高效的allocate 可变大小的cell 
* 没有内存碎片问题

缺点

* 两倍的内存开销
* Copy大对象(大内存拷贝)开销大

[copying_gc.pdf](http://nos.netease.com/doc/copying_gc.pdf)


#### 7 Generational Garbage Collection

* 观测到：大多数cell的生命周期都是很短的
*  将heap区分为多个区域，对生命周期短的cell进行更加频繁的GC
	* 在 GC期间不需要扫描所有的cell
	* 周期性得对“older generations” 进行GC

如下为“分代GC”的简单示例，将young GC阶段好久没有被垃圾回收的对象迁移到“older generation”。

![generationGc](/media/files/2016/04/generationGc.jpg)

以下分代GC机制结合Copying Collector的示例。

![generationGc2](/media/files/2016/04/generationGc2.jpg)


#### 8 附录|参考

* [copying_gc.pdf](http://nos.netease.com/doc/copying_gc.pdf)
* [《Garbage Collection》](http://www.cs.utexas.edu/~shmat/courses/cs345/11garbage.ppt)

