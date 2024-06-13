package prometheus

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/outputs"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Client output
type Client struct {
	Listen  string      `toml:"listen"`
	Path    string      `toml:"path"`
	Prefix  string      `toml:"prefix"`
	Metrics []MetricDef `toml:"metrics"`

	sync.Mutex
	server   *http.Server
	wg       sync.WaitGroup
	registry *prometheus.Registry
	metrics  []*Metric
}

func (p *Client) registerCounters() error {
	p.registry = prometheus.NewRegistry()
	p.metrics = make([]*Metric, 0, len(p.Metrics))
	for _, mDef := range p.Metrics {
		metric, err := NewMetric(p.Prefix, &mDef)
		if err != nil {
			return err
		}
		p.metrics = append(p.metrics, metric)
		p.registry.Register(metric.Counter)
		slog.Debug("new metric registered", "name", mDef.Name, "type", mDef.Type, "condition", mDef.Condition)
	}
	return nil
}

func (p *Client) startServer() {
	promHandler := promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})

	srv := &http.Server{Addr: p.Listen}
	pattern := "/metrics"
	if p.Path != "" {
		pattern = p.Path
	}
	http.Handle(pattern, promHandler)

	slog.Debug("starting HTTP server...", "addr", srv.Addr)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("unable to create Prometheus server endpoint", "error", err)
			os.Exit(1)
		}
	}()
}

// Connect activate the Prometheus writer
func (p *Client) Connect() error {
	p.registerCounters()
	p.startServer()
	return nil
}

// Close the output writer
func (p *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.server.Shutdown(ctx)
	p.wg.Wait()

	// unregister metrics
	for _, metric := range p.metrics {
		p.registry.Unregister(metric.Counter)
	}
	return err
}

// SendEvent send event to the Output
func (p *Client) SendEvent(event events.Event) error {
	p.Lock()
	defer p.Unlock()

	for _, metric := range p.metrics {
		metric.IncIfMatch(event)
	}

	return nil
}

func init() {
	outputs.Add("prom", func() outputs.Output {
		return &Client{
			Listen: ":9213",
			Path:   "/metrics",
			Prefix: "za_",
		}
	})
}
