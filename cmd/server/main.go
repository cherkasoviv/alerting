package main

import (
	"alerting/cmd/server/handlers"
	"alerting/cmd/server/mstorage"
	metric "alerting/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	var MemStorage mstorage.MetricStorage

	MemStorage = mstorage.InMemoryStorage{
		Storage: map[string]metric.AbstractMetric{},
	}

	var mHandler handlers.MetricHandler

	mHandler = handlers.MetricHandler{
		Storage: &MemStorage,
	}
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Post("/", mHandler.UpdateRequest)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", mHandler.ValueAllMetrics)
		r.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
			r.Get("/", mHandler.ValueMetricByName)
		})
	})

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}

}
