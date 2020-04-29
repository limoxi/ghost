# 设计理念&使用规范

## 概述
```
按照领域驱动设计，采用分层架构，分为api层、领域层和基础设施层 
```

### api层
> api层对外提供服务的能力，并对各种协议进行适配，支持基于http的restful api，
websocket和rpc调用

> api层的设计采用了一种称为ROA的设计，每一个接口都是一种资源，针对每个资源的
查、插、改、删分别对应了Get、Put、Post、Delete4种方法，当然api层还支持Head和Option
方法
例：
```go
package user

import (
	"github.com/limoxi/ghost"
)

type Users struct {
	ghost.ApiTemplate

}

func (this *Users) Resource() string{
	return "user.users"
}

func (this *Users) Post() ghost.Response{
	params := new(struct {
		Token string         `form:"token"`
		D     map[string]int `json:"d"`
	})
	this.Bind(params)
	return ghost.NewJsonResponse(map[string]interface{}{
		"id": 0,
		"name": "test",
		"token": params.Token,
		"d": params.D,
	})
}

func init(){
	ghost.RegisterApi(&Users{})
}
```
#### api的属性和方法
1. ghost.ApiTemplate 包装了api的各种方法
2. Resource() 定义了资源的路径，形式为"resource.sub_resource"，展现在url上即
/resource/sub_resource/，当然此定义支持更多层级
3. ghost.Response api支持多种响应形式，包括json、raw_string、xml、redirect，json
格式最终会被包装成以下结构
```json
{
    "code": 200,
    "data": {},
    "errCode": 200,
    "errMsg": "",
    "errStack": ""
}
```
4. 参数定义和检查使用gin的机制，ghost中使用Bind根据Content-Type自动解析参数
5. ghost.RegisterApi 此方法将资源注册到路由

### 领域层
> 领域层为服务的核心业务层，ghost中将领域对象分为领域模型和领域服务，具体细分如下图：
