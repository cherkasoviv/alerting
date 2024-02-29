package main

import (
	"alerting/internal/config"
	"alerting/internal/mstorage"
	"alerting/internal/server/handlers"
	"alerting/internal/server/middleware/compress"
	"alerting/internal/server/middleware/hash"
	mwLogger "alerting/internal/server/middleware/logger"
	_ "encoding/json"
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadServerConfig()

	storage, _ := mstorage.InitializePgStorage(cfg)

	updateHandler := handlers.NewUpdateHandler(storage)
	valueHandler := handlers.NewValueHandler(storage)
	pingHandler := handlers.NewPingHandler(storage)

	if cfg.DatabaseDSN == "" {
		storageInMemory := mstorage.Initialize(cfg)
		updateHandler = handlers.NewUpdateHandler(storageInMemory)
		valueHandler = handlers.NewValueHandler(storageInMemory)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()

	r := chi.NewRouter()
	r.Use(mwLogger.New(&sugar))
	r.Use(compress.GzipMiddleware())
	r.Use(hash.Hash256Middleware(cfg.HashSHA256Key))
	r.Mount("/debug", middleware.Profiler())
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
	r.Route("/ping", func(r chi.Router) {
		r.Get("/", pingHandler.Ping())
	})

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", updateHandler.CreateOrUpdateFromJSONArray())
	})
	err = http.ListenAndServe(cfg.Host, r)

	if err != nil {
		panic(err)
	}

}

func Profiler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.Handle("/vars", expvar.Handler())

	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	return r
}
