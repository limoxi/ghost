package ghost

import "github.com/gorilla/websocket"

const WS_MSG_PING = "ping"
const WS_MSG_PONG = "pong"

type WsMessage struct {
	Topic string `json:"topic"`
}

type iWebsocketHandler interface {
	Name() string
	GetUpgrader() *websocket.Upgrader
	Receive(content string)
	Emit(content string)
}

var registeredWSHandler = make(map[string]iWebsocketHandler)

func RegisterWSHandler(iwh iWebsocketHandler) {
	Info("register ws= ======>", iwh.Name())
	registeredWSHandler[iwh.Name()] = iwh
}

type WebsocketHandler struct {
}

func (this *WebsocketHandler) Name() string {
	panic("not_implemented")
}

func (this *WebsocketHandler) GetUpgrader() *websocket.Upgrader {
	return nil
}

func (this *WebsocketHandler) Receive(content string) {
	panic("not_implemented")
}
