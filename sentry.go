package ghost

import (
	go_context "context"
	"fmt"
	"github.com/getsentry/sentry-go"

	"runtime/debug"
	"time"
	
	"os"
)

var sentryChannel = make(chan map[string]interface{}, 2048)

func isEnableSentry() bool {
	return Config.GetBool("sentry.enabled", false)
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


func init() {
	if isEnableSentry() {
		dsn := Config.GetString("sentry.dsn")
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
			AttachStacktrace: true,
		})
		if err != nil {
			Errorf("sentry.Init: %s", err)
			panic(err)
		}
		Info(fmt.Sprintf("[sentry] dsn:%s ", dsn))
	} else {
		Warn("[sentry] sentry is DISABLED!!!")
	}
}
