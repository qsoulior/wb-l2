package http

import (
	"log/slog"
	"net/http"
	"time"
)

// Тип промежуточного HTTP-обработчика.
type Middleware func(next http.Handler) http.Handler

// Кастомный http.ResponseWriter для логирования.
type loggerWriter struct {
	http.ResponseWriter
	code int
	size int
}

// Write записывает байты в тело HTTP-ответа.
func (w *loggerWriter) Write(bytes []byte) (int, error) {
	n, err := w.ResponseWriter.Write(bytes)
	w.size += n
	return n, err
}

// Write записывает заголовок HTTP-ответа.
func (w *loggerWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LoggerMiddleware возвращает middleware для логирования запросов.
func LoggerMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggerWriter{ResponseWriter: w}
			start := time.Now()
			next.ServeHTTP(lw, r)
			logger.Info(r.RemoteAddr, "method", r.Method, "url", r.URL.String(), "proto", r.Proto, "code", lw.code, "size", lw.size, "time", time.Since(start))
		})
	}
}

// RecovererMiddleware возвращает middleware для обработки внутренних ошибок.
func RecovererMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rc := recover(); rc != nil {
					w.WriteHeader(http.StatusInternalServerError)
					if err, ok := rc.(error); ok {
						logger.Error(r.RemoteAddr, "err", err)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
