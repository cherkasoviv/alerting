package handlers

import (
	metric "alerting/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type valueHandler struct {
	storage metricGetter
}

type metricGetter interface {
	FindMetric(name string) (metric.AbstractMetric, bool, error)
	FindAllMetrics() (map[string]metric.AbstractMetric, error)
}

func NewValueHandler(str metricGetter) *valueHandler {
	return &valueHandler{storage: str}
}

func (vhandler *valueHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		allMetricsInStorage, _ := vhandler.storage.FindAllMetrics()
		for _, metric := range allMetricsInStorage {
			metricAsString := metric.String()
			w.Write([]byte(metricAsString))
			w.Write([]byte("\n"))
		}
	}
}

func (vhandler *valueHandler) GetByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricName := chi.URLParam(r, "metricName")
		metric, exist, err := vhandler.storage.FindMetric(metricName)
		if exist && err == nil {

			w.Write([]byte(metric.GetValue()))

		} else {
			http.Error(w, "No such metric", http.StatusNotFound)
		}
	}
}
