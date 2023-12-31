package handlers

import (
	metric "alerting/internal/metrics"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type valueHandler struct {
	storage metricGetter
	//TODO добавить логгер
}

type responseForJSONValueHandler struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type requestForJSONValueHandler struct {
	ID    string `json:"id"`
	MType string `json:"type"`
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

		w.Header().Add("Content-Type", "text/html")
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
		w.Header().Add("Content-Type", "text/html")
		metricName := chi.URLParam(r, "metricName")
		metric, exist, err := vhandler.storage.FindMetric(metricName)
		if exist && err == nil {

			w.Write([]byte(metric.GetValue()))

		} else {
			http.Error(w, "No such metric", http.StatusNotFound)
		}
	}
}

func (vhandler *valueHandler) GetJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requestForJSONValueHandler
		err := render.DecodeJSON(r.Body, &req)
		fmt.Println(req)
		if err != nil || len(req.ID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		findMetric, exist, err := vhandler.storage.FindMetric(req.ID)

		if !exist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if findMetric.GetType() != req.MType {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp := responseForJSONValueHandler{
			ID:    findMetric.GetName(),
			MType: findMetric.GetType(),
		}

		switch findMetric.GetType() {
		case "gauge":
			{

				val, err := strconv.ParseFloat(findMetric.GetValue(), 64)
				resp.Value = &val
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		case "counter":
			{
				val, err := strconv.ParseInt(findMetric.GetValue(), 10, 64)
				resp.Delta = &val
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

		}

		render.JSON(w, r, resp)

	}
}
