package ghost

import "github.com/gin-gonic/gin"

var registeredApis = make([]apiInterface, 0)
var registeredGroupedApis = make(map[string][]apiInterface)

type apiInterface interface {

	setCtx(*gin.Context)
	setParams(RequestParams)
	GetCtx() *gin.Context
	GetParams() RequestParams

	GetResource() string
	GetLock() string

	Head() Response
	Option() Response
	Get() Response
	Put() Response
	Post() Response
	Delete() Response
}

type ApiTemplate struct{

}

func (a ApiTemplate) GetResource() string{
	panic("method not implement")
}

func (a ApiTemplate) GetLock() string{
	return ""
}

func (a ApiTemplate) Head() Response{
	panic("method not implement")
}

func (a ApiTemplate) Option() Response{
	panic("method not implement")
}

func (a ApiTemplate) Get() Response{
	panic("method not implement")
}

func (a ApiTemplate) Put() Response{
	panic("method not implement")
}

func (a ApiTemplate) Post() Response{
	panic("method not implement")
}

func (a ApiTemplate) Delete() Response{
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