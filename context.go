package ghost

import (
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type Context struct {
	mu   sync.Mutex
	data GMap
}

func (this *Context) Set(k string, v interface{}) *Context {
	this.mu.Lock()

	if this.data == nil {
		this.data = make(GMap)
	}
	this.data[k] = v

	this.mu.Unlock()

	return this
}

func (this *Context) Get(k string) interface{} {
	return this.data.Get(k)
}

func (this *Context) Data() GMap {
	return this.data
}

func (this *Context) Clone() *Context {
	newCtx := NewContext()
	newCtx.data = this.data.Clone()
	return newCtx
}

func (this *Context) GetGinCtx() *gin.Context {
	iginCtx := this.Get(_GIN_CTX_KEY)
	if iginCtx == nil {
		return nil
	}
	return iginCtx.(*gin.Context)
}

// 以下实现golang的Context接口

func (this *Context) Err() error {
	return nil
}

func (this *Context) Done() <-chan struct{} {
	return nil
}

func (this *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (this *Context) Value(ik interface{}) interface{} {
	if ik == nil {
		return nil
	}

	if key, ok := ik.(string); ok {
		return this.data.Get(key)
	}

	return nil
}

func NewContext() *Context {
	inst := new(Context)
	inst.data = make(GMap)
	return inst
}

const _GHOST_CTX_KEY = "ghostCtx"
const _GIN_CTX_KEY = "ginCtx"

func GetContextFromGinCtx(ginCtx *gin.Context) *Context {
	ighostCtx, exist := ginCtx.Get(_GHOST_CTX_KEY)
	if exist {
		return ighostCtx.(*Context)
	}
	return nil
}
