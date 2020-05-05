package ghost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func recovery() gin.HandlerFunc{
	Info("recover func loaded...")
	return func (ctx *gin.Context){
		defer func() {
			if err := recover(); err != nil{
				errMsg := ""
				var specError *BaseError
				switch err.(type) {
				case *BaseError:
					specError = err.(*BaseError)
					errMsg = specError.ToString()
				case string:
					errMsg = err.(string)
					specError = DefaultError(errMsg)
				case *logrus.Entry:
					errMsg = err.(*logrus.Entry).Message
					specError = DefaultError(errMsg)
				default:
					errMsg = err.(error).Error()
					specError = DefaultError(errMsg)
				}
				if Config.Mode == gin.DebugMode{
					debug.PrintStack()
				}
				Error(fmt.Sprintf("recover from error: %s", errMsg))

				if idb, ok := ctx.Get("db"); ok && idb != nil{
					idb.(*gorm.DB).Rollback()
					Warn("db transaction rollback")
				}
				ctx.JSON(SERVICE_INNER_SUCCESS_CODE, specError.GetData())
			}
		}()
		ctx.Next()
	}
}