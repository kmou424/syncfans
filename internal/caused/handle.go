package caused

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/internal/core/global"
)

func Recover(exit bool) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			panic(r)
		}
		if global.Debug() && ero.IsEro(err) {
			trace := ero.AllTrace(err, true)
			slog.Error(trace)
		} else {
			debug.PrintStack()
			slog.Error(fmt.Sprintf("panic: %v", err))
		}
		if exit {
			os.Exit(1)
		}
	}
}
