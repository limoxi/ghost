package ghost

import (
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/limoxi/ghost/event"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

func beforeTerminate() {
	if isEnableSentry() {
		sentry.Flush(2 * time.Second)
	}
}

func graceRun(engine *gin.Engine) {
	host := Config.GetString("web_server.host", "")
	port := Config.GetInt("web_server.port", 8080)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: engine,
	}

	Info("server start: ", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			beforeTerminate()
			panic(fmt.Sprintf("listen: %s", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	c := <-quit
	switch c {
	case syscall.SIGINT, syscall.SIGTERM:
		Info("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer beforeTerminate()
		if err := server.Shutdown(ctx); err != nil {
			Info("shutdown failed: ", err)
		}
	}
}

// parseResource
// 将resource名称格式转换为router可用的path格式
// 即 user.users ==> /user/users/
func parseResource(resource string) string {
	return fmt.Sprintf("/%s/", strings.Replace(resource, ".", "/", -1))
}

func bindRestParams(paramsName string, apiHandler apiInterface, ginContext *gin.Context) {
	elemVal := reflect.ValueOf(apiHandler).Elem()
	elemField := elemVal.FieldByName(paramsName)
	if elemField.CanSet() {
		newParamsVal := reflect.New(elemField.Type().Elem())
		bindObj := newParamsVal.Interface()

		ct := ginContext.GetHeader("Content-Type")
		var err error

		if strings.HasPrefix(ct, "application/json") {
			err = ginContext.ShouldBindJSON(bindObj)
		} else {
			err = ginContext.ShouldBind(bindObj)
		}
		if err != nil && err.Error() != "EOF" {
			panic(fmt.Sprintf("invalid params: %s", err))
		}

		elemField.Set(newParamsVal)
	}
}

func bindRouter(group *gin.RouterGroup, routers []apiInterface) {
	for _, r := range routers {
		func(ir apiInterface) {
			rs := ir.Resource()
			t := reflect.Indirect(reflect.ValueOf(ir)).Type()
			group.Any(parseResource(rs), func(ctx *gin.Context) {
				Infof("coming request %s.%s", rs, ctx.Request.Method)
				apiHandler := reflect.New(t).Interface().(apiInterface)

				var resp Response
				tx := GetDB()
				txOn := false
				if !apiHandler.DisableTx() {
					switch ctx.Request.Method {
					case "PUT", "POST", "DELETE":
						tx = tx.Begin()
						if err := tx.Error; err != nil {
							panic(err)
						}
						txOn = true
						Info("db transaction begin...")
					}
				}

				ctx.Set("db", tx)
				ctx.Set("db_tx_on", txOn)
				apiHandler.SetCtx(ctx)
				switch ctx.Request.Method {
				case "HEAD":
					resp = apiHandler.Head()
				case "OPTIONS":
					resp = apiHandler.Options()
				case "", "GET":
					bindRestParams("GetParams", apiHandler, ctx)
					resp = apiHandler.Get()
				case "PUT":
					bindRestParams("PutParams", apiHandler, ctx)
					resp = apiHandler.Put()
				case "POST":
					bindRestParams("PostParams", apiHandler, ctx)
					resp = apiHandler.Post()
				case "DELETE":
					bindRestParams("DeleteParams", apiHandler, ctx)
					resp = apiHandler.Delete()
				default:
					ctx.String(404, "", "http method not implemented")
					return
				}

				if tx != nil && txOn {
					if err := tx.Commit().Error; err != nil {
						Warn(err.Error())
						panic(err)
					}
					Info("db transaction committed...")
				}
				// 处理暂存的事件 --start
				ctx.Set("db", GetDB()) // 重置db
				event.EmitAll(ctx)
				// -- end
				if resp == nil {
					ctx.JSON(SERVICE_INNER_SUCCESS_CODE, Map{
						"code":  SERVICE_INNER_SUCCESS_CODE,
						"state": "success",
						"data":  nil,
					})
				} else {
					switch resp.GetDataType() {
					case "json":
						ctx.JSON(SERVICE_INNER_SUCCESS_CODE, resp.GetData())
					default:
						ctx.String(SERVICE_INNER_SUCCESS_CODE, "", "empty response")
					}
				}
			})
		}(r)
	}
}

func bindDocRouter(group *gin.RouterGroup) {
	group.GET("/api/doc/", func(ctx *gin.Context) {
		Info("coming request api.doc", ctx.Request.Method)
		contents, err := ioutil.ReadFile("./update_log.md")
		if err != nil {
			panic(errors.New(fmt.Sprintf("gen doc failed: %s", err.Error())))
		}

		ctx.Header("Content-Type", "text/html; charset=UTF-8")
		ctx.String(SERVICE_INNER_SUCCESS_CODE, "%s", string(contents))
	})
}

func RunWebServer() {
	engine := gin.New()
	engine.Use(recovery())

	// 应用中间件
	for _, m := range registeredMiddlewares {
		engine.Use(func(im middlewareInterFace) gin.HandlerFunc {
			im.Init()
			return func(ctx *gin.Context) {
				im.PreRequest(ctx)
				ctx.Next()
				im.AfterResponse(ctx)
			}
		}(m))
	}

	// 应用路由、组
	for groupName, apis := range registeredGroupedApis {
		bindRouter(engine.Group(groupName), apis)
	}

	defaultGroup := &engine.RouterGroup
	bindRouter(defaultGroup, registeredApis)

	// 开发环境下生成接口文档
	if Config.Mode == gin.DebugMode {
		//GenApiDoc()
	}

	//engine.Run()
	graceRun(engine)
}
