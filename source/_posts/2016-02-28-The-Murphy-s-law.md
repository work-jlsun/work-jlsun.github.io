title: "墨菲定律(The-Murphy's-law)"
date: 2016-02-28 07:44:16
categories: 
tags: 
	- Murphy's law 
	- 分布式
---


### 墨菲定律（The-Murphy's-law）

“Anything that can go wrong, will go wrong.”－任何可能发生的事情，都将发生。

当然墨菲定律并没有说发生的事情是好是坏，因为好坏都是人为界定的；但是人们往往都用它来说明一些只有在极端苛刻条件下才能发生的糟糕事情。

对于我们IT工程师，往往需要对这个定律更加怀有一份敬畏之心。因为程序逻辑被调用的越多，那么理论上任何潜在的程序分支都会被跑到，任何潜在的bug都会发生。特别是在如今的云计算大环境下，越来越多的企业将自己的服务迁移到第三方厂商的云服务平台，所以云服务厂商程序的调用频繁度，和生命周期都是远远高于以往，墨菲定律的发生周期将被大大缩短。所以对于在大型云服务厂商coding的程序员们需要将格外敬畏。


### 潜伏Bug
此次谈一谈在我们系统中潜伏将近2年多的BUG。涉及的是对象存储服务提供的分块上传功能。

在对象存储系统中，为了提高大文件的上传成功率和对特大文件的支持，提供了“分块上传”功能，即将文件切成更小的分块进行上传，其功能的实现主要通过如下三个接口。

* InitMultiUpload: 初始化一次分块上传，返回唯一UploadID用于串联后续每个分块的上传。
* UploadPart: 将文件均匀拆分成几个块，上传每一分块数据。
* CompleteMultiUpload：所有分块上传完毕，提交所有分块的元信息合并成一个文件。

接口可参见[AWS OSS 分块接口](http://docs.aws.amazon.com/AmazonS3/latest/API/mpUploadInitiate.html)。


![multipartupload.jpg](/media/files/2016/03/multipartupload.jpg)


#### 数据库表设计

分块上传主要设计3张数据库表。

* NOS_UploadInfo表

```
CREATE TABLE `NOS_UploadInfo` (
  `UploadID` bigint(20) unsigned NOT NULL COMMENT '上传SessionID',
  `ObjectName` varchar(1000) NOT NULL COMMENT '对象名称，桶内唯一',
   .....略......
)
```
InitiateMultiUpload操作会向NOS_UploadInfo表里插入一条记录，同时给用户返回一个UploadId，以后该对象的所有分块上传都要带着这个UploadID。

* NOS_UploadPart表

```
CREATE TABLE `NOS_UploadPart` (
  `UploadID` bigint(20) unsigned NOT NULL COMMENT '所属上传SessionID',
  `DocID` bigint(20) NOT NULL COMMENT '对应的DFSDocID',
  `Sequence` smallint(5) unsigned NOT NULL COMMENT '分块序号',
  `Size` bigint(20) unsigned NOT NULL COMMENT '本分块长度',
    .....略......
)
```
每次调用UploardPart都会向NOS_UploadPart表中插入一条记录对应一个分块信息，DocId用于索引存储引擎文件的ID，Sequence为分块序号，Size为本分块的长度。

* NOS_AbandonUploadPart表

```
CREATE TABLE `NOS_AbandonUploadPart` (
  `DocID` bigint(20) unsigned NOT NULL COMMENT '废弃docid',
  `CTime` bigint(20) unsigned NOT NULL COMMENT '废弃时间',
    .....略......
) 
```
NOS_Uploadpart的主键为 (`UploadID`,`Sequence`)，当重复上传某个分块的时候，被覆盖的信息会插入NOS_AbandonUploadPart表中。

**注意：此数据库表中记录的DocID对应的存储引擎文件最终是没有实际对象会引用的,所以清理程序会周期性对NOS_AbandonUploadPart中的对象进行清理。**


#### 数据库逻辑

在所有分块上传成功之后，调用CompleteMultiUpload完成分块上传，此操作涉及的数据库逻辑如下所示。

```
(逻辑A) select DocID, Sequence, Size, ETag, LastModify from NOS_UploadPart where UploadID = ? and Sequence > ? order by Sequence asc limit ?" 
(逻辑B) 事务 {
	delete from NOS_UploadPart where UploadID = ?"
	delete from NOS_UploadInfo where UploadID = ? and IsAbort = 0
	使用上述DocID List 执行一次PUT Object SQL逻辑
}
```

其实仔细分析这个数据库逻辑，我们会发现是有问题的，因为在逻辑A和逻辑B之间有可能会有一次UploadPart操作覆盖前一次的分块，即逻辑A中获取到的DocId可能会被逻辑B之前插入的UploadPart操作丢入NOS_AbandonUploadPart列表。

对于此接口，接口的使用者一般都会选择不同分块之间进行并发上传(UploadPart)以提供上传速度，而不会对同一分块进行并发上传(没有必要)。**而实际上，发生问题的原因也并非用户并发调用同一分块上传所致, 而是由于我们系统的负载日益增大，内网万M上行网卡水位满，导致用户UploadPart调用超时从而进行重试，而某个组件处理此请求的时间由于网卡满导致执行的时间较长，在用户返回超时的时候还在执行。延时执行的请求和后续用户的重试形成了一次并发UploadPart操作。。。。**


#### 系统日志

错误发生场景：外部调用者设置的超时时间为15s，内部模块执行耗时15s以上。如下日志所示，第二条UPLOAD_PART日志即为由于内部负载过高而延迟执行的请求，其overhead为 16621，即16s以上。而第一条日志为后续用户重试的日志，然后悲剧就发生了。

```
[INFO ]2016-02-22 22:47:25,313,[Class]NosResponse, {nossvr:UPLOAD_PART, overhead:147, requestID:ee1e572b0aa000000153097300ee840f, productID:NosUp, resource:/imglf2/img/cUwreko1OTc1UG01RWR5K2kra0NWN1VXS1k2bjN1OHdqckxwMTloSEY4RitSUGRtYk4rVE5BPT0.jpg, logID:7c69b3f0-3b77-4ecc-96cc-5a75e424d102, CmdSpecInfo:UploadPartCommand [partNumber=3, uploadId=4561524247931778940, uploadPart=UploadPart [uploadID=4561524247931778940, docID=88970281898321691, sequence=3, size=1048576, lastModify=1456152445298, eTag=1e81ff2e17fdcf78fe4c9ef9d76d04f3, isCommit=false]], Date:null}[WARN ]2016-02-22 22:47:26,788,[Class]NosResponse, {nossvr:UPLOAD_PART, overhead:16621, requestID:1c0793450aa0000001530972c657849a, productID:NosUp, resource:/imglf2/img/cUwreko1OTc1UG01RWR5K2kra0NWN1VXS1k2bjN1OHdqckxwMTloSEY4RitSUGRtYk4rVE5BPT0.jpg, logID:a8cb0c09-44b7-4588-91b8-9cdb7b7e1b0b, CmdSpecInfo:UploadPartCommand [partNumber=3, uploadId=4561524247931778940, uploadPart=UploadPart [uploadID=4561524247931778940, docID=89389195827348221, sequence=3, size=1048576, lastModify=1456152446703, eTag=1e81ff2e17fdcf78fe4c9ef9d76d04f3, isCommit=false]], Date:null}```
