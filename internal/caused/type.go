package caused

import "github.com/kmou424/ero"

func TypeError(err any) error {
	return ero.Wrap(toError(err), "type error")
}

func ValueError(err any) error {
	return ero.Wrap(toError(err), "value error")
}

func InvalidTypeError(err any) error {
	return ero.Wrap(toError(err), "invalid type error")
}
