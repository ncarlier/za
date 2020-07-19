package prometheus

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ncarlier/trackr/pkg/model"
	"github.com/ncarlier/trackr/pkg/outputs"
	"github.com/ncarlier/trackr/pkg/serializers"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var pageviewsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "website_analytics_pageviews_total",
		Help: "Number of page views.",
	},
	[]string{"tid", "hostname", "path", "isNewVisitor"},
)
var referrersCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "website_analytics_referrers_total",
		Help: "Referrers.",
	},
	[]string{"tid", "referrer"},
)

// PrometheusClient output
type PrometheusClient struct {
	Listen             string            `toml:"listen"`
	Path               string            `toml:"path"`

	server    *http.Server
	collector Collector
	url       *url.URL
	wg        sync.WaitGroup
}

type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
	Add(metrics []telegraf.Metric) error
}

var sampleConfig = `
  ## Address to listen on
  listen = ":9273"
  ## Path to publish the metrics on.
  # path = "/metrics"
`

// Connect activate the Prometheus writer
func (p *PrometheusClient) Connect() error {
	registry := prometheus.NewRegistry()
	registry.Register(pageviewsCounter)
	registry.Register(referrersCounter)
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}

	srv := &http.Server{Addr: p.Listen}
	pattern := "/metrics"
	if path != "" {
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
func (p *PrometheusClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.server.Shutdown(ctx)
	p.wg.Wait()
	p.url = nil
	prometheus.Unregister(p.collector)
	return err
}

// SampleConfig get sample configuration
func (p *PrometheusClient) SampleConfig() string {
	return sampleConfig
}

// Description get output description
func (p *PrometheusClient) Description() string {
	return "Prometheus client"
}

func (p *PrometheusClient) Write(views []*model.PageView) error {
	for _, view := range views {
		pageviewsCounter.With(prometheus.Labels{
			"tid":          view.TrackingID,
			"hostname":     view.DocumentHostName,
			"path":         view.DocumentPath,
			"isNewVisitor": strconv.FormatBool(view.IsNewVisitor),
		}).Inc()
		if view.DocumentReferrer != "" {
			referrersCounter.With(prometheus.Labels{
				"tid":      view.TrackingID,
				"referrer": view.DocumentReferrer,
			}).Inc()
		}
	}

	return null
}

func init() {
	outputs.Add("prometheus", func() outputs.Output {
		return &PrometheusClient{
			Listen:             defaultListen,
			Path:               defaultPath,
			ExpirationInterval: defaultExpirationInterval,
			StringAsLabel:      true,
		}
	})
}