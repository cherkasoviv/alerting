package main

import (
	"alerting/internal/config"
	"alerting/internal/mstorage"
	"alerting/internal/server/handlers"
	mwLogger "alerting/internal/server/middleware/logger"
	_ "encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	cfg := config.LoadServerConfig()
	storage := mstorage.New()
	updateHandler := handlers.NewUpdateHandler(storage)
	valueHandler := handlers.NewValueHandler(storage)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()

	r := chi.NewRouter()
	r.Use(mwLogger.New(&sugar))
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Route("/update", func(r chi.Router) {
		r.Post("/", updateHandler.CreateOrUpdateFromJSON())
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Post("/", updateHandler.CreateOrUpdateFromURLPath())
		})
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/", valueHandler.GetAll())
		r.Route("/value", func(r chi.Router) {
			r.Post("/", valueHandler.GetJSON())

			r.Route("/{metricType}/{metricName}", func(r chi.Router) {
				r.Get("/", valueHandler.GetByName())
			})
		})

	})
	err = http.ListenAndServe(cfg.Host, r)
	if err != nil {
		panic(err)
	}

}
