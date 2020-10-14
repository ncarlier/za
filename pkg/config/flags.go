package config

// Flags contain global configuration
type Flags struct {
	ListenAddr string `flag:"listen-addr" desc:"HTTP listen address" default:":8080"`
	ConfigFile string `flag:"config-file" desc:"Config file" default:"zero-analytics.toml"`
	Debug      bool   `flag:"debug" desc:"Output debug logs" default:"false"`
}
