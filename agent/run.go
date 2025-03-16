package agent

import (
	"context"
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/agent/client"
	"github.com/kmou424/syncfans/agent/reporter"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/conf"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := conf.GetAgentConfig().Config

	slog.Infof("connecting to server at %s...", config.ServerAddr)
	c, err := client.NewClient(&client.Option{
		Url:    config.ServerAddr,
		Secret: config.Secret,
	})
	if err != nil {
		err = caused.RuntimeError(ero.Wrap(err, "failed to create client"))
		panic(err)
	}

	slog.Info("server connection established")

	slog.Info("starting reporter...")

	var r *reporter.Reporter
	r, err = reporter.NewReporter(c)
	if err != nil {
		err = caused.RuntimeError(ero.Wrap(err, "failed to create reporter"))
		slog.Error(ero.AllTrace(err))
	}
	r.Run(ctx)

	return
}
