package caused

import (
	"github.com/kmou424/ero"
)

func NetError(err any) error {
	return ero.Wrap(toError(err), "net error")
}

func UrlError(err any) error {
	return ero.Wrap(toError(err), "url error")
}

func WebsocketError(err any) error {
	return ero.Wrap(toError(err), "websocket error")
}
