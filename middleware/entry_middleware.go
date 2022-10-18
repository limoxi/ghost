package middleware

import (
	"github.com/limoxi/ghost"
	util "github.com/limoxi/ghost/utils"
	"strings"
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
		ctx.Set("__middleware_passed", true)
	}

	// 实现CORS
	anyHost := "*"
	corsWhiteList := ghost.Config.GetArray("cors.white_list")
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
}

func (this *EntryMiddleware) AfterResponse(ctx *ghost.Context) {

}
