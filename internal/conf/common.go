package conf

import (
	"strings"

	"github.com/gookit/goutil/fsutil"
	"github.com/kmou424/syncfans/internal/caused"
)

type AfterProcess interface {
	afterProcess() error
}

func parseBothFileText(text string) (string, error) {
	if !strings.HasPrefix(text, "file:") {
		return text, nil
	}
	path := strings.TrimPrefix(text, "file:")
	if !fsutil.IsFile(path) {
		return "", caused.FileNotFoundError("file not found")
	}
	content, err := fsutil.ReadStringOrErr(path)
	if err != nil {
		return "", caused.FileSystemError(err)
	}
	content = strings.TrimSpace(content)
	return content, nil
}
