package config

import (
	"fmt"
	"io"
	"os"

	"dario.cat/mergo"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/serializers"
	"github.com/ncarlier/za/pkg/usage"
)

// LoadFile loads the given config file and applies it to c
func (c *Config) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	data = []byte(os.ExpandEnv(string(data)))
	root, err := toml.Parse(data)
	if err != nil {
		return err
	}

	// Parse log section table
	if err := parseSectionTable(root, "log", &c.Log); err != nil {
		return err
	}
	// Parse global section table
	if err := parseSectionTable(root, "global", &c.Global); err != nil {
		return err
	}
	// Parse Geo IP section table
	if err := parseSectionTable(root, "geo-ip", &c.GeoIP); err != nil {
		return err
	}

	// Parse trackers section table
	if err := c.parseTrackersTable(root); err != nil {
		return err
	}

	// Parse outputs section table
	if err := c.parseOutputsTable(root); err != nil {
		return err
	}

	// Apply default configuration...
	return mergo.Merge(c, NewConfig())
}

func parseSectionTable(root *ast.Table, field string, v interface{}) error {
	if val, ok := root.Fields[field]; ok {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing '%s' table", field)
		}
		if err := toml.UnmarshalTable(subTable, v); err != nil {
			return fmt.Errorf("error parsing '%s' table: %w", field, err)
		}
	}
	return nil
}

func (c *Config) parseOutputsTable(tbl *ast.Table) error {
	if val, ok := tbl.Fields["outputs"]; ok {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing outputs section")
		}
		for pluginName, pluginVal := range subTable.Fields {
			switch pluginSubTable := pluginVal.(type) {
			case []*ast.Table:
				for _, t := range pluginSubTable {
					if err := c.addOutput(pluginName, t); err != nil {
						return fmt.Errorf("error parsing %s array, %s", pluginName, err)
					}
				}
			default:
				return fmt.Errorf("unsupported output config format: %s", pluginName)
			}
		}
	}
	return nil
}

func (c *Config) parseTrackersTable(tbl *ast.Table) error {
	if val, ok := tbl.Fields["trackers"]; ok {
		subTable, ok := val.([]*ast.Table)
		if !ok {
			return fmt.Errorf("invalid configuration, error parsing trackers section")
		}
		for _, trackerTable := range subTable {
			if err := c.addTracker(trackerTable); err != nil {
				return fmt.Errorf("error parsing trackers array, %s", err)
			}
		}
	}
	return nil
}

func (c *Config) addTracker(table *ast.Table) error {
	tracker := TrackerConfig{}
	if err := toml.UnmarshalTable(table, &tracker); err != nil {
		return err
	}

	if tracker.Badge == "" {
		tracker.Badge = "ZerØ|analytics|#00a5da"
	} else if !validateBadgeSyntaxe(tracker.Badge) {
		return fmt.Errorf("invalid badge format: expecting \"<title>|<label>|<color>\" got: %s", tracker.Badge)
	}
	rateLimitingConfig := usage.RateLimitingConfig{}
	if err := parseSectionTable(table, "rate_limiting", &rateLimitingConfig); err != nil {
		return fmt.Errorf("error parsing rate limiting section, %s", err)
	}
	rateLimiter, err := usage.NewRateLimiter(rateLimitingConfig)
	if err != nil {
		return err
	}
	tracker.RateLimiter = rateLimiter

	c.Trackers = append(c.Trackers, tracker)
	return nil
}

func (c *Config) addOutput(name string, table *ast.Table) error {
	creator, ok := outputs.Outputs[name]
	if !ok {
		return fmt.Errorf("undefined but requested output: %s", name)
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

	c.Outputs = append(c.Outputs, output)
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
