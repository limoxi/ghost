package ghost

import "github.com/gin-gonic/gin"

var registeredMiddlewares  = make([]middlewareInterFace, 0)

type middlewareInterFace interface {
	Init()
	PreRequest(*gin.Context)
	AfterResponse(*gin.Context)
}

func RegisterMiddleware(m middlewareInterFace){
	registeredMiddlewares = append(registeredMiddlewares, m)
}