package ghost

import "github.com/gin-gonic/gin"

func recovery() gin.HandlerFunc{
	Info("recover func loaded...")
	return func (ctx *gin.Context){
		defer func() {
			if err := recover(); err != nil{
				var specError *BaseError
				switch err.(type) {
				case *BaseError:
					specError = err.(*BaseError)
				case string:
					specError = DefaultError(err.(string))
				default:
					specError = DefaultError(err.(error).Error())
				}
				ctx.JSON(specError.GetCode(), specError)
			}
		}()
		ctx.Next()
	}
}