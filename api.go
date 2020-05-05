package ghost

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

var registeredApis = make([]apiInterface, 0)
var registeredGroupedApis = make(map[string][]apiInterface)

type apiInterface interface {

	setCtx(ctx context.Context)
	GetCtx() context.Context

	Resource() string
	GetLock() string
	DisableTx() bool

	Head() Response
	Options() Response
	Get() Response
	Put() Response
	Post() Response
	Delete() Response
}

type ApiTemplate struct{
	ctx context.Context
	GMap
}

func (a *ApiTemplate) setCtx(ctx context.Context) {
	a.ctx = ctx
}

func (a *ApiTemplate) GetCtx() context.Context {
	return a.ctx
}

// 绑定参数到struct
func (a *ApiTemplate) Bind(obj interface{}){
	ginContext := a.GetCtx().(*gin.Context)
	ct := ginContext.GetHeader("Content-Type")
	var err error
	if ginContext.Request.Method == "GET"{
		err = ginContext.ShouldBind(obj)
	}else{
		switch ct {
		case "application/json":
			fallthrough
		case "application/json;charset=utf-8", "application/json;charset=UTF-8":
			err = ginContext.ShouldBindJSON(obj)
		default:
			Infof("coming request Content-Type: %s", ct)
			err = ginContext.ShouldBind(obj)
		}
	}

	if err != nil{
		panic(fmt.Sprintf("invalid params: %s", err))
	}
}

func (a *ApiTemplate) Resource() string{
	panic("method not implement")
}

func (a *ApiTemplate) GetLock() string{
	return ""
}

func (a *ApiTemplate) DisableTx() bool{
	return false
}

func (a *ApiTemplate) Head() Response{
	panic("method not implement")
}

func (a *ApiTemplate) Options() Response{
	return NewRawResponse("")
}

func (a *ApiTemplate) Get() Response{
	panic("method not implement")
}

func (a *ApiTemplate) Put() Response{
	panic("method not implement")
}

func (a *ApiTemplate) Post() Response{
	panic("method not implement")
}

func (a *ApiTemplate) Delete() Response{
	panic("method not implement")
}

func RegisterApi(r apiInterface){
	registeredApis = append(registeredApis, r)
}

func RegisterGroupedApi(groupName string, r apiInterface){
	if g, ok := registeredGroupedApis[groupName]; ok && len(g) > 0{
		registeredGroupedApis[groupName] = append(registeredGroupedApis[groupName], r)
	}else{
		registeredGroupedApis[groupName] = []apiInterface{r}
	}
}