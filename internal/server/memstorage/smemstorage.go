package smemstorage

type MemStorage struct {
	dataGauge   map[string]float64
	dataCounter map[string]int64
}

func (m *MemStorage) SetMetricGauge(metricName string, metricValue float64) {
	m.dataGauge[metricName] = metricValue
}

func (m *MemStorage) SetMetricCounter(metricName string, metricValue int64) {
	m.dataCounter[metricName] += metricValue
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		dataGauge:   make(map[string]float64),
		dataCounter: make(map[string]int64),
	}
}
