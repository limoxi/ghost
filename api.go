package ghost

import (
	"context"
)

var registeredApis = make([]apiInterface, 0)
var registeredGroupedApis = make(map[string][]apiInterface)

type apiInterface interface {
	SetCtx(ctx context.Context)
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

type ApiTemplate struct {
	ContextObject
	GMap
}

func (a *ApiTemplate) Resource() string {
	panic("method not implement")
}

func (a *ApiTemplate) GetLock() string {
	return ""
}

func (a *ApiTemplate) DisableTx() bool {
	return false
}

func (a *ApiTemplate) Head() Response {
	panic("method not implement")
}

func (a *ApiTemplate) Options() Response {
	return NewRawResponse("")
}

func (a *ApiTemplate) Get() Response {
	panic("method not implement")
}

func (a *ApiTemplate) Put() Response {
	panic("method not implement")
}

func (a *ApiTemplate) Post() Response {
	panic("method not implement")
}

func (a *ApiTemplate) Delete() Response {
	panic("method not implement")
}

func RegisterApi(r apiInterface) {
	registeredApis = append(registeredApis, r)
}

func RegisterGroupedApi(groupName string, r apiInterface) {
	if g, ok := registeredGroupedApis[groupName]; ok && len(g) > 0 {
		registeredGroupedApis[groupName] = append(registeredGroupedApis[groupName], r)
	} else {
		registeredGroupedApis[groupName] = []apiInterface{r}
	}
}

func getAllApis() []apiInterface {
	apis := append([]apiInterface{}, registeredApis...)
	for _, as := range registeredGroupedApis {
		apis = append(apis, as...)
	}
	return apis
}
