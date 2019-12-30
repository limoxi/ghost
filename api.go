package ghost

import (
	"github.com/gin-gonic/gin"
)

var registeredApis = make([]apiInterface, 0)
var registeredGroupedApis = make(map[string][]apiInterface)

type apiInterface interface {

	setCtx(*gin.Context)
	GetCtx() *gin.Context

	GetResource() string
	GetLock() string

	Head() Response
	Options() Response
	Get() Response
	Put() Response
	Post() Response
	Delete() Response
}

type ApiTemplate struct{
	ctx *gin.Context
}

func (a *ApiTemplate) setCtx(ctx *gin.Context) {
	a.ctx = ctx
}

func (a *ApiTemplate) GetCtx() *gin.Context {
	return a.ctx
}

// 绑定参数到struct
func (a *ApiTemplate) Bind(obj interface{}){
	ginContext := a.GetCtx()
	ct := ginContext.GetHeader("Content-Type")
	var err error
	switch ct {
	case "application/json", "application/json;charset=UTF-8":
		err = ginContext.ShouldBindJSON(obj)
	case "application/xml":
		err = ginContext.ShouldBindXML(obj)
	case "application/x-www-form-urlencoded":
		err = ginContext.ShouldBind(obj)
	default:
		Infof("unhandled Content-Type: %s", ct)
	}
	if err != nil{
		Panicf("invalid params: %s", err)
	}
}

func (a *ApiTemplate) GetResource() string{
	panic("method not implement")
}

func (a *ApiTemplate) GetLock() string{
	return ""
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