package main

import (
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/server"
)

func main() {
	defer caused.Recover(true)
	parseArgs()
	server.Run()
}
