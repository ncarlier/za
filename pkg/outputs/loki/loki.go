package loki

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/logger"
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/outputs/loki/logproto"
	"github.com/ncarlier/za/pkg/serializers"
)

const maxEntriesChanSize = 5000

// Loki output
type Loki struct {
	URL           string          `toml:"url"`
	Timeout       config.Duration `toml:"timeout"`
	BatchSize     int             `toml:"batch_size"`
	BatchInterval config.Duration `toml:"batch_interval"`

	client     *Client
	serializer serializers.Serializer

	quit      chan struct{}
	entries   chan events.Event
	waitGroup sync.WaitGroup
}

var sampleConfig = `
  ## Loki URL
  url = "http://localhost:3001"
  ## Timeout
	timeout = "3s"
	## Batch interval
	batch_interval = "5s"
	## Batch size
	batch_size = 100
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

	l.quit = make(chan struct{})
	l.entries = make(chan events.Event, maxEntriesChanSize)

	go l.run()

	logger.Debug.Printf("using LOKI output: %s\n", u.String())

	return nil
}

// Close the output writer
func (l *Loki) Close() error {
	close(l.quit)
	l.waitGroup.Wait()
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

// SendEvent send event to the Output
func (l *Loki) SendEvent(event events.Event) error {
	l.entries <- event
	return nil
}

func (l *Loki) run() {
	l.waitGroup.Add(1)
	var batch []events.Event
	batchSize := 0
	maxWait := time.NewTimer(l.BatchInterval.Duration)

	defer func() {
		if batchSize > 0 {
			l.write(batch)
		}
		l.waitGroup.Done()
	}()

	for {
		select {
		case <-l.quit:
			return
		case entry := <-l.entries:
			batch = append(batch, entry)
			batchSize++
			if batchSize >= l.BatchSize {
				if err := l.write(batch); err != nil {
					logger.Error.Printf("unable to send batch of page view to Loki (%s): %v\n", l.URL, err)
				}
				batch = []events.Event{}
				batchSize = 0
				maxWait.Reset(l.BatchInterval.Duration)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				if err := l.write(batch); err != nil {
					logger.Error.Printf("unable to send batch of page view to Loki (%s): %v\n", l.URL, err)
				}
				batch = []events.Event{}
				batchSize = 0
			}
			maxWait.Reset(l.BatchInterval.Duration)
		}
	}
}

func (l *Loki) write(entries []events.Event) error {
	batch := map[string]*logproto.Stream{}

	for _, event := range entries {
		line, err := l.serializer.Serialize(event)
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
	return l.client.Send(streams)
}

func init() {
	outputs.Add("loki", func() outputs.Output {
		return &Loki{
			URL:           "http://localhost:3100",
			Timeout:       config.Duration{Duration: time.Second * 2},
			BatchInterval: config.Duration{Duration: 10 * time.Second},
			BatchSize:     10,
		}
	})
}
