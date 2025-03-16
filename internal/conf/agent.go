package conf

type Agent struct {
	Config struct {
		ServerAddr string `toml:"server"`
		Secret     string `toml:"secret"`
		Fan        string `toml:"fan"`

		CriticalTempRange []float64 `toml:"critical_temp_range"`
		CriticalMargin    float64   `toml:"critical_margin"`

		OverrideCurve bool    `toml:"override_curve"`
		CurveType     string  `toml:"curve_type"`
		CurveFactor   float64 `toml:"curve_factor"`
		DeadZoneRatio float64 `toml:"dead_zone_ratio"`
	} `toml:"config"`

	Sysinfo map[string]*SysInfoBase `toml:"sysinfo"`
}

type SysInfoBase struct {
	Method string `toml:"method"`
	Query  string `toml:"query"`
	Type   string `toml:"type"`
}

func GetAgentConfig() *Agent {
	return getConfig[Agent]()
}
