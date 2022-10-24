package ghost

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"runtime/debug"
)

func rollbackTxInCtx(ctx context.Context) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ighostCtx, _ := ginCtx.Get("ghostCtx")
		ctx = ighostCtx.(*Context)
	}

	var txOn bool
	var db interface{}
	if itx := ctx.Value("db_tx_on"); itx != nil && itx.(bool) {
		txOn = true
		if idb := ctx.Value("db"); idb != nil {
			db = idb
		}
	}

	if txOn && db != nil {
		db.(*gorm.DB).Rollback()
		Warn("db transaction rollback")
	}
}

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

				rollbackTxInCtx(ctx)
				ctx.JSON(SERVICE_INNER_SUCCESS_CODE, specError.GetData())
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func RecoverFromPanic(ctx context.Context, scene string) {
	if err := recover(); err != nil {
		Error(string(debug.Stack()))
		Warn(fmt.Sprintf("recover from %s panic...", scene), err)
		rollbackTxInCtx(ctx)

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
