package handlers

import (
	metrics "alerting/internal/metrics"
	"net/http"
	"strings"
)

func UpdateRequest(res http.ResponseWriter, req *http.Request) {
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
	urlMetricType := urlParams[2]
	metricRequestName := urlParams[3]
	metricRequestValue := urlParams[4]
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

	switch metricRequestType {
	case metrics.Counter:
		{
			cMetric := metrics.Metric{
				Name:  metricRequestName,
				Mtype: metrics.Counter,
			}
			newMetricValue = &metrics.CounterMetric{
				CMetric: cMetric,
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
	if err != nil {
		http.Error(res, "Wrong value", http.StatusBadRequest)
		return
	}

}
