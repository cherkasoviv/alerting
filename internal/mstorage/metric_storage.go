package mstorage

import (
	metric "alerting/internal/metrics"
)

type MetricStorage interface {
	CreateOrUpdateMetric(m *metric.Metric) error
	FindMetric(name string) (*metric.Metric, bool, error)
}
