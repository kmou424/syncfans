package init

import (
	"github.com/gookit/slog"
	"github.com/kmou424/syncfans/internal/core/global"
)

func initLogger() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = false
		f.TimeFormat = "2006-01-02 15:04:05.000"
		if !global.Debug() {
			f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n")
			slog.SetLogLevel(slog.InfoLevel)
		} else {
			f.SetTemplate("[{{datetime}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n")
			slog.SetLogLevel(slog.DebugLevel)
		}
	})
}
