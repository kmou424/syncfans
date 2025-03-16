package fantuner

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/kmou424/syncfans/internal/proto"
	"github.com/kmou424/syncfans/pkg/genericmap"
	"sync"
	"time"
)

type fansManager struct {
	wg sync.WaitGroup

	fanInfoMap    *genericmap.SafeMap[string, *FanInfo]
	cancelFuncMap *genericmap.SafeMap[string, context.CancelFunc]
}

var gFansManager = &fansManager{}

func (r *fansManager) Init() {
	r.wg = sync.WaitGroup{}
	r.fanInfoMap = genericmap.NewSafeMap[string, *FanInfo]()
	r.cancelFuncMap = genericmap.NewSafeMap[string, context.CancelFunc]()
}

func (r *fansManager) StopAll() {
	for _, routineName := range r.cancelFuncMap.Keys() {
		r.Stop(routineName)
	}
	r.wg.Wait()
}

func (r *fansManager) Stop(routineName string) {
	cancelFunc, ok := r.cancelFuncMap.Get(routineName)
	if ok {
		slog.Infof("[%s] exiting...", routineName)
		r.cancelFuncMap.Delete(routineName)
		cancelFunc()
	} else {
		slog.Infof("[%s] exited", routineName)
	}
}

func (r *fansManager) RegisterConn(conn *websocket.Conn, hello *proto.ReportHello) (err error) {
	fanName := hello.Fan
	fanInfo, ok := r.fanInfoMap.Get(fanName)
	if !ok {
		return caused.RuntimeError(fmt.Sprintf("fan [%s] not found", fanName))
	}
	if fanInfo.conn != nil {
		return caused.RuntimeError(fmt.Sprintf("fan [%s] already has a connection", fanName))
	}
	fanInfo.conn = conn

	healthChecker := newConnHealthChecker(conn)
	fanInfo.healthChecker = healthChecker

	defer func() {
		if err != nil {
			fanInfo.Clean()
		}
	}()

	if len(hello.CriticalTempRange) != 2 {
		return caused.ValueError(fmt.Sprintf("invalid config.critical_temp_range: %v", hello.CriticalTempRange))
	}
	if hello.CriticalTempRange[0] >= hello.CriticalTempRange[1] {
		return caused.ValueError(fmt.Sprintf("invalid config.critical_temp_range: %v", hello.CriticalTempRange))
	}
	fanInfo.criTemps = hello.CriticalTempRange

	if hello.CriticalMargin < 0 {
		return caused.ValueError(fmt.Sprintf("invalid config.critical_margin: %v", hello.CriticalMargin))
	}
	fanInfo.criMargin = hello.CriticalMargin

	config := conf.GetServerConfig()
	curveType := config.Default.CurveType
	curveFactor := config.Default.CurveFactor
	deadZoneRatio := config.Default.DeadZoneRatio
	if hello.OverrideCurve {
		curveType = hello.CurveType
		curveFactor = hello.CurveFactor
		deadZoneRatio = hello.DeadZoneRatio
	}
	switch curveType {
	case "linear", "s-curve", "exponential", "aggressive":
		break
	default:
		return caused.ValueError(fmt.Sprintf("invalid config.curve_type: %v", curveType))
	}
	if curveFactor < 0 {
		return caused.ValueError(fmt.Sprintf("invalid config.curve_factor: %v", curveFactor))
	}
	if deadZoneRatio < 0 || deadZoneRatio > 1 {
		return caused.ValueError(fmt.Sprintf("invalid config.dead_zone_ratio: %v", deadZoneRatio))
	}
	fanInfo.curveType = curveType
	fanInfo.curveFactor = curveFactor
	fanInfo.deadZoneRatio = deadZoneRatio

	return nil
}

func (r *fansManager) UnregisterConn(fanName string) {
	info, ok := r.fanInfoMap.Get(fanName)
	if !ok {
		return
	}
	if info.conn == nil {
		return
	}
	info.Clean()
}

func (r *fansManager) HealthCheckRoutine() {
	r.wg.Add(1)
	defer r.wg.Done()

	const RoutineName = "HealthCheckRoutine"
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFuncMap.Set(RoutineName, cancel)

	timer := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			r.Stop(RoutineName)
			return
		case <-timer.C:
			deadFans := make([]string, 0)
			r.fanInfoMap.ForEach(func(fanName string, info *FanInfo) bool {
				if info.healthChecker == nil {
					return true
				}
				checker := info.healthChecker
				if !checker.CanCheck() {
					return true
				}
				var ok bool
				for i := 0; i < 3; i++ { // retry 3 times
					ok = checker.Check()
					if ok {
						break
					}
					time.Sleep(time.Microsecond * 100)
				}
				if !ok {
					deadFans = append(deadFans, fanName)
				}
				return true
			})
			for _, fanName := range deadFans {
				r.UnregisterConn(fanName)
				slog.Warnf("[%s] connection lost, unregistering", fanName)
			}
		}
	}
}

func (r *fansManager) StartRoutines() {
	go r.HealthCheckRoutine()
}
