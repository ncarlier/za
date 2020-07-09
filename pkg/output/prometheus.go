package output

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/model"
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

// PrometheusWriter writes data to Prometheus endpoint
type PrometheusWriter struct {
	srv *http.Server
}

// NewPrometheusOutputWriter create new Prometheus output writer
func NewPrometheusOutputWriter(uri string) (*PrometheusWriter, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil || u.Scheme != "http" {
		return nil, fmt.Errorf("invalid listen URL: %s", uri)
	}
	prometheus.MustRegister(pageviewsCounter)
	prometheus.MustRegister(referrersCounter)
	srv := &http.Server{Addr: u.Hostname() + ":" + u.Port()}
	pattern := "/metrics"
	if u.Path != "" {
		pattern = u.Path
	}
	http.Handle(pattern, promhttp.Handler())
	go func() {
		logger.Debug.Printf("starting HTTP server (%s)...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Panicln("unable to create Prometheus server endpoint:", err)
		}
	}()
	return &PrometheusWriter{
		srv: srv,
	}, nil
}

// Write writes metric to Prometheus
func (w *PrometheusWriter) Write(hit model.PageView) {
	pageviewsCounter.With(prometheus.Labels{
		"tid":          hit.TrackingID,
		"hostname":     hit.DocumentHostName,
		"path":         hit.DocumentPath,
		"isNewVisitor": strconv.FormatBool(hit.IsNewVisitor),
	}).Inc()
	if hit.DocumentReferrer != "" {
		referrersCounter.With(prometheus.Labels{
			"tid":      hit.TrackingID,
			"referrer": hit.DocumentReferrer,
		})
	}
}

// Close close the metric writer
func (w *PrometheusWriter) Close() error {
	logger.Debug.Printf("stopping HTTP server (%s)...\n", w.srv.Addr)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return w.srv.Shutdown(ctx)
}
