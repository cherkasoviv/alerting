package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (m MetricHandler) ValueAllMetrics(res http.ResponseWriter, req *http.Request) {
	str := *m.Storage
	allMetricsInStorage, _ := str.FindAllMetrics()
	for _, metric := range allMetricsInStorage {
		metricAsString, _ := metric.String()
		res.Write([]byte(metricAsString))
		res.Write([]byte("\n"))
	}

}

func (m MetricHandler) ValueMetricByName(res http.ResponseWriter, req *http.Request) {
	str := *m.Storage
	metricName := chi.URLParam(req, "metricName")
	metric, exist, err := str.FindMetric(metricName)
	if exist && err == nil {
		metricAsString, _ := metric.String()
		res.Write([]byte(metricAsString))

	} else {
		http.Error(res, "No such metric", http.StatusNotFound)
	}
}
