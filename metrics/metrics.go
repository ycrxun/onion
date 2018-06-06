package metrics

type (
	Counter interface {
		With(labelValues ...string) Counter
		Add(delta float64)
	}

	Gauge interface {
		With(labelValues ...string) Gauge
		Set(value float64)
		Add(delta float64)
	}

	Summary interface {
		With(labelValues ...string) Summary
		Observe(value float64)
	}

	Histogram interface {
		With(labelValues ...string) Histogram
		Observe(value float64)
	}

	LabelValues []string

	Metrics struct {
		Counter   Counter
		Gauge     Gauge
		Summary   Summary
		Histogram Histogram
	}
)

func (lvs LabelValues) With(labelValues ...string) LabelValues {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}
	return append(lvs, labelValues...)
}
