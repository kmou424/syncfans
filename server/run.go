package server

import (
	"context"
	"github.com/gookit/slog"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/kmou424/syncfans/server/fantuner"
	"github.com/kmou424/syncfans/server/handler"
	"github.com/kmou424/syncfans/server/middleware"
	"github.com/labstack/echo/v4"
	"os/signal"
	"syscall"
)

func Run() {
	fantuner.Run()

	e := echo.New()

	e.Use(middleware.Recover)

	handler.NewReportHandler().Register(e)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := e.Start(conf.GetServerConfig().Config.Listen)
		if err != nil {
			panic(caused.RuntimeError(err))
		}
	}()

	<-ctx.Done()

	// on stop
	fantuner.Stop()

	slog.Info("server stopped")
}
