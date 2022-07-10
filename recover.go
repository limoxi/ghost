package ghost

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"runtime/debug"
)

func recovery() gin.HandlerFunc {
	Info("recover func loaded...")
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
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
				debug.PrintStack()
				Error(fmt.Sprintf("recover from panic: %s", errMsg))

				if itx, ok := ctx.Get("db_tx_on"); ok && itx.(bool) {
					if idb, ok := ctx.Get("db"); ok && idb != nil {
						idb.(*gorm.DB).Rollback()
						Warn("db transaction rollback")
					}
				}
				ctx.JSON(SERVICE_INNER_SUCCESS_CODE, specError.GetData())
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func RecoverFromCronTaskPanic(ctx context.Context) {
	db := GetDBFromCtx(ctx)
	if err := recover(); err != nil {
		Error(string(debug.Stack()))
		Warn("recover from cron task panic...", err)
		if db != nil {
			db.Rollback()
			Warn("[ORM] rollback transaction for cron task")
		}

		{
			// 推送日志到sentry
			errMsg := err.(error).Error()
			if be, ok := err.(*BaseError); ok {
				errMsg = fmt.Sprintf("%s - %s", be.ErrCode, be.ErrMsg)
			}
			Error(errMsg)
			//CaptureTaskErrorToSentry(ctx, errMsg)
		}
	}
}
