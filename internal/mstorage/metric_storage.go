package mstorage

import (
	metric "alerting/internal/metrics"
)

type MetricStorage interface {
	CreateOrUpdateMetric(m metric.AbstractMetric) error
	FindMetric(name string) (metric.AbstractMetric, bool, error)
	FindAllMetrics() (map[string]metric.AbstractMetric, error)
}
