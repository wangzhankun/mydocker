# mydocker

## 1. 概述

参考《自己动手写 docker》从零开始实现一个简易的 docker 以及配套教程。

具体差异如下：

* UnionFS替换：从AUFS 替换为 Overlayfs
* 依赖管理更新：从 go vendor 替换为 Go Module
* 一些写法上的优化调整



建议先了解一下 Docker 的核心原理大致分析，可以看这篇文章：[Docker教程(三)---核心实现原理分析](https://www.lixueduan.com/post/docker/03-container-core/)。

核心原理一共包含以下三个点：

* 1）Namespace
* 2）Cgroups
* 3）UnionFS

## 2. 基础知识

### Namespace

[namespace 初体验](https://www.lixueduan.com/post/docker/05-namespace/)



### Cgroups

[Cgroups-1-初体验](https://www.lixueduan.com/post/docker/06-cgroups-1/)

[Cgroups-2-subsystem演示](https://www.lixueduan.com/post/docker/07-cgroups-2/)

[Cgroups-3-相关命令汇总及Go Demo](https://www.lixueduan.com/post/docker/08-cgroups-3/)



### UnionFS

[ufs-1-初体验](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/3-ufs-1%E5%88%9D%E4%BD%93%E9%AA%8C.md)

[ufs-2-overlayfs2](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/3-ufs-2overlay.md)

[ufs-3-docker文件系统](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/3-ufs-3docker%E6%96%87%E4%BB%B6%E7%B3%BB%E7%BB%9F.md)



## 3. 具体实现

### 构造容器

本章构造了一个简单的容器，具有基本的 Namespace 隔离，确定了基本的开发架构，后续在此基础上继续完善即可。

[3-1 实现 mydocker run 命令](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/03-1-%E5%AE%9E%E7%8E%B0run%E5%91%BD%E4%BB%A4.md)

[3-2 增加Cgroups实现资源限制](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/03-2-%E5%A2%9E%E5%8A%A0cgroups.md)



### 构造镜像

本章首先使用 busybox 作为基础镜像创建了一个容器，理解了什么是 rootfs，以及如何使用 rootfs 来打造容器的基本运行环境。

然后，使用 OverlayFS 来构建了一个拥有二层模式的镜像，对于最上层可写层的修改不会影响到基础层。这里就基本解释了镜像分层存储的原理。

之后使用 -v 参数做了一个 volume 挂载的例子，介绍了如何将容器外部的文件系统挂载到容器中，并且让它可以访问。

最后实现了一个简单版本的容器镜像打包。

这一章主要针对镜像的存储及文件系统做了基本的原理性介绍，通过这几个例子，可以很好地理解镜像是如何构建的，第 5 章会基于这些基础做更多的扩展。

[04-1-使用busybox做rootfs](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-1-rootfs.md)

[04-2-overlayfs](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-2-overlayfs.md)

[04-3-实现数据卷挂载 mydocker -v](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-3-volume.md)

[04-4-实现简单镜像打包 mydocker commit](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-4-%E5%AE%9E%E7%8E%B0%E7%AE%80%E5%8D%95%E9%95%9C%E5%83%8F%E6%89%93%E5%8C%85.md)



### 构建容器进阶

本章实现了容器操作的基本功能。

* 首先实现了容器的后台运行，然后将容器的状态在文件系统上做了存储。
* 通过这些存储信息，又可以实现列出当前容器信息的功能。
* 并且， 基于后台运行的容器，我们可以去手动停止容器，并清除掉容器的存储信息。
* 最后修改了上一章镜像的存储结构，使得多个容器可以并存，且存储的内容互不干扰。

[05-1-实现容器后台运行：mydocker run -d](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-1-%E5%AE%9E%E7%8E%B0%E5%AE%B9%E5%99%A8%E5%90%8E%E5%8F%B0%E8%BF%90%E8%A1%8C.md)

[05-2-实现查看运行中的容器：mydocker ps](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-2-%E5%AE%9E%E7%8E%B0%E6%9F%A5%E7%9C%8B%E8%BF%90%E8%A1%8C%E4%B8%AD%E7%9A%84%E5%AE%B9%E5%99%A8.md)

[05-3-实现查看容器日志：mydocker log](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-3-%E5%AE%9E%E7%8E%B0%E6%9F%A5%E7%9C%8B%E5%AE%B9%E5%99%A8%E6%97%A5%E5%BF%97.md)

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
