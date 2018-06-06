package prome

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ycrxun/onion/metrics"
)

type (
	Counter struct {
		cv  *prometheus.CounterVec
		lvs metrics.LabelValues
	}

	Gauge struct {
		gv  *prometheus.GaugeVec
		lvs metrics.LabelValues
	}

	Summary struct {
		sv  *prometheus.SummaryVec
		lvs metrics.LabelValues
	}

	Histogram struct {
		hv  *prometheus.HistogramVec
		lvs metrics.LabelValues
	}
)

func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
