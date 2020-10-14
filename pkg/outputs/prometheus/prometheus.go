package prometheus

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/logger"
	"github.com/ncarlier/za/pkg/outputs"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Client output
type Client struct {
	Listen string `toml:"listen"`
	Path   string `toml:"path"`

	sync.Mutex
	server           *http.Server
	pageviewsCounter *prometheus.CounterVec
	referersCounter  *prometheus.CounterVec
	wg               sync.WaitGroup
}

var sampleConfig = `
  ## Address to listen on
  listen = ":9213"
  ## Path to publish the metrics on.
  # path = "/metrics"
`

func newPageviewsCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "website_analytics_pageviews_total",
			Help: "Number of page views.",
		},
		[]string{"tid", "hostname", "path", "isNewVisitor"},
	)
}

func newReferersCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "website_analytics_referers_total",
			Help: "Referers.",
		},
		[]string{"tid", "referer"},
	)
}

// Connect activate the Prometheus writer
func (p *Client) Connect() error {
	registry := prometheus.NewRegistry()
	p.pageviewsCounter = newPageviewsCounter()
	registry.Register(p.pageviewsCounter)
	p.referersCounter = newReferersCounter()
	registry.Register(p.referersCounter)
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})

	srv := &http.Server{Addr: p.Listen}
	pattern := "/metrics"
	if p.Path != "" {
		pattern = p.Path
	}
	http.Handle(pattern, promHandler)

	logger.Debug.Printf("starting HTTP server (%s)...\n", srv.Addr)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Panicln("unable to create Prometheus server endpoint:", err)
		}
	}()

	return nil
}

// Close the output writer
func (p *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.server.Shutdown(ctx)
	p.wg.Wait()
	prometheus.Unregister(p.pageviewsCounter)
	prometheus.Unregister(p.referersCounter)
	return err
}

// SampleConfig get sample configuration
func (p *Client) SampleConfig() string {
	return sampleConfig
}

// Description get output description
func (p *Client) Description() string {
	return "Prometheus client"
}

// SendEvent send event to the Output
func (p *Client) SendEvent(event events.Event) error {
	p.Lock()
	defer p.Unlock()

	labels := make(prometheus.Labels, len(event.Labels()))
	for k, v := range event.Labels() {
		labels[k] = v
	}

	p.pageviewsCounter.With(labels).Inc()

	return nil
}

func init() {
	outputs.Add("prometheus", func() outputs.Output {
		return &Client{
			Listen: ":9213",
			Path:   "/metrics",
		}
	})
}
