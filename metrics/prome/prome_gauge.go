package prome

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ycrxun/onion/metrics"
)

func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return NewGauge(gv)
}

func NewGauge(gv *prometheus.GaugeVec) *Gauge {
	return &Gauge{
		gv: gv,
	}
}

func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}

func (g *Gauge) Set(value float64) {
	g.gv.With(makeLabels(g.lvs...)).Set(value)
}

func (g *Gauge) Add(delta float64) {
	g.gv.With(makeLabels(g.lvs...)).Add(delta)
}
