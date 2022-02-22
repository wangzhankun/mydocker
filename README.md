# mydocker

## 1. 概述


跟着《自己动手写 docker》从零开始实现一个简易的 docker。

建议先了解一下 Docker 的核心原理大致分析，可以看这篇文章：[Docker教程(三)---核心实现原理分析](https://www.lixueduan.com/post/docker/03-container-core/)。

核心原理一共包含以下三个点：

* 1）Namespace  
* 2）Cgroups  
* 3）UnionFS



## 2. 基础知识

### namespace

[namespace 初体验](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/1-namespace-2%E5%88%9D%E4%BD%93%E9%AA%8C.md)



### Cgroups

[cgroup-1-初体验](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgoups-1%E5%88%9D%E4%BD%93%E9%AA%8C.md)

[cgroup-2-相关操作](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgroups-2%E7%9B%B8%E5%85%B3%E6%93%8D%E4%BD%9C.md)

[cgroup-3-演示](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgroups-3%E6%BC%94%E7%A4%BA.md)

[cgroup-4-go操作cgroup](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgroups-4go%E8%AF%AD%E8%A8%80%E6%93%8D%E4%BD%9C.md)



### ufs

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

之后使用 -v参 数做了一个 volume 挂载的例子，介绍了如何将容器外部的文件系统挂载到容器中，并且让它可以访问。

最后实现了一个简单版本的容器镜像打包。

这一章主要针对镜像的存储及文件系统做了基本的原理性介绍，通过这几个例子，可以很好地理解镜像是如何构建的，第 5 章会基于这些基础做更多的扩展。

[04-1-使用busybox做rootfs](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-1-rootfs.md)

[04-2-overlayfs](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-2-overlayfs.md)

[04-3-实现数据卷挂载 mydocker -v](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-3-volume.md)

[04-4-实现简单镜像打包 mydocker commit](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/04-4-%E5%AE%9E%E7%8E%B0%E7%AE%80%E5%8D%95%E9%95%9C%E5%83%8F%E6%89%93%E5%8C%85.md)



### 构建容器进阶

[05-1-实现容器后台运行 mydocker run -d](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-1-%E5%AE%9E%E7%8E%B0%E5%AE%B9%E5%99%A8%E5%90%8E%E5%8F%B0%E8%BF%90%E8%A1%8C.md)

[05-2-实现查看运行中的容器 mydocker ps](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-2-%E5%AE%9E%E7%8E%B0%E6%9F%A5%E7%9C%8B%E8%BF%90%E8%A1%8C%E4%B8%AD%E7%9A%84%E5%AE%B9%E5%99%A8.md)

[05-3-实现查看容器日志 mydocker log](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/05-3-%E5%AE%9E%E7%8E%B0%E6%9F%A5%E7%9C%8B%E5%AE%B9%E5%99%A8%E6%97%A5%E5%BF%97.md)

[05-4-实现进入容器 Namespace docker exec]()
