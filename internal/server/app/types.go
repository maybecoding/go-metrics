package app

type (
	TypeGauge   float64
	TypeCounter int64

	MetricGauge struct {
		Name  string
		Value TypeGauge
	}

	MetricCounter struct {
		Name  string
		Value TypeCounter
	}

	Metric struct {
		Type  string
		Name  string
		Value string
	}
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)
