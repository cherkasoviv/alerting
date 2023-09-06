package mstorage

import (
	metric "alerting/cmd/server/metrics"
	"fmt"
)

type InMemoryStorage struct {
	Storage map[string]metric.AbstractMetric
}

func (st *InMemoryStorage) CreateOrUpdate(m metric.AbstractMetric) error {
	name := m.GetName()
	fmt.Println(name)
	st.Storage[name] = m
	return nil
}

func (st *InMemoryStorage) FindMetric(name string) (*metric.AbstractMetric, bool, error) {
	m, exists := st.Storage[name]
	return &m, exists, nil
}
