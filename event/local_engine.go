package event

import (
	"context"
	"fmt"
	"runtime/debug"
)

// IRestContext 解决循环依赖
type IRestContext interface {
	AddEvent(eventData map[string]interface{})
}

type localEngine struct{}

func (this *localEngine) GetType() string {
	return "local"
}

func (this *localEngine) Emit(ctx context.Context, e *Event) {
	for _, h := range getHandlersForEvent(e.Name) {
		fmt.Println(fmt.Sprintf("[event-local] %s emit", e.Name))
		go func(handler localEventHandler) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(string(debug.Stack()))
				}
			}()

			err := handler.Handle(ctx, e.data)
			if err != nil {
				fmt.Println(err)
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
	Handle(ctx context.Context, eventData map[string]interface{}) error
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
