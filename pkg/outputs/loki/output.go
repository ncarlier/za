package loki

import (
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/ncarlier/za/pkg/conditional"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/outputs/loki/logproto"
	"github.com/ncarlier/za/pkg/serializers"
)

const maxEntriesChanSize = 5000

// Output for Loki
type Output struct {
	URL           string          `toml:"url"`
	Timeout       config.Duration `toml:"timeout"`
	BatchSize     int             `toml:"batch_size"`
	BatchInterval config.Duration `toml:"batch_interval"`

	client     *Client
	serializer serializers.Serializer
	condition  conditional.Expression

	quit      chan struct{}
	entries   chan events.Event
	waitGroup sync.WaitGroup
}

// SetSerializer set data serializer
func (o *Output) SetSerializer(serializer serializers.Serializer) {
	o.serializer = serializer
}

// SetCondition set condition expression
func (o *Output) SetCondition(condition conditional.Expression) {
	o.condition = condition
}

// Connect activate the Loki writer
func (o *Output) Connect() error {
	u, err := url.Parse(o.URL)
	if err != nil {
		return fmt.Errorf("invalid Loki URL: %v", err)
	}
	u.Path = "/loki/api/v1/push"
	cfg := Config{
		URL:     u.String(),
		Timeout: o.Timeout.Duration,
	}
	o.client = NewClient(cfg)

	o.quit = make(chan struct{})
	o.entries = make(chan events.Event, maxEntriesChanSize)

	go o.run()

	slog.Debug("using LOKI output", "uri", u.String())

	return nil
}

// Close the output writer
func (o *Output) Close() error {
	close(o.quit)
	o.waitGroup.Wait()
	return nil
}

// Description get output description
func (o *Output) Description() string {
	return "Loki client"
}

// SendEvent send event to the Output
func (o *Output) SendEvent(event events.Event) error {
	if !o.condition.Match(event) {
		return nil
	}
	o.entries <- event
	return nil
}

func (o *Output) run() {
	o.waitGroup.Add(1)
	var batch []events.Event
	batchSize := 0
	maxWait := time.NewTimer(o.BatchInterval.Duration)

	defer func() {
		if batchSize > 0 {
			o.write(batch)
		}
		o.waitGroup.Done()
	}()

	for {
		select {
		case <-o.quit:
			return
		case entry := <-o.entries:
			batch = append(batch, entry)
			batchSize++
			if batchSize >= o.BatchSize {
				if err := o.write(batch); err != nil {
					slog.Error("unable to send batch of page view to Loki", "uri", o.URL, "error", err)
				}
				batch = []events.Event{}
				batchSize = 0
				maxWait.Reset(o.BatchInterval.Duration)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				if err := o.write(batch); err != nil {
					slog.Error("unable to send batch of page view to Loki", "uri", o.URL, "error", err)
				}
				batch = []events.Event{}
				batchSize = 0
			}
			maxWait.Reset(o.BatchInterval.Duration)
		}
	}
}

func (o *Output) write(entries []events.Event) error {
	batch := map[string]*logproto.Stream{}

	for _, event := range entries {
		line, err := o.serializer.Serialize(event)
		if err != nil {
			return err
		}
		labels := event.Labels().String()
		stream, ok := batch[labels]
		if !ok {
			stream = &logproto.Stream{
				Labels: labels,
			}
			batch[labels] = stream
		}
		entry := logproto.Entry{
			Timestamp: event.TS(),
			Line:      string(line),
		}
		stream.Entries = append(stream.Entries, entry)
	}
	streams := []*logproto.Stream{}
	for _, stream := range batch {
		streams = append(streams, stream)
	}
	return o.client.Send(streams)
}

func init() {
	outputs.Add("loki", func() outputs.Output {
		return &Output{
			URL:           "http://localhost:3100",
			Timeout:       config.Duration{Duration: time.Second * 2},
			BatchInterval: config.Duration{Duration: 10 * time.Second},
			BatchSize:     10,
		}
	})
}
