package conf

import (
	"reflect"

	"github.com/gookit/goutil/fsutil"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/pelletier/go-toml/v2"
)

var config any = nil

func getConfig[T any]() *T {
	if config == nil {
		panic(caused.RuntimeError("config not loaded"))
	}
	ret, ok := config.(T)
	if !ok {
		panic(caused.TypeError("config type not match"))
	}
	return &ret
}

func AutoLoad[T any](searchPaths []string) error {
	cfg, err := Load[T](searchPaths)
	if err != nil {
		return err
	}
	config = cfg
	return nil
}

func Load[T any](searchPaths []string) (*T, error) {
	cfgType := reflect.TypeFor[T]()
	cfg := ptr(reflect.New(cfgType).Elem().Interface().(T))

	for _, searchPath := range searchPaths {
		if !fsutil.IsFile(searchPath) {
			continue
		}
		err := loadConfig(searchPath, cfg)
		if err != nil {
			return cfg, err
		}
		return cfg, nil
	}

	return cfg, caused.FileNotFoundError("config file not found")
}

func loadConfig[T any](filePath string, cfg *T) error {
	content, err := fsutil.ReadOrErr(filePath)
	if err != nil {
		return caused.FileSystemError(err)
	}
	err = toml.Unmarshal(content, cfg)
	if err != nil {
		return caused.ValueError(err)
	}
	return nil
}

func ptr[T any](v T) *T { return &v }
