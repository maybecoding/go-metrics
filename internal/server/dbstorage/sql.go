package dbstorage

import _ "embed"

//go:embed sql/set_metric_gauge.sql
var sqlSetMetricGauge string

//go:embed sql/set_metric_counter.sql
var sqlSetMetricCounter string

//go:embed sql/get_metric.sql
var sqlGetMetric string

//go:embed sql/get_metric_all.sql
var sqlGetMetricAll string
