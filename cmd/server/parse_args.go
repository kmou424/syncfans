package main

import (
	"flag"
	"fmt"
	"github.com/gookit/goutil/sysutil"
	"github.com/kmou424/syncfans/internal/conf"
	_ "github.com/kmou424/syncfans/internal/core/init"
)

var (
	configFile string
)

var configSearchPaths = []string{
	"./server.toml",
	fmt.Sprintf("%s/.config/syncfans/server.toml", sysutil.HomeDir()),
	"/etc/syncfans/server.toml",
}

func parseArgs() {
	flag.StringVar(&configFile, "config", "", "config file path")
	flag.Parse()

	postParse()
}

func postParse() {
	var (
		err error
	)
	if configFile == "" {
		err = conf.AutoLoad[conf.Server](configSearchPaths)
		if err != nil {
			panic(err)
		}
	} else {
		err = conf.AutoLoad[conf.Server]([]string{configFile})
		if err != nil {
			panic(err)
		}
	}
}
