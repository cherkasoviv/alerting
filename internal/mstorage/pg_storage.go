package mstorage

import (
	"alerting/internal/config"
	metric "alerting/internal/metrics"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"strconv"
	"sync"
	"time"
)

type PgStorage struct {
	connString string
	mx         sync.Mutex
}

func InitializePgStorage(cfg *config.ServerConfig) (*PgStorage, error) {
	storage := PgStorage{connString: cfg.DatabaseDSN}
	db, err := sql.Open("pgx", storage.connString)
	db.ExecContext(context.Background(), ""+
		"create schema if not exists public;")
	db.ExecContext(context.Background(),
		"create table if not exists public.metrics"+
			"    (name  varchar(50)  not null constraint metrics_pk primary key,"+
			"    type  varchar(50)  not null,"+
			"    value varchar(50) not null);")
	if err != nil {
		return nil, err
	}
	return &storage, err
}

func (pgStorage *PgStorage) HealthCheck() error {
	db, err := sql.Open("pgx", pgStorage.connString)

	if err != nil {
		return err
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}

	return nil

}

func (pgStorage *PgStorage) CreateOrUpdateMetric(m metric.AbstractMetric) error {
	pgStorage.mx.Lock()
	db, err := sql.Open("pgx", pgStorage.connString)
	if err != nil {
		return err
	}
	defer db.Close()
	defer pgStorage.mx.Unlock()
	_, err = db.ExecContext(context.Background(), ""+
		"INSERT INTO public.metrics VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE  SET value = $3", m.GetName(), m.GetType(), m.GetValue())
	if err != nil {
		return err
	}

	return err
}
func (pgStorage *PgStorage) FindMetric(name string) (metric.AbstractMetric, bool, error) {
	db, err := sql.Open("pgx", pgStorage.connString)
	if err != nil {
		return nil, false, err
	}
	defer db.Close()

	metricRow := db.QueryRowContext(context.Background(), ""+
		"SELECT name, type, value\nfrom public.metrics where name = $1",
		name)
	var metricName, metricType, metricValue string
	err = metricRow.Scan(&metricName, &metricType, &metricValue)
	if err != nil {
		return nil, false, err
	}
	metricData := metric.Metric{
		Name:  metricName,
		Mtype: metric.MetricType(metricType),
	}

	switch metricType {
	case "gauge":
		{
			val, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				return nil, false, err
			}
			metricToReturn := metric.GaugeMetric{
				Metric: metricData,
				Value:  val}
			return &metricToReturn, true, nil

		}
	case "counter":
		{
			val, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				return nil, false, err
			}

			metricToReturn := metric.CounterMetric{
				Metric: metricData,
				Value:  uint64(val)}
			return &metricToReturn, true, nil

		}

	}
	return nil, false, err
}

func (pgStorage *PgStorage) FindAllMetrics() (map[string]metric.AbstractMetric, error) {
	db, err := sql.Open("pgx", pgStorage.connString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	metricRows, err := db.QueryContext(context.Background(), ""+
		"SELECT name, type, value\nfrom public.metrics")

	if err != nil {
		return nil, err
	}
	err = metricRows.Err()

	if err != nil {
		return nil, err
	}

	metricsToReturn := map[string]metric.AbstractMetric{}
	for metricRows.Next() {
		var metricName, metricType, metricValue string
		err = metricRows.Scan(&metricName, &metricType, &metricValue)
		if err != nil {
			return nil, err
		}
		metricData := metric.Metric{
			Name:  metricName,
			Mtype: metric.MetricType(metricType),
		}

		switch metricType {
		case "gauge":
			{
				val, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					return nil, err
				}
				metricToReturn := metric.GaugeMetric{
					Metric: metricData,
					Value:  val}
				metricsToReturn[metricName] = &metricToReturn

			}
		case "counter":
			{
				val, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					return nil, err
				}

				metricToReturn := metric.CounterMetric{
					Metric: metricData,
					Value:  uint64(val)}
				metricsToReturn[metricName] = &metricToReturn

			}

		}
	}

	return metricsToReturn, nil
}

func (pgStorage *PgStorage) CreateOrUpdateSeveralMetrics(metrics map[string]metric.AbstractMetric) error {

	db, err := sql.Open("pgx", pgStorage.connString)
	if err != nil {
		return err
	}
	defer db.Close()

	transaction, err := db.Begin()

	if err != nil {
		return err
	}

	for _, m := range metrics {
		_, err = transaction.ExecContext(context.Background(), ""+
			"INSERT INTO public.metrics VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE  SET value = $3", m.GetName(), m.GetType(), m.GetValue())
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	err = transaction.Commit()

	return err
}
