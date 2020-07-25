package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/ncarlier/trackr/pkg/model"
	"github.com/ncarlier/trackr/pkg/outputs"
	"github.com/ncarlier/trackr/pkg/serializers"
)

// Config is the root of the configuration
type Config struct {
	Global   GlobalConfig
	WebSites []model.WebSite
	Outputs  []*outputs.Output
}

// GlobalConfig is the global section fo the configuration
type GlobalConfig struct {
	GeoIPDatabase string `toml:"geo_ip_database"`
	BatchSize     int
	BatchInterval Duration
}

// NewConfig create new configuration
func NewConfig() *Config {
	c := &Config{
		Global: GlobalConfig{
			BatchSize:     10,
			BatchInterval: Duration{Duration: 10 * time.Second},
		},
		WebSites: make([]model.WebSite, 0),
		Outputs:  make([]*outputs.Output, 0),
	}
	return c
}

// ValidateTrackingID validate that origin matches with the tracking ID
func (c *Config) ValidateTrackingID(origin, trackingID string) bool {
	for _, website := range c.WebSites {
		if strings.HasPrefix(origin, website.Origin) && website.TrackingID == trackingID {
			return true
		}
	}
	return false
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tbl, err := toml.Parse(data)
	if err != nil {
		return err
	}

	// Parse global table:
	if val, ok := tbl.Fields["global"]; ok {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing global table")
		}
		if err = toml.UnmarshalTable(subTable, &c.Global); err != nil {
			return fmt.Errorf("error parsing global table: %w", err)
		}
	}

	// Parse websites table:
	if val, ok := tbl.Fields["websites"]; ok {
		subTable, ok := val.([]*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing websites table")
		}
		for _, websiteTable := range subTable {
			if err = c.addWebsite(websiteTable); err != nil {
				return fmt.Errorf("Error parsing website array, %s", err)
			}
		}
		delete(tbl.Fields, "websites")
	}

	// Parse rest
	for name, val := range tbl.Fields {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing field %q as table", name)
		}

		switch name {
		case "global":
		case "outputs":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addOutput(pluginName, t); err != nil {
							return fmt.Errorf("Error parsing %s array, %s", pluginName, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s",
						pluginName)
				}
			}
		default:
			return fmt.Errorf("Error parsing %s, %s", name, err)
		}
	}
	return nil
}

func (c *Config) addWebsite(table *ast.Table) error {
	website := model.WebSite{}
	if err := toml.UnmarshalTable(table, &website); err != nil {
		return err
	}

	c.WebSites = append(c.WebSites, website)
	return nil
}

func (c *Config) addOutput(name string, table *ast.Table) error {
	creator, ok := outputs.Outputs[name]
	if !ok {
		return fmt.Errorf("Undefined but requested output: %s", name)
	}
	output := creator()

	// If the output has a SetSerializer function, then this means it can write
	// arbitrary types of output, so build the serializer and set it.
	switch t := output.(type) {
	case serializers.SerializerOutput:
		serializer, err := buildSerializer(name, table)
		if err != nil {
			return err
		}
		t.SetSerializer(serializer)
	}

	if err := toml.UnmarshalTable(table, output); err != nil {
		return err
	}

	c.Outputs = append(c.Outputs, &output)
	return nil
}

func buildSerializer(name string, tbl *ast.Table) (serializers.Serializer, error) {
	c := &serializers.Config{}

	if node, ok := tbl.Fields["data_format"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				c.DataFormat = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["data_format_template"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				c.DataFormatTemplate = str.Value
			}
		}
	}

	if c.DataFormat == "" {
		c.DataFormat = "json"
	}

	delete(tbl.Fields, "data_format")
	delete(tbl.Fields, "data_format_template")
	return serializers.NewSerializer(c)
}
