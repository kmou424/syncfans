package fantuner

import (
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/core/global"
	"golang.org/x/time/rate"
	"time"
)

type connHealthChecker struct {
	limiter *rate.Limiter
	conn    *websocket.Conn
}

func (c *connHealthChecker) CanCheck() bool {
	r := c.limiter.Reserve()
	ok := r.OK()
	r.Cancel()
	return ok
}

func (c *connHealthChecker) Check() bool {
	err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(global.WSTimeout*time.Second))
	if err != nil {
		err = caused.WebsocketError(ero.Wrap(err, "Failed to write ping frame"))
		slog.Warn(ero.AllTrace(err))
		return false
	}
	return true
}

func newConnHealthChecker(conn *websocket.Conn) *connHealthChecker {
	return &connHealthChecker{
		limiter: rate.NewLimiter(rate.Every(global.WSHealthCheckInterval*time.Second), 1),
		conn:    conn,
	}
}
