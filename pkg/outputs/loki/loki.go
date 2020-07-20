package loki

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/model"
	"github.com/ncarlier/trackr/pkg/outputs"
	"github.com/ncarlier/trackr/pkg/outputs/loki/logproto"
	"github.com/ncarlier/trackr/pkg/serializers"
)

// Loki output
type Loki struct {
	URL     string          `toml:"url"`
	Timeout config.Duration `toml:"timeout"`

	client     *Client
	serializer serializers.Serializer
}

var sampleConfig = `
  ## Loki URL
  url = "http://localhost:3001"
  ## Timeout
  timeout = "3s"
`

// SetSerializer set data serializer
func (l *Loki) SetSerializer(serializer serializers.Serializer) {
	l.serializer = serializer
}

// Connect activate the Loki writer
func (l *Loki) Connect() error {
	u, err := url.Parse(l.URL)
	if err != nil {
		return fmt.Errorf("invalid Loki URL: %v", err)
	}
	u.Path = "/loki/api/v1/push"
	cfg := Config{
		URL:     u.String(),
		Timeout: l.Timeout.Duration,
	}
	l.client = NewClient(cfg)
	return nil
}

// Close the output writer
func (l *Loki) Close() error {
	return nil
}

// SampleConfig get sample configuration
func (l *Loki) SampleConfig() string {
	return sampleConfig
}

// Description get output description
func (l *Loki) Description() string {
	return "Loki client"
}

func (l *Loki) Write(views []*model.PageView) error {
	batch := map[string]*logproto.Stream{}

	for _, view := range views {
		line, err := l.serializer.Serialize(*view)
		if err != nil {
			return err
		}
		labels := view.Labels().String()
		stream, ok := batch[labels]
		if !ok {
			stream = &logproto.Stream{
				Labels: labels,
			}
			batch[labels] = stream
		}
		entry := logproto.Entry{
			Timestamp: view.Timestamp,
			Line:      string(line),
		}
		stream.Entries = append(stream.Entries, entry)
	}
	streams := []*logproto.Stream{}
	for _, stream := range batch {
		streams = append(streams, stream)
	}
	return l.client.Send(streams)
}

func init() {
	outputs.Add("loki", func() outputs.Output {
		return &Loki{
			URL:     "http://localhost:3100",
			Timeout: config.Duration{Duration: time.Second * 2},
		}
	})
}
