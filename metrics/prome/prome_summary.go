package prome

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ycrxun/onion/metrics"
)

func NewSummaryFrom(opts prometheus.SummaryOpts, labelNames []string) *Summary {
	sv := prometheus.NewSummaryVec(opts, labelNames)
	prometheus.MustRegister(sv)
	return NewSummary(sv)
}

func NewSummary(sv *prometheus.SummaryVec) *Summary {
	return &Summary{
		sv: sv,
	}
}

func (s *Summary) With(labelValues ...string) metrics.Histogram {
	return &Summary{
		sv:  s.sv,
		lvs: s.lvs.With(labelValues...),
	}
}

func (s *Summary) Observe(value float64) {
	s.sv.With(makeLabels(s.lvs...)).Observe(value)
}
