package mstorage

import (
	metric "alerting/internal/metrics"
)

type InMemoryStorage struct {
	storage map[string]metric.AbstractMetric
}

func New() *InMemoryStorage {
	str := map[string]metric.AbstractMetric{}
	var memStorage = InMemoryStorage{
		storage: str}
	return &memStorage
}

func (st *InMemoryStorage) FindAllMetrics() (map[string]metric.AbstractMetric, error) {
	return st.storage, nil
}

func (st *InMemoryStorage) CreateOrUpdateMetric(m metric.AbstractMetric) error {
	name := m.GetName()
	str := st.storage
	str[name] = m
	return nil
}

func (st *InMemoryStorage) FindMetric(name string) (metric.AbstractMetric, bool, error) {
	str := st.storage
	m, exists := str[name]
	return m, exists, nil
}
