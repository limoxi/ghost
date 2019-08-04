package ghost

import "github.com/gin-gonic/gin"

var registeredApis = make([]apiInterface, 0)
var registeredGroupedApis = make(map[string][]apiInterface)

type apiInterface interface {
	GetResource() string
	GetLock() string

	Head(*gin.Context, RequestParam) Response
	Option(*gin.Context, RequestParam) Response
	Get(*gin.Context, RequestParam) Response
	Put(*gin.Context, RequestParam) Response
	Post(*gin.Context, RequestParam) Response
	Delete(*gin.Context, RequestParam) Response
}

type ApiTemplate struct{

}

func (a ApiTemplate) GetResource() string{
	panic("method not implement")
}

func (a ApiTemplate) GetLock() string{
	return ""
}

func (a ApiTemplate) Head(ctx *gin.Context, param RequestParam) Response{
	panic("method not implement")
}

func (a ApiTemplate) Option(ctx *gin.Context, param RequestParam) Response{
	panic("method not implement")
}

func (a ApiTemplate) Get(ctx *gin.Context, param RequestParam) Response{
	panic("method not implement")
}

func (a ApiTemplate) Put(ctx *gin.Context, param RequestParam) Response{
	panic("method not implement")
}

func (a ApiTemplate) Post(ctx *gin.Context, param RequestParam) Response{
	panic("method not implement")
}

func (a ApiTemplate) Delete(ctx *gin.Context, param RequestParam) Response{
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