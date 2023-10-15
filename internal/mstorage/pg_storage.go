package mstorage

import (
	"alerting/internal/config"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type PgStorage struct {
	connString string
}

func InitializePgStorage(cfg *config.ServerConfig) *PgStorage {
	storage := PgStorage{connString: cfg.DatabaseDSN}
	return &storage
}

func (storage *PgStorage) HealthCheck() error {
	db, err := sql.Open("pgx", storage.connString)

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
