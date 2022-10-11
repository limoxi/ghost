package event

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Event struct {
	Name string
	Tag  string

	data map[string]interface{}
}

func NewEvent(name, tag string) *Event {
	event := new(Event)
	event.Name = name
	event.Tag = tag
	return event
}

// Emit 将事件暂存到ctx中，等当前请求完成后再触发
func Emit(ctx context.Context, e *Event, data map[string]interface{}) {
	var existEvents []*Event
	ginCtx := ctx.(*gin.Context)
	if ie, ok := ginCtx.Get("events"); ok {
		existEvents = ie.([]*Event)
	}

	e.data = data
	existEvents = append(existEvents, e)
	ginCtx.Set("events", existEvents)
}

// EmitAll 触发暂存的事件
func EmitAll(ctx context.Context) {
	ginCtx := ctx.(*gin.Context)
	if ies, ok := ginCtx.Get("events"); ok {
		for _, e := range ies.([]*Event) {
			for _, ve := range getValidEngines() {
				ve.Emit(ginCtx.Copy(), e)
			}
		}
	}
}
