package middleware

import (
	"strings"

	"github.com/limoxi/ghost"
	util "github.com/limoxi/ghost/utils"
)

type EntryMiddleware struct {
}

func (this *EntryMiddleware) Init() {
	ghost.Info("EntryMiddleware loaded")
}

func (this *EntryMiddleware) PreRequest(ctx *ghost.Context) {
	// 检查是否需要经由中间件
	ginCtx := ctx.GetGinCtx()
	if ginCtx == nil {
		return
	}
	if strings.ToUpper(ginCtx.Request.Method) == "OPTIONS" {
		ctx.SetMiddlewareIgnored()
	}

	// 实现CORS
	anyHost := "*"
	corsWhiteList := ghost.Config.GetArray("web.cors.white_list")
	if len(corsWhiteList) > 0 {
		validHost := ""
		curHost := ginCtx.Request.Host
		if corsWhiteList[0].(string) == anyHost {
			validHost = anyHost
		}
		if util.NewLister(corsWhiteList).Has(curHost) {
			validHost = curHost
		}
		if validHost != "" {
			ginCtx.Header("Access-Control-Allow-Origin", validHost)
			ginCtx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
			ginCtx.Header("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type, Accept")
		}
	}

	// 不走中间件的资源
	ignoredResources := ghost.Config.GetArray("web.middleware.ignored_resources")
	fullPath := strings.ToLower(ginCtx.FullPath())
	curResource := strings.Trim(fullPath, "/")
	if util.NewLister(ignoredResources).Has(curResource) {
		ctx.SetMiddlewareIgnored()
	}
}

func (this *EntryMiddleware) AfterResponse(ctx *ghost.Context) {

}
