package global

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/sysutil"
	"sync"
	"sync/atomic"
)

var (
	debug atomic.Bool
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		debug.Store(func() bool {
			ret, _ := goutil.ToBool(sysutil.Getenv("SYNCFANS_DEBUG", "false"))
			return ret
		}())
	})
}

func Debug() bool {
	return debug.Load()
}

func SetDebug(b bool) {
	debug.Store(b)
}
