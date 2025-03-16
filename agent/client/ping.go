package client

import (
	"github.com/gorilla/websocket"
	"github.com/kmou424/syncfans/internal/core/global"
	"time"
)

var getPingHandler = func(conn *websocket.Conn) func(appData string) error {
	return func(appData string) error {
		return conn.WriteControl(websocket.PingMessage, []byte(appData), time.Now().Add(time.Second*global.WSTimeout))
	}
}
