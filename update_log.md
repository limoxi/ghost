# ghost - 版本日志

### 版本号格式
```
主版本.子版本.修订版次.日期_版本阶段
```

#### 版本阶段
```
Alpha   内部开发版
Beta    公开测试
Release 正式版
```

### 升级日志
> v0.0.1.220213_Alpha
 - 支持postgresql
> v0.0.1.220114_Alpha
 - 升级gorm到v2
> v0.0.1.210828_Alpha
 - 调整结构
 - 修复接口内panic后任然继续执行框架流程的问题
> v0.0.1.210101_Alpha
 - 调整api层实现
 - 去除文件变动监听相关代码，不再支持在开发模式下的自动重启功能
> v0.0.1.200811_Alpha
 - 增加cron机制
 - 修复事务问题
 
> v0.0.1.200503_Alpha
 - rest事务控制
 - 调整GetDB
 - 调整ctx
 
> v0.0.1.200426_Alpha
 - 配置文件格式换成yaml
 - 调整分页工具
 - 一些小的改动
 
> v0.0.1.200117_Alpha
 - 增加开发模式下的文件变更监控
 - 修复一些bug
 
> v0.0.1.191216_Alpha
 - 完善各模块
    - db orm
    - GMap && Map
    - rest response
    - exception
    
> v0.0.1.191002_Alpha
 - 完善各模块
     - rest api
     - db orm
     - config
     - server
     
> v0.0.1.190803_Alpha
 - code base