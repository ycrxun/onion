package prome

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ycrxun/onion/metrics"
)

func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *Counter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return NewCounter(cv)
}

func NewCounter(cv *prometheus.CounterVec) *Counter {
	return &Counter{
		cv: cv,
	}
}

func (c *Counter) With(labelValues ...string) metrics.Counter {
	return &Counter{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}

func (c *Counter) Add(delta float64) {
	c.cv.With(makeLabels(c.lvs...)).Add(delta)
}
