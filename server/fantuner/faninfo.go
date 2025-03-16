package fantuner

import (
	"github.com/gorilla/websocket"
	"github.com/kmou424/syncfans/internal/conf"
	"golang.org/x/time/rate"
	"sync"
)

const SlidingWindowSize = 5

type FanInfo struct {
	lk sync.Mutex

	criTemps  []float64 // criTemps: The range of critical temperatures: from hello message.
	criMargin float64   // criMargin: The critical temperature margin: from hello message.

	// fan parameters
	*conf.FanParams

	// curve parameters
	curveType     string  // Speed curve algorithm
	curveFactor   float64 // Curve shape modifier
	deadZoneRatio float64 // Speed change tolerance

	// cached values
	smoothedSpeed float64                    // The smoothed speed from the last update.
	lastTarget    int                        // The last target speed from the last update.
	tempHist      [SlidingWindowSize]float64 // The history of temperatures.

	// connection
	conn          *websocket.Conn
	healthChecker *connHealthChecker

	limiter *rate.Limiter // Limit the maximum rate for updating the fan status.
}

func (fi *FanInfo) Clean() {
	fi.smoothedSpeed = 0
	fi.tempHist = [5]float64{}
	fi.lastTarget = 0

	fi.conn = nil
	fi.healthChecker = nil
}
