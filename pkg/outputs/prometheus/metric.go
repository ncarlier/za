package prometheus

import (
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/expr"
	"github.com/ncarlier/za/pkg/reflection"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	Type      string
	Counter   *prometheus.CounterVec
	Condition *expr.ConditionalExpression
	Labels    map[string]string
}

func (m *Metric) IncIfMatch(event events.Event) bool {
	if m.Type == event.Type() && m.Condition.Match(event) {
		obj := event.ToMap()
		labels := make(prometheus.Labels, len(m.Labels))
		for k, v := range m.Labels {
			field := reflection.GetField(obj, v)
			field.String()
			if val, ok := field.String(); ok {
				labels[k] = val
			}
		}

		m.Counter.With(labels).Inc()
		return true
	}
	return false
}

func NewMetric(prefix string, def *MetricDef) (*Metric, error) {
	var input events.Event
	switch def.Type {
	case events.Types.PageView:
		input = &events.PageView{}
	case events.Types.Exception:
		input = &events.Exception{}
	default:
		input = &events.SimpleEvent{}
	}

	condition, err := expr.NewConditionalExpression(def.Condition, input)
	if err != nil {
		return nil, err
	}

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
		Type:      def.Type,
		Condition: condition,
		Counter:   counter,
		Labels:    def.Labels,
	}, nil
}
