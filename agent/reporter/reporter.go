package reporter

import (
	"context"
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/agent/client"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/kmou424/syncfans/internal/proto"
	"time"
)

type Reporter struct {
	client    *client.Client
	helloSent bool
}

func NewReporter(client *client.Client) (*Reporter, error) {
	if client == nil {
		return nil, caused.RuntimeError(ero.New("client is nil"))
	}
	if !client.IsAlive() {
		return nil, caused.NetError(ero.New("client is dead"))
	}
	return &Reporter{
		client: client,
	}, nil
}

func (r *Reporter) Run(ctx context.Context) {
	defer func() {
		caused.Recover(false)
		slog.Info("reporter is stopped...")
	}()
	// every 500ms, limit to 1 report
	ticker := time.NewTicker(500 * time.Millisecond)

	slog.Info("reporter is running...")
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !r.checkClient() {
				continue
			}
			r.Report()
		}
	}
}

func (r *Reporter) checkClient() bool {
	if !r.client.IsAlive() {
		// TODO: reconnect

		// should resend hello after reconnect
		r.helloSent = false
		return false
	}
	return true
}

func (r *Reporter) Report() {
	r.onReportHello()
	r.onReportSysInfo()
}

func (r *Reporter) onReportHello() {
	if !r.helloSent {
		config := conf.GetAgentConfig().Config
		hello := &proto.ReportHello{
			Fan:               config.Fan,
			CriticalTempRange: config.CriticalTempRange,
			CriticalMargin:    config.CriticalMargin,
			OverrideCurve:     config.OverrideCurve,
			CurveType:         config.CurveType,
			CurveFactor:       config.CurveFactor,
			DeadZoneRatio:     config.DeadZoneRatio,
		}
		err := r.client.Send(hello)
		if err != nil {
			panic(caused.WebsocketError(ero.Wrap(err, "failed to send hello to server")))
		}
		r.helloSent = true
	}
}

func (r *Reporter) onReportSysInfo() {
	sysInfo, err := GetSysInfo()
	if err != nil {
		panic(caused.RuntimeError(ero.Wrap(err, "failed to get system info")))
	}
	reportSysInfo, err := sysInfo.Report()
	if err != nil {
		panic(caused.RuntimeError(ero.Wrap(err, "failed to convert system info to report format")))
	}
	err = r.client.Send(reportSysInfo)
	if err != nil {
		panic(caused.WebsocketError(ero.Wrap(err, "failed to send system info to server")))
	}
}
