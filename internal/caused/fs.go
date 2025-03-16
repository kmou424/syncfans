package caused

import "github.com/kmou424/ero"

func FileSystemError(err any) error {
	return ero.Wrap(toError(err), "file system error")
}

func IOError(err any) error {
	return ero.Wrap(toError(err), "io error")
}

func FileNotFoundError(path string) error {
	return ero.Wrap(toError(path), "file not found error")
}
