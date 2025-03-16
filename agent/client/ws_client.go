package client

import (
	"context"
	"encoding/json"
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/core/global"
	"github.com/kmou424/syncfans/internal/proto"
	"sync/atomic"
	"time"
)

type Client struct {
	alive atomic.Bool
	opt   *Option
	conn  *websocket.Conn
}

type Option struct {
	Url    string
	Secret string
}

func NewClient(opt *Option) (*Client, error) {
	c := &Client{
		opt: opt,
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) IsAlive() bool {
	return c.alive.Load()
}

func (c *Client) connect() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: global.WSTimeout * time.Second,
		ReadBufferSize:   global.WSReadBufferSize,
		WriteBufferSize:  global.WSWriteBufferSize,
	}

	conn, _, err := dialer.Dial(c.opt.Url, map[string][]string{
		"Authorization": {c.opt.Secret},
	})
	if err != nil {
		return caused.WebsocketError(err)
	}

	conn.SetPingHandler(getPingHandler(conn))
	c.alive.Store(true)
	go c.keepAlive()

	c.conn = conn
	return nil
}

func (c *Client) Reconnect(ctx context.Context) {
	var err error
	for i := 1; i <= 3; i++ {
		err = c.connect()
		if err == nil {
			slog.Info("client reconnected")
			break
		}
		slog.Warnf("failed to reconnect client, retrying... (attempt %d)", i)
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		err = ero.Wrap(err, "failed to reconnect client")
		slog.Error(ero.AllTrace(err))
		ctx.Done()
	}
	return
}

func (c *Client) keepAlive() {
	defer func() {
		if c.conn != nil {
			_ = c.conn.Close()
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			var err error
			for i := 0; i <= 3; i++ {
				err := c.conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(global.WSTimeout*time.Second))
				if err != nil {
					c.alive.Store(false)
					if i > 0 {
						slog.Warnf("failed to write ping frame, retrying %d...", i)
					}
					time.Sleep(time.Second)
					continue
				}
				c.alive.Store(true)
				break
			}
			if err != nil {
				c.alive.Store(false)
				err = caused.WebsocketError(ero.Wrap(err, "failed to write ping frame"))
				slog.Error(ero.AllTrace(err))
				ticker.Stop()
				return
			}
		}
	}
}

func (c *Client) Send(msg any) error {
	message := &proto.Message{}
	switch msg.(type) {
	case *proto.ReportHello:
		message.Type = proto.MsgTypeHello
	case *proto.ReportSysInfo:
		message.Type = proto.MsgTypeSysInfo
	default:
		return caused.TypeError("unsupported message type")
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return caused.ValueError(err)
	}
	message.Payload = bytes

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return caused.ValueError(err)
	}
	err = c.conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		return caused.WebsocketError(err)
	}
	return nil
}
