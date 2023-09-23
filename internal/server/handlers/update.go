package handlers

import (
	"alerting/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type updateHandler struct {
	storage metricSaver
}

type metricSaver interface {
	CreateOrUpdateMetric(m metrics.AbstractMetric) error
	FindMetric(name string) (metrics.AbstractMetric, bool, error)
}

func NewUpdateHandler(str metricSaver) *updateHandler {
	return &updateHandler{storage: str}
}

func (uhandler *updateHandler) CreateOrUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			return
		}
		var metricRequestType metrics.MetricType
		var newMetricValue metrics.AbstractMetric

		urlMetricType := chi.URLParam(r, "metricType")
		metricRequestName := chi.URLParam(r, "metricName")
		metricRequestValue := chi.URLParam(r, "metricValue")
		if len(metricRequestName) == 0 {
			http.Error(w, "Wrong metric Name", http.StatusNotFound)
			return
		}
		switch urlMetricType {
		case "counter":
			{
				metricRequestType = metrics.Counter
			}
		case "gauge":
			{
				metricRequestType = metrics.Gauge
			}
		default:
			{
				http.Error(w, "Wrong metric type", http.StatusBadRequest)
				return
			}
		}

		switch metricRequestType {
		case metrics.Counter:
			{
				var exists bool
				newMetricValue, exists, _ = uhandler.storage.FindMetric(metricRequestName)
				if !exists {
					cMetric := metrics.Metric{
						Name:  metricRequestName,
						Mtype: metrics.Counter,
					}
					newMetricValue = &metrics.CounterMetric{
						CMetric: cMetric,
					}
				}

			}
		case metrics.Gauge:
			{
				gMetric := metrics.Metric{
					Name:  metricRequestName,
					Mtype: metrics.Gauge,
				}
				newMetricValue = &metrics.GaugeMetric{
					GMetric: gMetric,
				}
			}
		}
		err := newMetricValue.UpdateValue(metricRequestValue)
		uhandler.storage.CreateOrUpdateMetric(newMetricValue)

		if err != nil {
			http.Error(w, "Wrong value", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
