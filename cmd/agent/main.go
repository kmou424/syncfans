package main

import (
	"github.com/kmou424/syncfans/agent"
	"github.com/kmou424/syncfans/internal/caused"
)

func main() {
	defer caused.Recover(true)
	parseArgs()
	agent.Run()
}
