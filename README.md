# mydocker

## 1. 概述

参考《自己动手写 docker》从零开始实现一个简易的 docker 以及配套教程。

具体差异如下：

* UnionFS替换：从AUFS 替换为 Overlayfs
* 依赖管理更新：从 go vendor 替换为 Go Module
* 一些写法上的优化调整


### 微信公众号：探索云原生

> 鸽了很久之后，终于开通了，欢迎关注。

一个云原生打工人的探索之路，专注云原生，Go，坚持分享最佳实践、经验干货。

扫描下面的二维码关注我的微信公众帐号，一起`探索云原生`吧~

![](https://img.lixueduan.com/about/wechat/qrcode_search.png)

> `从零开始写 Docker` 系列更新中~

### 个人博客：指月小筑(探索云原生)
在线阅读：[指月小筑(探索云原生)](https://www.lixueduan.com/categories/docker/)


## 2. 基础知识

推荐阅读以下文章对 Docker 核心原理有一个大致认识：
* **核心原理**：[深入理解 Docker 核心原理：Namespace、Cgroups 和 Rootfs](https://www.lixueduan.com/posts/docker/03-container-core/)
* **基于 namespace 的视图隔离**：[探索 Linux Namespace：Docker 隔离的神奇背后](https://www.lixueduan.com/posts/docker/05-namespace/)
* **基于 cgroups 的资源限制**
    * [初探 Linux Cgroups：资源控制的奇妙世界](https://www.lixueduan.com/posts/docker/06-cgroups-1/)
    * [深入剖析 Linux Cgroups 子系统：资源精细管理](https://www.lixueduan.com/posts/docker/07-cgroups-2/)
    * [Docker 与 Linux Cgroups：资源隔离的魔法之旅](https://www.lixueduan.com/posts/docker/08-cgroups-3/)
* **基于 overlayfs 的文件系统**：[Docker 魔法解密：探索 UnionFS 与 OverlayFS](https://www.lixueduan.com/posts/docker/09-ufs-overlayfs/)
* **基于 veth pair、bridge、iptables 等等技术的 Docker 网络**：[揭秘 Docker 网络：手动实现 Docker 桥接网络](https://www.lixueduan.com/posts/docker/10-bridge-network/)

通过上述文章，大家对 Docker 的实现原理已经有了初步的认知，接下来我们就用 Golang 手动实现一下自己的 docker(mydocker)。


## 3. 具体实现

### 构造容器

本章构造了一个简单的容器，具有基本的 Namespace 隔离，确定了基本的开发架构，后续在此基础上继续完善即可。

第一篇：
* [从零开始写 Docker：实现 run 命令](https://www.lixueduan.com/posts/docker/mydocker/01-mydocker-run/)
* 代码分支 [feat-run](https://github.com/lixd/mydocker/tree/feat-run)

第二篇：
* [从零开始写 Docker(二)---优化：使用匿名管道传参](https://www.lixueduan.com/posts/docker/mydocker/02-passing-param-by-pipe/)
* 代码分支 [opt-passing-param-by-pipe](https://github.com/lixd/mydocker/tree/opt-passing-param-by-pipe)

第三篇：
* [从零开始写 Docker(三)---基于 cgroups 实现资源限制](https://www.lixueduan.com/posts/docker/mydocker/03-resource-limit-by-cgroups/)
* 代码分支 [feat-cgroup](https://github.com/lixd/mydocker/tree/feat-cgroup)





### 构造镜像

本章首先使用 busybox 作为基础镜像创建了一个容器，理解了什么是 rootfs，以及如何使用 rootfs 来打造容器的基本运行环境。

然后，使用 OverlayFS 来构建了一个拥有二层模式的镜像，对于最上层可写层的修改不会影响到基础层。这里就基本解释了镜像分层存储的原理。

之后使用 -v 参数做了一个 volume 挂载的例子，介绍了如何将容器外部的文件系统挂载到容器中，并且让它可以访问。

最后实现了一个简单版本的容器镜像打包。

这一章主要针对镜像的存储及文件系统做了基本的原理性介绍，通过这几个例子，可以很好地理解镜像是如何构建的，第 5 章会基于这些基础做更多的扩展。

第四篇：

* [从零开始写 Docker(四)---使用 pivotRoot 切换 rootfs 实现文件系统隔离](https://www.lixueduan.com/posts/docker/mydocker/04-change-rootfs-by-pivot-root/)
* 代码分支 [feat-rootfs](https://github.com/lixd/mydocker/tree/feat-rootfs)

第五篇：

* [从零开始写 Docker(五)---基于 overlayfs 实现写操作隔离](https://www.lixueduan.com/posts/docker/mydocker/05-isolate-operate-by-overlayfs/)
* 代码分支 [feat-overlayfs](https://github.com/lixd/mydocker/tree/feat-overlayfs)

第六篇：

* [从零开始写 Docker(六)---实现 mydocker run -v 支持数据卷挂载](https://www.lixueduan.com/posts/docker/mydocker/06-volume-by-bind-mount/)
* 代码分支 [feat-volume](https://github.com/lixd/mydocker/tree/feat-volume)

第七篇：

* [从零开始写 Docker(七)---实现 mydocker commit 打包容器成镜像](https://www.lixueduan.com/posts/docker/mydocker/07-mydocker-commit/)
* 代码分支 [feat-commit](https://github.com/lixd/mydocker/tree/feat-commit)


### 构建容器进阶

本章实现了容器操作的基本功能。

* 首先实现了容器的后台运行，然后将容器的状态在文件系统上做了存储。
* 通过这些存储信息，又可以实现列出当前容器信息的功能。
* 并且， 基于后台运行的容器，我们可以去手动停止容器，并清除掉容器的存储信息。
* 最后修改了上一章镜像的存储结构，使得多个容器可以并存，且存储的内容互不干扰。

第八篇：

* [从零开始写 Docker(八)---实现 mydocker run -d 支持后台运行容器](https://www.lixueduan.com/posts/docker/mydocker/08-mydocker-run-d/)
* 代码分支 [feat-run-d](https://github.com/lixd/mydocker/tree/feat-run-d)

第九篇：

* [从零开始写 Docker(九)---实现 mydocker ps 查看运行中的容器](https://www.lixueduan.com/posts/docker/mydocker/09-mydocker-ps/)
* 代码分支 [feat-ps](https://github.com/lixd/mydocker/tree/feat-ps)


第十篇：

* [从零开始写 Docker(十)---实现 mydocker logs 查看容器日志](https://www.lixueduan.com/posts/docker/mydocker/10-mydocker-logs/)
* 代码分支 [feat-logs](https://github.com/lixd/mydocker/tree/feat-logs)


---
[05-4-实现进入容器：mydocker exec](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-4-%E5%AE%9E%E7%8E%B0%E8%BF%9B%E5%85%A5%E5%AE%B9%E5%99%A8%20Namespace.md)

[05-5-实现停止容器：mydocker stop](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-5-%E5%AE%9E%E7%8E%B0%E5%81%9C%E6%AD%A2%E5%AE%B9%E5%99%A8.md)

[05-6-实现删除容器：mydocker rm]()

[05-7-文件系统重构](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-7-%E5%AE%9E%E7%8E%B0%E9%80%9A%E8%BF%87%E5%AE%B9%E5%99%A8%E5%88%B6%E4%BD%9C%E9%95%9C%E5%83%8F.md)

> refactor: 文件系统重构,为不同容器提供独立的rootfs. feat: 更新rm命令，删除容器时移除对应文件系统. feat: 更新commit命令，实现对不同容器打包.

[05-8 实现环境变量注入功能：mydocker run -e ](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-8-%E5%AE%9E%E7%8E%B0%E6%8C%87%E5%AE%9A%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F%E8%BF%90%E8%A1%8C.md)



### 容器网络

在这一章中，首先手动给一个容器配置了网路，并通过这个过程了解了 Linux 虚拟网络设备和操作。然后构建了容器网络的概念模型和模块调用关系、IP 地址分配方案，以及网络模块的接口设计和实现，并且通过实现 Bridge
驱动给容器连上了“网线”。

[06-1-网络虚拟化技术](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/06-1-%E7%BD%91%E7%BB%9C%E8%99%9A%E6%8B%9F%E5%8C%96%E6%8A%80%E6%9C%AF.md)

[06-2-构建容器网络模型](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/06-2-%E6%9E%84%E5%BB%BA%E5%AE%B9%E5%99%A8%E7%BD%91%E7%BB%9C%E6%A8%A1%E5%9E%8B.md)

[06-3-容器地址分配](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/06-3-%E5%AE%B9%E5%99%A8%E5%9C%B0%E5%9D%80%E5%88%86%E9%85%8D.md)

[06-4-创建Bridge网络](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/06-4-%E5%88%9B%E5%BB%BABridge%E7%BD%91%E7%BB%9C.md)

[06-5-在Bridge网络创建容器](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/06-5-%E5%9C%A8Bridge%E7%BD%91%E7%BB%9C%E5%88%9B%E5%BB%BA%E5%AE%B9%E5%99%A8.md)
