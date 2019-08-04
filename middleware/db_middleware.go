package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/limoxi/ghost"
)

// DbMiddleware 使用gorm
// 根据api的配置，决定开启或关闭事务，默认(非GET请求)开启
type DbMiddleware struct{
	db *gorm.DB
}

func (this *DbMiddleware) Init(){

}

func (this *DbMiddleware) PreRequest(ctx *gin.Context){
	db, err := gorm.Open(
		ghost.Config.GetString("database.engine"),
		fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
			ghost.Config.GetString("user"),
			ghost.Config.GetString("password"),
			ghost.Config.GetString("name"),
		),
	)
	if err != nil{
		panic("connect mysql failed")
	}
	this.db = db
	ctx.Set("db", db)
}

func (this *DbMiddleware) AfterResponse(ctx *gin.Context){

	defer this.db.Close()
}