package ghost

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}else{
			log.Println("server start: ", server.Addr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	c := <- quit
	switch c {
	case syscall.SIGINT, syscall.SIGTERM:
		log.Println("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	case syscall.SIGUSR2:
		log.Println("Reload Server ...")
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

// unifyRequestParams
// 将 uri、query、postForm、file等参数合并到一处，统一通过key-value获取
func unifyRequestParams(ctx *gin.Context) RequestParam{
	return RequestParam{}
}

func bindRouter(group *gin.RouterGroup, routers []apiInterface){
	for _, r := range routers{
		group.Any(parseResource(r.GetResource()), func (ctx *gin.Context){
			unifiedParams := unifyRequestParams(ctx)
			var resp Response
			switch ctx.Request.Method {
			case "HEAD":
				resp = r.Head(ctx, unifiedParams)
			case "OPTION":
				resp = r.Option(ctx, unifiedParams)
			case "", "GET":
				resp = r.Get(ctx, unifiedParams)
			case "PUT":
				resp = r.Put(ctx, unifiedParams)
			case "POST":
				resp = r.Post(ctx, unifiedParams)
			case "DELETE":
				resp = r.Delete(ctx, unifiedParams)
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