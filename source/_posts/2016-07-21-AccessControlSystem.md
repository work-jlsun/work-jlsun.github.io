title: Access-Control-System
date: 2016-07-21 17:39:59
tags:
---


给企业提供云计算服务，提供基础设施服务，满足基本的功能需求之外，还需要很重要一点是层级的权限控制体系，如果没有权限控制体系，对企业用户来说其实就只有一个root用户,这是远远不能满足企业对数据和资源的管理和各种安全需求的。当今国内外的云服务平台很多，但是在权限控制上做得较为完善的企业并不是很多。国内的话阿里云的 RAM服务、国外的就属AWS的IAM服务。如下为权限控制系统的一些基本原理。以下为个人就wiki上访问控制的一些粗浅的理解。

### 1. 基本理论

要了解访问控制系统，必须首先了解几个基本元素，身份、认证、授权、访问控制、审计。

 * identification
 * authentication
 * authorization
 * access approval (a narrow definition)
 * audit

Access control systems provide the essential services of authorization, identification and authentication (I&A), access approval, and accountability.

* **identification and authentication:** ensure that only legitimate subjects can log on to a system.

* **authorization:** :specifies what a subject can do.

* **access approval:** :grants access during operations, by association of users with the resources that they are allowed to access, based on the **Authorization policy(授权策略)**.

* **audit:** accountability identifies what a subject (or all subjects associated with a user) did

![xxx](/media/files/2016/07/acl.jpg)

#### 1.1 Authorization(授权)

Authorization involves the act of defining access-rights for subjects. An **Authorization policy(授权策略)** specifies the operations are allowed to execute within a system. Most modern operation systems implement authorization policies as permissions that are variations or extensions of three basic types of access.

* Read (R): The subject can
	* Read file contents
	* List directory contents
* Write (W): The subject can change the contents of a file or directory with the following tasks
	* Add
	* Update
	* Delete
	* Rename
* Execute (X): If the file is a program, the subject can cause the program to be run. (In Unix-style systems, the "execute" permission doubles as a "traverse directory" permission when granted for a directory.)

如下为AWS IAM系统或者Aliyun RAM系统中授权策略的基本描述形势,如下授权策略的意思是允许对队该账户下的所有OSS资源进行只读访问。

```
{
  "Statement": [
    {
      "Action": [  //允许执行的动作：列举、访问对象等操作
        "oss:Get*",
        "oss:List*"
      ],
      "Effect": "Allow",   //授予或者不授予
      "Resource": "*" //允许访问的资源范围
    }
  ],
  "Version": "1" //策略版本(系统使用)
}
```

当然这个策略需要attach授予给某个具体的用户才能够使得对应用户拥有该权限。

#### 1.2 Identification and Authentication (I&A)（身份及验证）

Identification and Authentication is the process of verifying that an identity is bound to the entity that makes an assertion or claim of identity. 
The I&A process assumes that there was an initial validation of the identity, commonly called identity proofing. Various methods of identity proofing are available, ranging from **in-person validation using government issued identification(真实身份)**, **to anonymous methods that allow the claimant to remain anonymous(虚拟身份)**, but known to the system if they return. The method used for identity proofing and validation should provide an assurance level commensurate with the intended use of the identity within the system(**提供预期用途相称的安全验证水平和方案**). Subsequently, the entity asserts an identity together with an authenticator as a means for validation. The only requirements for the identifier is that it must be unique within its security domain.

Authenticators are commonly based on at least one of the following four factors.

* Something you know, such as a password or a personal identification number (PIN). This assumes that only the owner of the account knows the password or PIN needed to access the account.（**密码**）

* Something you have, such as a smart card or security token. This assumes that only the owner of the account has the necessary smart card or token needed to unlock the account.(**MFA 多因素认证**)

* Something you are, such as fingerprint, voice, retina, or iris characteristics.(**人体身份**)

* Where you are, for example inside or outside a company firewall, or proximity of login location to a personal GPS device.(**其它条件限制**)

#### 1.3 Access approval(访问控制)

Access approval is the function that actually grants or rejects access during operations.

During access approval, the system compares the **formal representation of the authorization policy** with the **access request**, to determine whether the request shall be granted or rejected. Moreover, the access evaluation can be done online/ongoing.


#### 1.4 Accountability(审计)

Accountability uses such system components as audit trails (records) and logs, to associate a subject with its actions. The information recorded should be sufficient to map the subject to a controlling user. Audit trails and logs are important for

* Detecting security violations(检查安全回归)
* Re-creating security incidents(重现安全事故)

If no one is regularly reviewing your logs and they are not maintained in a secure and consistent manner, they may not be admissible as evidence.

Many systems can generate automated reports, based on certain predefined criteria or thresholds, known as clipping levels. For example, a clipping level may be set to generate a report for the following.

* More than three failed logon attempts in a given period(多次登陆失败)
* Any attempt to use a disabled user account（禁用账户）

These reports help a system administrator or security administrator to more easily identify possible break-in attempts.
 

### 2 访问控制基本模型

主流的访问控制包含一下几大类

* DAC(Discretionary Access Control Model):自由访问控制
* MAC(Mandatory Access Control Model):强制访问控制
* RBAC(Role Based Access Control):基于角色的访问控制

#### 2.1 DAC(Discretionary Access Control Model)

Discretionary access control (DAC) is a policy determined by the owner of an object. The owner decides who is allowed to access the object, and what privileges they have.(**基于资源实际拥有者的访问权限控制**)

Two important concepts in DAC are

