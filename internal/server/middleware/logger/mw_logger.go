package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func New(sugarLogger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		sugarLogger.With(zap.String("component", "middleware/logger"))
		sugarLogger.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			startRequestHandlingTime := time.Now()

			entry := sugarLogger.With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("encoding", r.Header.Get("Accept-Encoding")),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				entry = entry.With(
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.String("duration", time.Since(startRequestHandlingTime).String()),
				)
				entry.Info("request complited")

			}()

			next.ServeHTTP(ww, r)

		}
		return http.HandlerFunc(fn)
	}
}
