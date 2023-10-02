package mstorage

import (
	"alerting/internal/config"
	metric "alerting/internal/metrics"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type InMemoryStorage struct {
	storage          map[string]metric.AbstractMetric `json:"storage"`
	mx               sync.Mutex
	storagePath      string
	storeInterval    time.Duration
	syncSave         bool
	needPeriodicSave bool
}

type jsonMetric struct {
	Metric metric.Metric
	Value  json.RawMessage
}

func Initialize(cfg *config.ServerConfig) *InMemoryStorage {

	str := map[string]jsonMetric{}
	stor := map[string]metric.AbstractMetric{}
	if cfg.NeedToRestore && cfg.FileStoragePath != "" {
		storageFromFile, _ := os.ReadFile(cfg.FileStoragePath)
		err := json.Unmarshal(storageFromFile, &str)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(str)
	for n, m := range str {
		switch m.Metric.Mtype {
		case metric.Gauge:
			{
				v, _ := strconv.ParseFloat(string(m.Value), 64)
				stor[n] = &metric.GaugeMetric{
					Metric: m.Metric,
					Value:  v,
				}
			}
		case metric.Counter:
			{
				v, _ := strconv.ParseUint(string(m.Value), 10, 64)

				stor[n] = &metric.CounterMetric{
					Metric: m.Metric,
					Value:  v,
				}
			}

		}

	}

	var memStorage = InMemoryStorage{
		storage: stor}
	if cfg.StoreInterval > 0 && cfg.FileStoragePath != "" {

		memStorage.needPeriodicSave = true
		memStorage.syncSave = false
		memStorage.storeInterval = time.Duration(cfg.StoreInterval) * time.Second
		memStorage.storagePath = cfg.FileStoragePath

	} else if cfg.StoreInterval == 0 && cfg.FileStoragePath != "" {

		memStorage.needPeriodicSave = false
		memStorage.syncSave = true
		memStorage.storagePath = cfg.FileStoragePath

	} else {
		memStorage.needPeriodicSave = false
		memStorage.syncSave = false
	}

	go func(st *InMemoryStorage) {
		for st.needPeriodicSave {
			st.mx.Lock()
			err := st.saveToFile()
			if err != nil {
				fmt.Println(err)
			}
			st.mx.Unlock()
			time.Sleep(st.storeInterval)
		}
	}(&memStorage)

	return &memStorage
}

func (st *InMemoryStorage) FindAllMetrics() (map[string]metric.AbstractMetric, error) {
	return st.storage, nil
}

func (st *InMemoryStorage) CreateOrUpdateMetric(m metric.AbstractMetric) error {
	st.mx.Lock()
	defer st.mx.Unlock()
	name := m.GetName()
	str := st.storage
	str[name] = m
	if st.syncSave {
		st.saveToFile()
	}
	return nil
}

func (st *InMemoryStorage) FindMetric(name string) (metric.AbstractMetric, bool, error) {
	str := st.storage
	m, exists := str[name]
	return m, exists, nil
}

func (st *InMemoryStorage) saveToFile() error {

	strJSON, _ := json.MarshalIndent(st.storage, "", "")
	os.MkdirAll(filepath.Dir(st.storagePath), 0777)
	db, err := os.Create(st.storagePath)
	db.Write(strJSON)
	db.Close()

	return err
}
