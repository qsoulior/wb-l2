package app

import (
	"context"
	"dev11/app/repo"
	"dev11/app/service"
	"dev11/app/transport/http"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run настраивает и запускает приложение.
func Run(host string, port string) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := repo.NewEventMemory()
	service := service.NewEventV1(repo)

	server := http.NewServer(host, port, service, logger)
	server.Start(ctx)
	logger.Info("http server started", "host", host, "port", port)

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := server.Stop(ctx)
		if err != nil {
			logger.Error("failed to stop http server", "err", err)
		} else {
			logger.Info("http server has been stopped")
		}
		cancel()
	case err := <-server.Err():
		logger.Error("http server returned error", "err", err)
	}

	stop()
}
