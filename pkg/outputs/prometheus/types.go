package prometheus

// MetricDef is a metric definition
type MetricDef struct {
	Name      string            `toml:"name"`
	Help      string            `toml:"help"`
	Type      string            `toml:"type"`
	Condition string            `toml:"condition"`
	Labels    map[string]string `toml:"labels"`
}
