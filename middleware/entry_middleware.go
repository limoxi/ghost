package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limoxi/ghost"
	util "github.com/limoxi/ghost/utils"
	"strings"
)

type EntryMiddleware struct{

}

func (this *EntryMiddleware) Init(){
	ghost.Info("EntryMiddleware loaded")
}

func (this *EntryMiddleware) PreRequest(ctx *gin.Context){
	// 检查是否需要经由中间价
	if strings.ToUpper(ctx.Request.Method) == "OPTIONS"{
		ctx.Set("__middleware_passed", true)
	}

	// 实现CORS
	anyHost := "*"
	corsWhiteList := ghost.Config.GetArray("cors.white_list")
	if len(corsWhiteList) > 0{
		validHost := ""
		curHost := ctx.Request.Host
		if corsWhiteList[0].(string) == anyHost{
			validHost = anyHost
		}
		if util.StringInList(curHost, corsWhiteList){
			validHost = curHost
		}
		if validHost != ""{
			ctx.Header("Access-Control-Allow-Origin", validHost)
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
			ctx.Header("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type, Accept")
		}
	}
}


func (this *EntryMiddleware) AfterResponse(ctx *gin.Context){

}