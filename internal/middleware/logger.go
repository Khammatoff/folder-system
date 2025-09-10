package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				latency := time.Since(start)
				logger.WithFields(logrus.Fields{
					"status":     ww.Status(),
					"method":     r.Method,
					"path":       r.URL.Path,
					"ip":         r.RemoteAddr,
					"latency":    latency,
					"user_agent": r.UserAgent(),
				}).Info("HTTP request")
			}()
			next.ServeHTTP(ww, r)
		})
	}
}
