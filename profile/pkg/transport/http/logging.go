package http

import (
	"net/http"
	"time"

	"github.com/go-kit/log"
)

// LoggingMiddleware returns a handler that logs HTTP requests and responses.
func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			begin := time.Now()

			// Create a wrapper for the response writer to capture status code
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			defer func() {
				logger.Log(
					"method", r.Method,
					"path", r.URL.Path,
					"status", wrapper.statusCode,
					"took", time.Since(begin),
					"remote_addr", r.RemoteAddr,
					"user_agent", r.UserAgent(),
				)
			}()

			next.ServeHTTP(wrapper, r)
		})
	}
}

// responseWriterWrapper captures the status code of a response
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (rww *responseWriterWrapper) WriteHeader(code int) {
	rww.statusCode = code
	rww.ResponseWriter.WriteHeader(code)
}
