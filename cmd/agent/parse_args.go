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
	"./agent.toml",
	fmt.Sprintf("%s/.config/syncfans/agent.toml", sysutil.HomeDir()),
	"/etc/syncfans/agent.toml",
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
		err = conf.AutoLoad[conf.Agent](configSearchPaths)
		if err != nil {
			panic(err)
		}
	} else {
		err = conf.AutoLoad[conf.Agent]([]string{configFile})
		if err != nil {
			panic(err)
		}
	}
}
