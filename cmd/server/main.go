package main

import (
	"alerting/internal/config"
	"alerting/internal/handlers/metrics/update"
	"alerting/internal/handlers/metrics/value"
	mstorage "alerting/internal/mstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {

	cfg := config.LoadServerConfig()
	storage := mstorage.New()

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Post("/", update.CreateOrUpdate(storage))
		})
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/", value.GetAll(storage))
		r.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
			r.Get("/", value.GetByName(storage))
		})
	})

	err := http.ListenAndServe(cfg.Host, r)
	if err != nil {
		panic(err)
	}

}
