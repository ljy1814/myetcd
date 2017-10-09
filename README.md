# my etcd service client

## Core
+ etcd cluster
+ etcd client
+ etcd init
    + etcd start
        - watch heartbeats
        - watch partitions
        - watch services

+ service start
    - register service
    - watch meta events
    - start http server

### 微服务系统底座
+ 一个完整的微服务系统，它的底座最少要包含以下功能：
    + 日志和审计，主要是日志的汇总，分类和查询
    + 监控和告警，主要是监控每个服务的状态，必要时产生告警
    + 消息总线，轻量级的MQ或HTTP
    + 注册发现
    + 负载均衡
    + 部署和升级
    + 事件调度机制
    + 资源管理，如：底层的虚拟机，物理机和网络管理
+ 以下功能不是最小集的一部分，但也属于底座功能：
    + 认证和鉴权
    + 微服务统一代码框架，支持多种编程语言
    + 统一服务构建和打包
    + 统一服务测试
    + 微服务CI/CD流水线
    + 服务依赖关系管理
    + 统一问题跟踪调试框架，俗称调用链
    + 灰度发布
    + 蓝绿部署
