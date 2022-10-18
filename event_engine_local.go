package ghost

import (
	"context"
	"fmt"
	"runtime/debug"
)

type localEngine struct{}

func (this *localEngine) GetType() string {
	return "local"
}

func (this *localEngine) Emit(ctx *Context, e *Event) {
	for _, h := range getHandlersForEvent(e.Name) {
		fmt.Println(fmt.Sprintf("[event-local] %s emit", e.Name))
		go func(handler localEventHandler) {
			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
					Error(err)
				}
			}()

			err := handler.Handle(ctx, NewGMapFromData(e.data))
			if err != nil {
				Error(err)
			}
		}(h)
	}
}

func newLocalEngine() *localEngine {
	return new(localEngine)
}

// localEngine发出事件的处理器
type localEventHandler interface {
	GetEventName() string
	Handle(ctx context.Context, eventData GMap) error
}

var event2Handlers = make(map[string][]localEventHandler)

func getHandlersForEvent(eventName string) []localEventHandler {
	return event2Handlers[eventName]
}

func RegisterEventHandler(handler localEventHandler) {
	if v, ok := event2Handlers[handler.GetEventName()]; ok {
		event2Handlers[handler.GetEventName()] = append(v, handler)
	} else {
		event2Handlers[handler.GetEventName()] = []localEventHandler{handler}
	}
}

func init() {
	registerEngine(newLocalEngine())
}
