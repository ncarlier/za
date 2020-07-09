package config

// Config contain global configuration
type Config struct {
	ListenAddr string `flag:"listen-addr" desc:"HTTP listen address" default:":8080"`
	Debug      bool   `flag:"debug" desc:"Output debug logs" default:"false"`
}
