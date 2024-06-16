package prometheus

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ncarlier/za/pkg/conditional"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/outputs"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Output for Prometheus
type Output struct {
	Listen  string      `toml:"listen"`
	Path    string      `toml:"path"`
	Prefix  string      `toml:"prefix"`
	Metrics []MetricDef `toml:"metrics"`

	condition conditional.Expression

	sync.Mutex
	server   *http.Server
	wg       sync.WaitGroup
	registry *prometheus.Registry
	metrics  []*Metric
}

func (o *Output) registerCounters() error {
	o.registry = prometheus.NewRegistry()
	o.metrics = make([]*Metric, 0, len(o.Metrics))
	for _, mDef := range o.Metrics {
		metric, err := NewMetric(o.Prefix, &mDef)
		if err != nil {
			return err
		}
		o.metrics = append(o.metrics, metric)
		o.registry.Register(metric.Counter)
		slog.Debug("new metric registered", "name", mDef.Name, "type", mDef.Type)
	}
	return nil
}

func (o *Output) startServer() {
	promHandler := promhttp.HandlerFor(o.registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})

	srv := &http.Server{Addr: o.Listen}
	pattern := "/metrics"
	if o.Path != "" {
		pattern = o.Path
	}
	http.Handle(pattern, promHandler)

	slog.Debug("starting HTTP server...", "addr", srv.Addr)

	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("unable to create Prometheus server endpoint", "error", err)
			os.Exit(1)
		}
	}()
}

// SetCondition set condition expression
func (o *Output) SetCondition(condition conditional.Expression) {
	o.condition = condition
}

// Connect activate the Prometheus writer
func (o *Output) Connect() error {
	o.registerCounters()
	o.startServer()
	return nil
}

// Close the output writer
func (o *Output) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := o.server.Shutdown(ctx)
	o.wg.Wait()

	// unregister metrics
	for _, metric := range o.metrics {
		o.registry.Unregister(metric.Counter)
	}
	return err
}

// SendEvent send event to the Output
func (o *Output) SendEvent(event events.Event) error {
	if !o.condition.Match(event) {
		return nil
	}

	o.Lock()
	defer o.Unlock()

	for _, metric := range o.metrics {
		metric.IncIfMatch(event)
	}

	return nil
}

func init() {
	outputs.Add("prom", func() outputs.Output {
		return &Output{
			Listen: ":9213",
			Path:   "/metrics",
			Prefix: "za_",
		}
	})
}
