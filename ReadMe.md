#Dousheng
## 一、项目介绍
- 项目简介：基于kitex RPC微服务 + gin HTPP服务完成的第五届字节跳动青训营——极简抖音后端项目
- 项目服务地址：
- Github地址：https://github.com/siye1024/JccJ/

## 二、项目分工
| 团队成员 |  主要贡献 |
| --- | --- |
| 夏志良 | 部署基于kitex的RPC框架，实现ETCD微服务相关功能, 负责服务器的环境搭建与配置, 负责开发publish模块、feed模块 |
| 李思怡 | 负责开发user模块、comment模块、relation模块、favorite模块，负责数据库的设计和github仓库的管理 |

## 三、项目实现
#### 1. 项目特点
1. 采用RPC框架（kitex）脚手架生成代码进行开发，基于RPC微服务 + gin提供HTTP服务
2. 基于[极简版抖音-各功能对应的接口说明文档](https://www.apifox.cn/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c)提供的接口进行开发，使用[极简抖音APP使用说明](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)提供的APK进行Demo测试，功能**完整实现**，前端接口匹配良好
3. 代码结构采用（HTTP API层 + RPC Service层 + Dal层）结构，项目结构清晰，代码符合规范
4. 使用JWT进行用户Token验证
5. 使用ETCD进行服务发现于服务注册
6. 使用Minio实现视频文件和图片的对象存储
7. 使用Gorm对Mysql进行ORM操作
8. 数据库表建立了外键约束，并使用事务对数据库进行操作，保证数据的一致性和安全性

#### 2. 架构设计

#### 3. 项目代码介绍

## 四、测试结果
#### 1. 功能测试

#### 2. 性能测试

## 五、Demo演示视频

## 六、项目总结与反思
#### 1.目前仍存在的问题
#### 2.已识别出的优化项
#### 3. 架构演进的可能性
#### 4. 项目过程中的反思与总结
