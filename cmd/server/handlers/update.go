package handlers

import (
	"alerting/cmd/server/mstorage"
	"alerting/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type MetricHandler struct {
	Storage *mstorage.MetricStorage
}

func (m MetricHandler) UpdateRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		return
	}
	var metricRequestType metrics.MetricType
	var newMetricValue metrics.AbstractMetric
	urlParams := strings.Split(req.URL.String(), "/")
	if len(urlParams) < 5 {
		http.Error(res, "Not enough data", http.StatusNotFound)
		return
	}
	urlMetricType := chi.URLParam(req, "metricType")
	metricRequestName := chi.URLParam(req, "metricName")
	metricRequestValue := chi.URLParam(req, "metricValue")
	if len(metricRequestName) == 0 {
		http.Error(res, "Wrong metric Name", http.StatusNotFound)
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
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
			return
		}
	}
	str := *m.Storage
	switch metricRequestType {
	case metrics.Counter:
		{
			var exists bool
			newMetricValue, exists, _ = str.FindMetric(metricRequestName)
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
	str.CreateOrUpdateMetric(newMetricValue)

	if err != nil {
		http.Error(res, "Wrong value", http.StatusBadRequest)
		return
	}

}
