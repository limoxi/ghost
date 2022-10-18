package ghost

var registeredMiddlewares = make([]middlewareInterFace, 0)

type middlewareInterFace interface {
	Init()
	PreRequest(*Context)
	AfterResponse(*Context)
}

func RegisterMiddleware(m middlewareInterFace) {
	registeredMiddlewares = append(registeredMiddlewares, m)
}
