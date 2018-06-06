package prome

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ycrxun/onion/metrics"
)

func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *Histogram {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(hv)
	return NewHistogram(hv)
}

func NewHistogram(hv *prometheus.HistogramVec) *Histogram {
	return &Histogram{
		hv: hv,
	}
}
func (h *Histogram) With(labelValues ...string) metrics.Histogram {
	return &Histogram{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
}

func (h *Histogram) Observe(value float64) {
	h.hv.With(makeLabels(h.lvs...)).Observe(value)
}
