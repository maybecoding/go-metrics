package app

type Storage interface {
	SetMetricGauge(gauge *MetricGauge)
	SetMetricCounter(counter *MetricCounter)

	GetMetricGauge(name string) (TypeGauge, error)
	GetMetricCounter(name string) (TypeCounter, error)

	GetMetricGaugeAll() []*MetricGauge
	GetMetricCounterAll() []*MetricCounter
}
