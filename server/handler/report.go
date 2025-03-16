package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/core/global"
	"github.com/kmou424/syncfans/internal/proto"
	"github.com/kmou424/syncfans/server/fantuner"
	"github.com/kmou424/syncfans/server/middleware"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync/atomic"
	"time"
)

type ReportHandler struct {
}

func NewReportHandler() IHandler {
	return &ReportHandler{}
}

func (h *ReportHandler) Register(e *echo.Echo) {
	e.GET("/report", h.Report, middleware.Auth)
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * global.WSTimeout,
	ReadBufferSize:   global.WSReadBufferSize,
	WriteBufferSize:  global.WSReadBufferSize,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func (h *ReportHandler) Report(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			panic(caused.WebsocketError(err))
		}
	}()

	var fanName string
	var fanRegistered = atomic.Bool{}
	defer func() {
		if fanRegistered.Load() {
			fantuner.Goodbye(fanName)
		}
	}()
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return caused.WebsocketError(err)
		}
		if msgType != websocket.TextMessage {
			continue
		}
		message := &proto.Message{}
		err = json.Unmarshal(msg, message)
		if err != nil {
			return caused.WebsocketError(err)
		}

		switch message.Type {
		case proto.MsgTypeHello:
			if fanRegistered.Load() {
				return caused.RuntimeError("agent already registered")
			}
			hello := &proto.ReportHello{}
			err = json.Unmarshal(message.Payload, hello)
			if err != nil {
				return caused.WebsocketError(err)
			}
			fanName = hello.Fan
			err = fantuner.Hello(ws, hello)
			if err != nil {
				return caused.WebsocketError(err)
			}
			fanRegistered.Store(true)
		case proto.MsgTypeSysInfo:
			if !fanRegistered.Load() {
				return caused.RuntimeError("agent not registered")
			}
			sysInfo := &proto.ReportSysInfo{}
			err = json.Unmarshal(message.Payload, sysInfo)
			if err != nil {
				return caused.WebsocketError(err)
			}
			fantuner.ReportSysInfo(fanName, sysInfo)
		}
	}
}
