package user

import (
	"github.com/gin-gonic/gin"
	"github.com/limoxi/ghost"
)

type AUsers struct {
	ghost.ApiTemplate
}

func (this *AUsers) GetResource() string{
	return "user.users"
}

func (this *AUsers) GET(ctx *gin.Context, param ghost.Map) ghost.Response{

	token := param.GetString("token")

	return ghost.NewJsonResponse(map[string]interface{}{
		"id": 0,
		"name": "test",
		"token": token,
	})
}

func init(){
	ghost.RegisterApi(&AUsers{})
}