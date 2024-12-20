package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)
		fullURL := r.Host + r.URL.Path
		logrus.WithFields(logrus.Fields{
			"Method":     r.Method,
			"Path":       fullURL,
			"StatusCode": wrappedWriter.statusCode,
			"Duration":   time.Since(start),
		}).Info()
	})
}
