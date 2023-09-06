package mstorage

import (
	metric "alerting/cmd/server/metrics"
)

type MetricStorage interface {
	CreateOrUpdateMetric(m *metric.Metric) error
	FindMetric(name string) (*metric.Metric, bool, error)
}
