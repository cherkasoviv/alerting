package main

import (
	"alerting/internal/config"
	"alerting/internal/mstorage"
	"alerting/internal/server/handlers"
	mwLogger "alerting/internal/server/middleware/logger"
	"github.com/go-chi/chi/v5"
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

	err = http.ListenAndServe(cfg.Host, r)
	if err != nil {
		panic(err)
	}

}
