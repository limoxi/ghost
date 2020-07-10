package ghost

import (
	go_context "context"
	"fmt"
	//"github.com/getsentry/raven-go"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
	
	"os"
)

var sentryChannel = make(chan map[string]interface{}, 2048)

func isEnableSentry() bool {
	return Config.GetBool("sentry.enabled", false)
}

func CaptureErrorToSentry(ctx go_context.Context, err string) {
	if !isEnableSentry() {
		Info("Sentry is not enabled")
		return
	}
	
	data := make(map[string]interface{})
	data["err_msg"] = err
	data["service_name"] = Config.GetString("service_name")
	
	data["stack"] = string(debug.Stack())

	select {
	case sentryChannel <- data:
	
	case <-time.After(time.Millisecond * 500):
		Warn("[sentry] push timeout")
	}
	
}

func CaptureTaskErrorToSentry(ctx go_context.Context, errMsg string) {
	if !isEnableSentry() {
		beegoMode := os.Getenv("BEEGO_RUNMODE")
		if beegoMode == "prod" {
			Warn("Sentry is not enabled under prod mode, Please enable it!!!!")
		}
		return
	}
	
	data := make(map[string]interface{})
	data["err_msg"] = errMsg
	data["service_name"] = Config.GetString("service_name")
	
	data["stack"] = string(debug.Stack())
	
	select {
	case sentryChannel <- data:
	
	case <-time.After(time.Millisecond * 500):
		Warn("[sentry] push timeout")
	}
}

func PushErrorToSentry(errMsg string, req *http.Request) {
	if !isEnableSentry() {
		return
	}
	
	data := make(map[string]interface{})
	data["err_msg"] = errMsg
	data["service_name"] = Config.GetString("service_name")
	
	stack := string(debug.Stack())
	data["stack"] = stack
	if req != nil {
		data["raven_http"] = raven.NewHttp(req)
		data["http_request"] = req
	}
	select {
	case sentryChannel <- data:
	
	case <-time.After(time.Millisecond * 500):
		Warn("[sentry] push timeout")
	}
}

func PushErrorWithExtraDataToSentry(errMsg string, extra map[string]interface{}, req *http.Request) {
	if !isEnableSentry() {
		return
	}
	
	data := make(map[string]interface{})
	data["err_msg"] = errMsg
	data["service_name"] = Config.GetString("service_name")
	
	//stack := string(debug.Stack())
	data["stack"] = "ignore"
	data["extra"] = extra
	if req != nil {
		data["raven_http"] = raven.NewHttp(req)
		data["http_request"] = req
	}
	select {
	case sentryChannel <- data:
	
	case <-time.After(time.Millisecond * 500):
		Warn("[sentry] push timeout")
	}
}


func sendSentryPacketV2(data map[string]interface{}) {
	var packet *raven.Packet
	errMsg := data["err_msg"].(string)
	
	//封装http request
	httpRequest, ok := data["http_request"].(*http.Request)
	if ok {
		ravenHttp := raven.NewHttp(httpRequest)
		
		method := strings.ToLower(httpRequest.Method)
		if method == "post" || method == "put" || method == "delete" {
			data := make(map[string]string)
			for key, _ := range httpRequest.PostForm {
				value := httpRequest.PostForm.Get(key)
				if len(value) >= 100 {
					value = value[:100] + "..."
				}
				data[key] = value
			}
			ravenHttp.Data = data
		}
		
		packet = raven.NewPacket(errMsg, ravenHttp)
	} else {
		packet = raven.NewPacket(errMsg)
	}
	
	//确定extra
	if extra, ok := data["extra"]; ok {
		packet.Extra = extra.(map[string]interface{})
	} else {
		packet.Extra = make(map[string]interface{})
	}
	
	//确定堆栈信息
	stack, ok := data["stack"].(string)
	if !ok {
		stack = "no stack"
	}
	packet.Extra["stacktrace"] = stack
	
	//其他Tag
	tags := map[string]string{
		"service_name": data["service_name"].(string),
	}
	
	//发送给Raven
	raven.Capture(packet, tags)
}

func runSentryWorker(ch chan map[string]interface{}) {
	Info("[sentry] push-worker is ready to receive message...")
	
	for {
		data := <-sentryChannel
		
		sendSentryPacketV2(data)
	}
}

func startSentryWorker() {
	Info("[sentry] start push-worker")
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			fmt.Printf("\n>>>>>>>>>>>>>>>>>>>>\n%v\n%s\n<<<<<<<<<<<<<<<<<<<<\n", err, string(stack))
			//restart worker
			go startSentryWorker()
		}
	}()
	
	runSentryWorker(sentryChannel)
}

func init() {
	if isEnableSentry() {
		dsn := Config.GetString("sentry.dsn")
		raven.SetDSN(dsn)
		Info(fmt.Sprintf("[sentry] enable:%t, dsn:%s ", isEnableSentry(), dsn))
		go startSentryWorker()
	} else {
		Warn("[sentry] sentry is DISABLED!!!")
	}
}
