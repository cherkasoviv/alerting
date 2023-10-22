package handlers

import (
	"alerting/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type UpdateHandler struct {
	storage metricSaver
}

type metricSaver interface {
	CreateOrUpdateMetric(m metrics.AbstractMetric) error
	FindMetric(name string) (metrics.AbstractMetric, bool, error)
	CreateOrUpdateSeveralMetrics(metrics []metrics.AbstractMetric) error
}

type responseForJSONUpdateHandler struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type requestForJSONUpdateHandler struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}

func NewUpdateHandler(str metricSaver) *UpdateHandler {
	return &UpdateHandler{storage: str}
}

func (uhandler *UpdateHandler) CreateOrUpdateFromURLPath() http.HandlerFunc {
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
						Metric: cMetric,
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
					Metric: gMetric,
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

func (uhandler *UpdateHandler) CreateOrUpdateFromJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			return
		}
		var req requestForJSONUpdateHandler
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(req.ID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metricRequestName := req.ID

		var newMetricValue metrics.AbstractMetric
		var newRequestValueForMetric string
		switch req.MType {
		case "counter":
			{
				var exists bool
				newMetricValue, exists, _ = uhandler.storage.FindMetric(metricRequestName)
				if !exists {
					cMetric := metrics.Metric{
						Name:  metricRequestName,
						Mtype: metrics.Counter,
					}
					newMetricValue = &metrics.CounterMetric{
						Metric: cMetric,
					}

				}
				newRequestValueForMetric = strconv.FormatInt(req.Delta, 10)
			}
		case "gauge":
			{

				gMetric := metrics.Metric{
					Name:  metricRequestName,
					Mtype: metrics.Gauge,
				}
				newMetricValue = &metrics.GaugeMetric{
					Metric: gMetric,
				}

				newRequestValueForMetric = strconv.FormatFloat(req.Value, 'f', 20, 64)
			}
		default:
			{
				http.Error(w, "Wrong metric type", http.StatusBadRequest)
				return
			}
		}

		err = newMetricValue.UpdateValue(newRequestValueForMetric)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = uhandler.storage.CreateOrUpdateMetric(newMetricValue)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		metric, _, err := uhandler.storage.FindMetric(req.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := responseForJSONUpdateHandler{
			ID:    metric.GetName(),
			MType: metric.GetType(),
		}

		switch metric.GetType() {
		case "gauge":
			{
				val, _ := strconv.ParseFloat(metric.GetValue(), 64)
				resp.Value = &val
			}
		case "counter":
			{
				delta, _ := strconv.ParseInt(metric.GetValue(), 10, 64)
				resp.Delta = &delta
			}

		}
		render.JSON(w, r, resp)

	}
}

func (uhandler *UpdateHandler) CreateOrUpdateFromJSONArray() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reqMetrics := []requestForJSONUpdateHandler{}
		err := render.DecodeJSON(r.Body, &reqMetrics)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(reqMetrics) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var metricsToSave []metrics.AbstractMetric

		for _, req := range reqMetrics {

			if len(req.ID) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metricRequestName := req.ID

			var newMetricValue metrics.AbstractMetric
			var newRequestValueForMetric string
			switch req.MType {
			case "counter":
				{
					var exists bool
					newMetricValue, exists, _ = uhandler.storage.FindMetric(metricRequestName)
					if !exists {
						cMetric := metrics.Metric{
							Name:  metricRequestName,
							Mtype: metrics.Counter,
						}
						newMetricValue = &metrics.CounterMetric{
							Metric: cMetric,
						}

					}
					newRequestValueForMetric = strconv.FormatInt(req.Delta, 10)
				}
			case "gauge":
				{

					gMetric := metrics.Metric{
						Name:  metricRequestName,
						Mtype: metrics.Gauge,
					}
					newMetricValue = &metrics.GaugeMetric{
						Metric: gMetric,
					}

					newRequestValueForMetric = strconv.FormatFloat(req.Value, 'f', 20, 64)
				}
			default:
				{
					http.Error(w, "Wrong metric type", http.StatusBadRequest)
					return
				}
			}
			err = newMetricValue.UpdateValue(newRequestValueForMetric)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			metricsToSave = append(metricsToSave, newMetricValue)

		}
		err = uhandler.storage.CreateOrUpdateSeveralMetrics(metricsToSave)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
}
