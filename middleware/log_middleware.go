package middleware

import (
	"github.com/gin-gonic/gin"
)

type LogMiddleware struct{

}

func (this *LogMiddleware) Init(){

}

func (this *LogMiddleware) PreRequest(ctx *gin.Context){

}

func (this *LogMiddleware) AfterResponse(ctx *gin.Context){

}