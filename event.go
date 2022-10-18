package ghost

import (
	"context"
	"fmt"
)

var type2Engine map[string]iEngine
var validEngines []iEngine

type iEngine interface {
	GetType() string
	Emit(ctx *Context, e *Event)
}

func registerEngine(eg iEngine) {
	if type2Engine == nil {
		type2Engine = make(map[string]iEngine)
	}
	if _, ok := type2Engine[eg.GetType()]; !ok {
		type2Engine[eg.GetType()] = eg
		validEngines = append(validEngines, eg)
		fmt.Printf("event engine: %s registered", eg.GetType())
	}
}

func getValidEngines() (ies []iEngine) {
	return validEngines
}

type Event struct {
	Name string
	Tag  string

	data Map
}

func NewEvent(name, tag string) *Event {
	event := new(Event)
	event.Name = name
	event.Tag = tag
	return event
}

// Emit 将事件暂存到ctx中，等当前请求完成后再触发
func Emit(ctx context.Context, e *Event, data Map) {
	ghostCtx := ctx.(*Context)
	var existEvents []*Event
	if ies := ghostCtx.Get("events"); ies != nil {
		existEvents = ies.([]*Event)
		for _, ee := range existEvents {
			if ee.Name == e.Name { // 防止已存在的事件多次触发
				return
			}
		}
	}

	e.data = data
	existEvents = append(existEvents, e)
	ghostCtx.Set("events", existEvents)
}

// EmitAll 触发暂存的事件
func EmitAll(ghostCtx *Context) {
	if ies := ghostCtx.Get("events"); ies != nil {
		ghostCtx.Set("db", GetDB())
		for _, e := range ies.([]*Event) {
			for _, ve := range getValidEngines() {
				ve.Emit(ghostCtx.Clone(), e)
			}
		}
	}
}
