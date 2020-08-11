package ghost

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func reloadServer(){

}

func beforeTerminate(){
	CloseAllDBConnections()
	if isEnableSentry(){
		sentry.Flush(2 * time.Second)
	}
}

func graceRun(engine *gin.Engine) {
	//go watchFs()
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
			Panicf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	c := <- quit
	switch c {
	case syscall.SIGINT, syscall.SIGTERM:
		Info("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer beforeTerminate()
		defer func() {
			if fsWatcher != nil{
				fsWatcher.Close()
			}
		}()
		if err := server.Shutdown(ctx); err != nil {
			Info("shutdown failed: ", err)
		}
	}
}

// 开发环境下监听项目文件变化，触发signal
func watchFs(){
	if Config.Mode == gin.DebugMode{
		watchProject()
	}
}

// parseResource
// 将resource名称格式转换为router可用的path格式
// 即 user.users ==> /user/users/
func parseResource(resource string) string{
	return fmt.Sprintf("/%s/", strings.Replace(resource, ".", "/", -1))
}

func bindRouter(group *gin.RouterGroup, routers []apiInterface){
	for _, r := range routers{
		func (ir apiInterface){
			rs := ir.Resource()
			group.Any(parseResource(rs), func (ctx *gin.Context){
				Infof("coming request %s.%s", rs, ctx.Request.Method)
				var resp Response
				var tx *gorm.DB
				if !ir.DisableTx(){
					switch ctx.Request.Method {
					case "PUT", "POST", "DELETE":
						tx = GetDB().Begin()
						Info("db transaction begin...")
						if err := tx.Error; err != nil{
							panic(err)
						}
						ctx.Set("db", tx)
					}
				}

				ir.setCtx(ctx)
				switch ctx.Request.Method {
				case "HEAD":
					resp = ir.Head()
				case "OPTIONS":
					resp = ir.Options()
				case "", "GET":
					resp = ir.Get()
				case "PUT":
					resp = ir.Put()
				case "POST":
					resp = ir.Post()
				case "DELETE":
					resp = ir.Delete()
				default:
					ctx.String(404, "","http method not implemented")
					return
				}

				if tx != nil{
					if err := tx.Commit().Error; err != nil{
						panic(err)
					}
					Info("db transaction committed...")
				}

				switch resp.GetDataType() {
				case "json":
					ctx.JSON(SERVICE_INNER_SUCCESS_CODE, resp.GetData())
				default:
					ctx.String(SERVICE_INNER_SUCCESS_CODE, "","empty response")
				}
			})
		}(r)
	}
}

func RunWebServer(){
	switch Config.Mode {

	}
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(recovery())

	// 应用中间件
	for _, m := range registeredMiddlewares{
		engine.Use(func(im middlewareInterFace) gin.HandlerFunc{
			im.Init()
			return func (ctx *gin.Context){
				im.PreRequest(ctx)
				ctx.Next()
				im.AfterResponse(ctx)
			}
		}(m))
	}

	// 应用路由、组
	for groupName, apis := range registeredGroupedApis{
		bindRouter(engine.Group(groupName), apis)
	}
	bindRouter(&engine.RouterGroup, registeredApis)

	//engine.Run()
	graceRun(engine)
}

func init(){

}