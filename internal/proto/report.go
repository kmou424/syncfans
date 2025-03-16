package proto

import (
	"encoding/json"
)

const (
	ReportVersion = 1
)

type ReportHello struct {
	Version           int       `json:"version"`
	Fan               string    `json:"fan"`
	CriticalTempRange []float64 `json:"critical_temp_range"`
	CriticalMargin    float64   `json:"critical_margin"`
	OverrideCurve     bool      `json:"override_curve"`
	CurveType         string    `toml:"curve_type"`
	CurveFactor       float64   `toml:"curve_factor"`
	DeadZoneRatio     float64   `toml:"dead_zone_ratio"`
}

func (r ReportHello) MarshalJSON() ([]byte, error) {
	type alias ReportHello
	v := alias(r)
	v.Version = ReportVersion
	return json.Marshal(v)
}

type ReportSysInfo struct {
	// V1
	Temperature float64 `json:"temperature"`
	Usage       float64 `json:"usage"`
}
