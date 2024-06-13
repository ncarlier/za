package serializers

import (
	"fmt"

	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/serializers/json"
	"github.com/ncarlier/za/pkg/serializers/template"
)

// SerializerOutput is an interface for output plugins that are able to
// serialize events into arbitrary data formats.
type SerializerOutput interface {
	// SetSerializer sets the serializer function for the interface.
	SetSerializer(serializer Serializer)
}

// Serializer is an interface defining functions that a serializer plugin must satisfy.
type Serializer interface {
	// Serialize takes a single event and turns it into a byte buffer.
	// separate metrics should be separated by a newline, and there should be
	// a newline at the end of the buffer.
	Serialize(event events.Event) ([]byte, error)
	// ContentType returns content-type used by the serializer
	ContentType() string
}

// Config is a struct that covers the data types needed for all serializer types,
// and can be used to instantiate _any_ of the serializers.
type Config struct {
	// DataFormat can be one of the serializer types listed in NewSerializer.
	DataFormat string `toml:"data_format"`
	// DataFormatTemplate is a Golang template used by template dataformat
	// It is only used with template dataformat
	DataFormatTemplate string `toml:"data_format_template"`
}

// NewSerializer a Serializer interface based on the given config.
func NewSerializer(config *Config) (Serializer, error) {
	var err error
	var serializer Serializer
	switch config.DataFormat {
	case "json":
		serializer, err = json.NewSerializer()
	case "template":
		serializer, err = template.NewSerializer(config.DataFormatTemplate)
	default:
		err = fmt.Errorf("invalid data format: %s", config.DataFormat)
	}
	return serializer, err
}
