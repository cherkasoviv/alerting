package handlers

import (
	"alerting/internal/mstorage"
	"net/http"
)

type PingHandler struct {
	pg *mstorage.PgStorage
}

func NewPingHandler(storage *mstorage.PgStorage) *PingHandler {
	return &PingHandler{
		pg: storage,
	}
}

func (pHandler *PingHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pHandler.pg.HealthCheck(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

}
