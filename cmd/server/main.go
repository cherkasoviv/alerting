package main

import (
	"alerting/internal/config"
	"alerting/internal/handlers"
	mstorage "alerting/internal/mstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {

	cfg := config.LoadServerConfig()
	storage := mstorage.New()
	updateHandler := handlers.NewUpdateHandler(storage)
	valueHandler := handlers.NewValueHandler(storage)

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Post("/", updateHandler.CreateOrUpdate())
		})
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/", valueHandler.GetAll())
		r.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
			r.Get("/", valueHandler.GetByName())
		})
	})

	err := http.ListenAndServe(cfg.Host, r)
	if err != nil {
		panic(err)
	}

}
