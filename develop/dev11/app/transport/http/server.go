package http

import (
	"context"
	"dev11/app/service"
	"dev11/app/transport/http/handler"
	"log/slog"
	"net"
	"net/http"
)

// Обертка над http-сервером с маршрутами, промежуточными слоями и методами Start, Stop, Err.
type Server struct {
	httpServer *http.Server
	errCh      chan error
}

// NewServer возвращает новый http-сервер, если service и logger не равны nil.
func NewServer(host string, port string, service service.Event, logger *slog.Logger) *Server {
	if service == nil || logger == nil {
		return nil
	}

	router := http.NewServeMux()
	router.Handle("POST /create_event", handler.EventCreate{Service: service})
	router.Handle("POST /update_event", handler.EventUpdate{Service: service})
	router.Handle("POST /delete_event", handler.EventDelete{Service: service})
	router.Handle("GET /events_for_day", handler.EventGetForDay{Service: service})
	router.Handle("GET /events_for_week", handler.EventGetForWeek{Service: service})
	router.Handle("GET /events_for_month", handler.EventGetForMonth{Service: service})

	var mux http.Handler = router
	middlewares := []Middleware{RecovererMiddleware(logger), LoggerMiddleware(logger)}
	for _, middleware := range middlewares {
		mux = middleware(mux)
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	return &Server{
		httpServer: httpServer,
		errCh:      make(chan error, 1),
	}
}

// Start запускает http-сервер в отдельной горутине.
func (s *Server) Start(ctx context.Context) {
	s.httpServer.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	go func() {
		s.errCh <- s.httpServer.ListenAndServe()
		close(s.errCh)
	}()
}

// Stop останавливает http-сервер.
func (s *Server) Stop(ctx context.Context) error { return s.httpServer.Shutdown(ctx) }

// Err возвращает канал с ошибками http-сервера.
func (s *Server) Err() <-chan error { return s.errCh }
