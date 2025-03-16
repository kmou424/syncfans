package conf

type Server struct {
	Config struct {
		Listen          string  `toml:"listen"`
		Secret          string  `toml:"secret"`
		Interval        int     `toml:"interval"`
		SmoothingFactor float64 `toml:"smoothing_factor"`
	} `toml:"config"`

	Default FanCurve `toml:"default"`

	Sysfans map[string]*FanParams `toml:"sysfans"`
}

type FanCurve struct {
	CurveType     string  `toml:"curve_type"`
	CurveFactor   float64 `toml:"curve_factor"`
	DeadZoneRatio float64 `toml:"dead_zone_ratio"`
}

type FanParams struct {
	Path     string `toml:"path"`
	MaxSpeed int    `toml:"max_speed"`
	MinSpeed int    `toml:"min_speed"`
}

func GetServerConfig() *Server {
	return getConfig[Server]()
}
