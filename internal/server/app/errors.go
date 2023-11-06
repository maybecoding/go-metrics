package app

import "errors"

var (
	ErrMetricTypeIncorrect = errors.New("metric type incorrect")
	ErrNoMetricValue       = errors.New("no metric value")
)
