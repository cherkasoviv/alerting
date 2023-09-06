package handlers

import (
	metrics2 "alerting/cmd/server/metrics"
	"net/http"
	"strings"
)

func UpdateRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		return
	}
	var metricRequestType metrics2.MetricType
	var newMetricValue metrics2.AbstractMetric
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
			metricRequestType = metrics2.Counter
		}
	case "gauge":
		{
			metricRequestType = metrics2.Gauge
		}
	default:
		{
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
			return
		}
	}

	switch metricRequestType {
	case metrics2.Counter:
		{
			cMetric := metrics2.Metric{
				Name:  metricRequestName,
				Mtype: metrics2.Counter,
			}
			newMetricValue = &metrics2.CounterMetric{
				CMetric: cMetric,
			}
		}
	case metrics2.Gauge:
		{
			gMetric := metrics2.Metric{
				Name:  metricRequestName,
				Mtype: metrics2.Gauge,
			}
			newMetricValue = &metrics2.GaugeMetric{
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
