# ghost - 为gin套上工程化DDD的外衣-。-
> 使个人或团队能够更简单、快速地搭建起一个工程化的DDD架构的web服务

### 设计目标
```
在领域驱动设计的指导下按照分层架构，
实现统一的服务调用、可扩展的适配器插槽、清晰规范的领域方法论
```

### 依赖
- go ^1.13
- go mod

### 安装
```shell script
go get -u github.com/limoxi/ghost
```

### 脚手架
> 使用glo生成项目框架以及添加资源
```shell script
pip install -U git+https://github.com/limoxi/glo.git

cd ~/workspace
glo init

```

### 使用
>[设计理念&使用规范](./guide_lines.md)
>[glo: 框架生成器](https://github.com/limoxi/glo)

### 工程化目标
- 

### TODO
- [ ] 支持多种服务调用, restApi, gRPC, websocket ...
- [x] resource api设计
- [x] 中间件设计
- [x] orm（不采用依赖倒置原则进行接口设计，因为大部分场景下存储层的选择是固定的且不会改变）
- [x] 数据库事务应用策略
- [x] 配置文件设计
- [x] 异常处理
- [x] DDD规范
- [ ] 分布式锁
- [x] 日志
- [ ] docker部署
- [ ] 自动生成api文档

> [升级日志](./update_log.md)

### 项目参考
>[《实现领域驱动设计》[美] Vaughn Vernon 著；滕云 译](https://item.jd.com/11423256.html)    
>[gin](https://github.com/gin-gonic/gin)    
>[gorm](https://github.com/jinzhu/gorm)     
>[logrus](https://github.com/sirupsen/logrus)