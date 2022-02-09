# mydocker

## 1. 概述

《自己动手写 docker》笔记和源码


建议先了解一下 Docker 的核心原理大致分析，可以看这篇文章：[Docker教程(三)---核心实现原理分析](https://www.lixueduan.com/post/docker/03-container-core/)。

核心原理一共包含以下三个点：

* 1）Namespace  
* 2）Cgroups  
* 3）UnionFS



## 2. 基础知识

### namespace

[namespace-1概述](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/1-namespace-1%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86.md)

[namespace-2初体验](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/1-namespace-2%E5%88%9D%E4%BD%93%E9%AA%8C.md)



### Cgroups

[cgroup-1-初体验](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgoups-1%E5%88%9D%E4%BD%93%E9%AA%8C.md)

[cgroup-2-相关操作](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgroups-2%E7%9B%B8%E5%85%B3%E6%93%8D%E4%BD%9C.md)

[cgroup-3-演示](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgroups-3%E6%BC%94%E7%A4%BA.md)

[cgroup-4-go操作cgroup](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/2-cgroups-4go%E8%AF%AD%E8%A8%80%E6%93%8D%E4%BD%9C.md)



### ufs

[ufs-1-2初体验](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/3-ufs-1%E5%88%9D%E4%BD%93%E9%AA%8C.md)

[ufs-2-overlayfs2](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/3-ufs-2overlay.md)

[ufs-3-docker文件系统](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/3-ufs-3docker%E6%96%87%E4%BB%B6%E7%B3%BB%E7%BB%9F.md)



## 3. 具体实现

### 构造容器
本章构造了一个简单的容器，具有基本的Namespace隔离，确定了基本的开发架构，后续在此基础上继续完善即可。


[3-1 实现 mydocker run 命令](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/03-1-%E5%AE%9E%E7%8E%B0run%E5%91%BD%E4%BB%A4.md)

[3-2 增加Cgroups实现资源限制](https://github.com/lixd/daily-notes/blob/master/Golang/mydocker/03-2-%E5%A2%9E%E5%8A%A0cgroups.md)

