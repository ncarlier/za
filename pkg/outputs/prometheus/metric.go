package prometheus

import (
	"log/slog"

	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/reflection"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	Type    string
	Counter *prometheus.CounterVec
	Labels  map[string]string
}

func (m *Metric) IncIfMatch(event events.Event) bool {
	if m.Type == event.Type() {
		obj := event.ToMap()
		labels := make(prometheus.Labels, len(m.Labels))
		for k, v := range m.Labels {
			field := reflection.GetField(obj, v)
			field.String()
			if val, ok := field.String(); ok {
				labels[k] = val
			}
		}

		if counter, err := m.Counter.GetMetricWith(labels); err == nil {
			counter.Inc()
			return true
		} else {
			slog.Error("invalid metric definition, please check your labels configuration", "error", err)
		}
	}
	return false
}

func NewMetric(prefix string, def *MetricDef) (*Metric, error) {
	// TODO add global labels
	labelNames := make([]string, 0, len(def.Labels))
	for label := range def.Labels {
		labelNames = append(labelNames, label)
	}

	metricName := prefix + def.Name
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricName,
			Help: def.Help,
		},
		labelNames,
	)
	return &Metric{
		Type:    def.Type,
		Counter: counter,
		Labels:  def.Labels,
	}, nil
}
