package ghost

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func reloadServer(){

}

func graceRun(engine *gin.Engine) {
	watchFs()
	host := Config.GetString("web_server.host", "")
	port := Config.GetInt("web_server.port", 8080)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: engine,
	}

	Info("server start: ", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			CloseAllDBConnections()
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
		defer CloseAllDBConnections()
		if err := server.Shutdown(ctx); err != nil {
			Panic("Server Shutdown: ", err)
		}
	case syscall.SIGUSR2:
		CloseAllDBConnections()
		Info("Reload Server ...")
	}
}

// 开发环境下监听项目文件变化，触发signal
func watchFs(){
	if Config.Mode == DEV_MODE{
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
		group.Any(parseResource(r.GetResource()), func (ctx *gin.Context){
			r.setCtx(ctx)
			var resp Response
			switch ctx.Request.Method {
			case "HEAD":
				resp = r.Head()
			case "OPTION":
				resp = r.Option()
			case "", "GET":
				resp = r.Get()
			case "PUT":
				resp = r.Put()
			case "POST":
				resp = r.Post()
			case "DELETE":
				resp = r.Delete()
			default:
				panic("method not implement")
			}

			switch resp.GetDataType() {
			case "json":
				ctx.JSON(200, resp.GetData())
			default:
				ctx.String(500, "","empty response")
			}
		})
	}
}

func RunWebServer(){
	watchFs()

	engine := gin.New()

	// 应用中间件
	for _, m := range registeredMiddlewares{
		engine.Use(func() gin.HandlerFunc{
			m.Init()
			return func (ctx *gin.Context){
				m.PreRequest(ctx)
				ctx.Next()
				m.AfterResponse(ctx)
			}
		}())
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