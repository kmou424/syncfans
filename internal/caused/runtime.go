package caused

import "github.com/kmou424/ero"

func RuntimeError(err any) error {
	return ero.Wrap(toError(err), "runtime error")
}
