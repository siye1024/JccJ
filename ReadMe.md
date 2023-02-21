# Dousheng
## 一、项目介绍
- 项目简介：主要基于Kitex RPC+ Gin HTTP + MySQL实现的第五届字节跳动青训营极简抖音后端微服务项目 ——抖声
- Github地址：https://github.com/siye1024/JccJ/

## 二、项目实现
### 1. 技术选型分析
- **微服务架构vs单体架构**：

  对于本轻量级项目而言，单体架构或许具有更高的性能和代码易读性。但考虑到项目未来的功能扩展，可能的分布式部署，以及出于学习分布式架构的目的，我们最终选择微服务架构作为整体架构。
  
- **RPC框架**：

  本着简洁的原则，最初考虑使用最基础的gRPC。在学习Kitex RPC框架后，了解到Kitex提供丰富的服务治理功能（尤其是多种服务发现模式的支持），整合了两种代码生成工具，十分方便开发，因此适合作为本项目的主体框架。参考：[Kitex文档](https://www.cloudwego.io/zh/docs/kitex/)
  
- **Web框架**： 

  考虑到本项目中HTTP功能主要是暴露给用户接口，以及对用户请求进行路由转发，功能较为简单，本着简洁的原则，使用Gin作为HTTP框架 。 参考：[Gin文档](https://gin-gonic.com/zh-cn/docs/)

- **ORM框架**：

  Gorm是Golang中比较成熟的ORM框架，方便与数据库进行交互。 参考：[gorm文档](https://gorm.io/zh_CN/docs/index.html)
  
- **底层存储**：

  我们使用关系型数据库MySQL存储用户与服务相关信息，服务数据（视频，封面等）存储在高性能对象存储Minio中。参考：[Mini文档](https://min.io/)
  
  - 为什么用关系型数据库？
  
    我们可以从用户、视频、各种服务列表中分析出显而易见的关联关系
    
  - 视频数据的存储方案
  
    使用nginx反向代理 + Multi Minio可以弹性扩展存储设备，配置各种负载均衡策略，应对不断增长的业务需求
    
- 服务注册与发现：

  综合比较磁盘、网络、CPU、内存的性能开销，etcd的表现比较优越，并且受到Kitex框架支持。参考：etcd、ZookeeperConsul性能对比

- 快速部署：

  docker能够快速部署etcd，mysql，minio，nginx，redis等，操作快捷，使用方便

### 2. 架构设计
#### 2.1 系统架构设计
![框架架构](https://user-images.githubusercontent.com/52773233/220381271-269ce674-684d-4a40-8fbe-e094e79e4dda.jpg)

#### 2.2 关系型数据库设计
参考[https://github.com/siye1024/JccJ/tree/master/pic](https://github.com/siye1024/JccJ/tree/master/pic)

### 3. 项目代码介绍

直接依赖，需提前安装：
- Go 1.19.5
- Kitex v0.4.4
- libprotoc 3.21.12
- ffmpeg 5.1
- Docker 23.0.1

操作流程：
1. 在`./pkg/minio/init`中配置Minio参数 (注：ip填0则自动获取本机内网ip，公网ip请手动配置)
2. 启动docker  `docker compose up -d` 
3. 运行服务 `go build && ./dousheng` (注：仅为方便单机测试的启动方式，采用多协程，sync.WaitGroup管理)

```
├── config                     //目前只有nginx配置,计划完善配置管理
├── controller
│   ├── xhttp                  //打包请求给rpc client
│   ├── xrpc                   //rpc client发送rpc请求
├── db
├── pkg
│   ├── jwt
│   ├── minio
│   ├── pack                   //数据打包，结构转换
├── rpcserver
│   ├── kitex_gen              //Kitex生成的脚手架代码
│   ├── comment                //rpc服务端：评论，获取评论列表
│   ├── favorite               //rpc服务端：点赞，获取喜欢列表
│   ├── feed                   //rpc服务端：视频流
│   ├── publish                //rpc服务端：发布视频，获取发布列表
│   ├── relation               //rpc服务端：关注，获取关注列表
│   ├── user                   //rpc服务端：注册、登录、获取用户信息
├── docker-compose.yaml
├── main.go                    //启动服务
├── router.go                  //路由
```
## 三、测试结果
### 1. 功能测试

### 2. 性能测试

## 四、Demo演示视频

## 五、项目总结与反思
### 1.目前仍存在的问题
a. 视频发布耗时待优化
  
  视频体积较大时，发布视频的耗时较长

b. 视频点赞操作待优化
  
  用户进行视频点赞或取消点赞操作时，会先从用户喜爱的视频列表中查看用户是否点赞过该视频，在根据是否点赞过该视频的状态，判断用户后续点赞或取消点赞的合法性

### 2.已识别出的优化项
a. 视频点赞优化操作
  
  维护一个UserID-VideoID表格，表格中存储的内容为用户的UserID及用户点赞过的视频的VideoID，并对该表建立索引，以此提高查询速度

### 3. 架构演进的可能性

### 4. 项目过程中的反思与总结
a. 理解**Gorm的蛇形复数约定**含义
  
  此处踩坑，Gorm默认情况下使用结构体名的**蛇形复数**作为表名，字段名的**蛇形**作为列名，因此**数据库设计时，最好表名都采用复数形式**，不然可能导致对数据库的操作失败。

b. 善用github的分支
  
  使用github进行代码管理时，其他成员可以使用**fork**操作进行自己的开发，当对于仓库所有者，无法**fork**自己的仓库，此时应该创建分支，在分支上进行项目的开发，问就是管理规范，踩坑教训。
  

