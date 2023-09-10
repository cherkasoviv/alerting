package main

import (
	"alerting/cmd/server/handlers"
	metric "alerting/internal/metrics"
	mstorage2 "alerting/internal/mstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	var MemStorage mstorage2.MetricStorage = mstorage2.InMemoryStorage{
		Storage: map[string]metric.AbstractMetric{}}

	mHandler := handlers.MetricHandler{
		Storage: &MemStorage,
	}
	parseFlags()
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

	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err)
	}

}
