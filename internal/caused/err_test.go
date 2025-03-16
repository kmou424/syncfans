package caused

import (
	"errors"
	"github.com/kmou424/ero"
	"testing"
)

func TestError(t *testing.T) {
	err := errors.New("test error")
	err = NetError(ero.Wrap(err, "timeout exceeded"))
	t.Log(ero.AllTrace(err, true))
	t.Log(ero.StackTrace(err))
}
