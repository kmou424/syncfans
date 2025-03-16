package caused

import (
	"github.com/kmou424/ero"
	"runtime"
)

func funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

func toError(err any) error {
	if err, ok := err.(error); ok {
		return err
	}
	return ero.Newf("%v", err)
}
