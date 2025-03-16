package fantuner

import (
	"github.com/gorilla/websocket"
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/kmou424/syncfans/internal/proto"
	"golang.org/x/time/rate"
	"time"
)

func Hello(conn *websocket.Conn, hello *proto.ReportHello) error {
	return gFansManager.RegisterConn(conn, hello)
}

func Goodbye(fanName string) {
	gFansManager.UnregisterConn(fanName)
}

func ReportSysInfo(fanName string, sysInfo *proto.ReportSysInfo) {
	gFansManager.Tuning(fanName, sysInfo)
}

func Stop() {
	gFansManager.StopAll()
}

func Run() {
	gFansManager.Init()

	config := conf.GetServerConfig()

	limit := rate.Every(time.Duration(config.Config.Interval) * time.Millisecond)

	for fanName, fanConf := range config.Sysfans {
		gFansManager.fanInfoMap.Set(fanName, &FanInfo{
			FanParams: fanConf,
			limiter:   rate.NewLimiter(limit, 1),
		})
	}

	gFansManager.StartRoutines()
}