* **File and data ownership**: Every object in the system has an owner. In most DAC systems, each object's initial owner is the subject that caused it to be created. The access policy for an object is determined by its owner.(此类系统中每个object(资源)都有自己的owner，owener一般来说都是资源的创建者，资源的创建者可以控制资源的访问控制)

* **Access rights and permissions**: These are the controls that an owner can assign to other subjects for specific resources.

Access controls may be discretionary in ACL-based or capability-based access control systems. (In capability-based systems, there is usually no explicit concept of 'owner', but the creator of an object has a similar degree of control over its access policy.)

DAC一般可以通过**ACLs(Access Control Lists)**的方式来进行实现。

![tab1](/media/files/2016/07/tb1.jpg)

ACL是最早的一种访问控制机制，它的原理非常简单，每一项资源都有一个列表，这个列表纪录的就是那些用户可以对着项资源执行CURD操作。当系统试图访问这项资源时，会首先检查这个列表中是否有关于当前用户的访问权限，从而确定是否可以执行相应的操作。ACL是一种**面向资源的访问控制模型**，它的机制都是围绕着“资源”展开的。
　　由于ACL的简单性，使得它几乎不需要任何基础设施就可以完成访问控制。但同时它的缺点也是很明显的，由于需要维护大量的访问权限列表，ACL在性能上有明显的缺陷。另外，对于拥有大量用户与众多资源的应用，管理访问控制列表本身就变成非常繁重的工作。

ps:自由访问控制之所以称之为"自由"的原因为：权限的传递较为随意和自由，比如linux下的文件权限控制，文件的拥有者可以非常方便赋予其它任意用户进行读写删的权限。但是某些情况下其实这种自由是非常可怕的，比如说机构中有各种角色，某些用户属于绝密级别、有些用户属于普通用户(扫地阿姨)，一般来说访问控制系统不能够提供任何一种途径使得绝密用户的数据透漏给普通用户；那么下面MAC既为与DAC相对应的访问控制系统。


#### 2.2 MAC(Mandatory Access Control Model)

Mandatory access control refers to allowing access to a resource if and only if rules exist that allow a given user to access the resource. It is difficult to manage, but its use is usually justified when used to protect highly sensitive information. Examples include certain government and military information. Management is often simplified (over what can be required) if the information can be protected using hierarchical access control, or by implementing sensitivity labels. What makes the method "mandatory" is the use of either rules or sensitivity labels.

(PS:资源上必须显式明确某个用户是否有权限操作这个资源，而且权限的控制一般不是资源的owner可以控制的，所以纯粹的这种访问控制一般很难管理，所以很多简化的实现是给用户一个类似的保密级别，资源一个保密级别，用户的保密级别只有高于资源的保密级别才能够访问)


Subject表：

![tab2](/media/files/2016/07/tb2.jpg)



Object表：

![tab3](/media/files/2016/07/tb3.jpg)


* Sensitivity labels

In such a system subjects and objects must have labels assigned to them. A subject's sensitivity label specifies its level of trust. An object's sensitivity label specifies the level of trust required for access. In order to access a given object, the subject must have a sensitivity level equal to or higher than the requested object.

* Data import and export

Controlling the import of information from other systems and export to other systems (including printers) is a critical function of these systems, which must ensure that sensitivity labels are properly maintained and implemented so that sensitive information is appropriately protected at all times.

Two methods are commonly used for applying mandatory access control.

#### 2.3 RBAC(Role-based Access Model)

Role-based access control (RBAC) is an access policy determined by the system, not by the owner. RBAC is used in commercial applications and also in military systems, where multi-level security requirements may also exist.

**RBAC differs from DAC** in that DAC allows users to control access to their resources. while in RBAC, access is controlled at the system level, outside of the user's control. Although RBAC is non-discretionary, **it can be distinguished from MAC** primarily in the way permissions are handled. MAC controls read and write permissions based on a user's clearance level and additional labels.RBAC controls collections of permissions that may include complex operations such as an e-commerce transaction, or may be as simple as read or write. A role in RBAC can be viewed as a set of permissions.
Three primary rules are defined for RBAC:

Role assignment: A subject can execute a transaction only if the subject has selected or been assigned a suitable role.
Role authorization: A subject's active role must be authorized for the subject. With rule 1 above, this rule ensures that users can take on only roles for which they are authorized.
Transaction authorization: A subject can execute a transaction only if the transaction is authorized for the subject's active role. With rules 1 and 2, this rule ensures that users can execute only transactions for which they are authorized.

Additional constraints may be applied as well, and roles can be combined in a hierarchy where higher-level roles subsume permissions owned by lower-level sub-roles.

Most IT vendors offer RBAC in one or more products.

总结:RBAC是把用户按角色进行归类，通过用户的角色来确定用户能否针对某项资源进行某项操作。RBAC相对于ACL最大的优势就是简化了用户与权限的管理，通过对用户进行分类，使得角色与权限关联起来，用户与权限变成了间接的关联。RBAC模型使得访问控制，特别是对用户的授权管理变得较为简单和易于维护，因此有广泛的使用。（但是相对的缺陷就是，如果要是细粒度的控制，比如某个角色下的个别用户需要特别的权限定制，如同加入一些其它角色的小部分权限或去除当前角色的一些权限是，就可能需要新建角色，然后分配新的角色给这个单独的用户，这样以来，如果这种需求越来越多，那么就会慢慢退化为跟ACL一样繁琐了）。

### 3. 参考资料

* [wiki](https://en.wikipedia.org/wiki/Computer_access_control)
* [AccessControlList](http://nos.ntease.com/doc/accesscontrolpresentation-130630194314-phpapp01.pdf)
* [访问控制](http://www.cec-ceda.org.cn/information/book/info_6.htm)
* [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)


