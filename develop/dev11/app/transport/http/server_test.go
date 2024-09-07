package http

import (
	"context"
	"dev11/app/service"
	"log/slog"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestNewServer(t *testing.T) {
	t.Run("InvalidArguments", func(t *testing.T) {
		var want *Server = nil
		if got := NewServer("", "", nil, nil); got != want {
			t.Errorf("NewServer() = %v, want %v", got, want)
		}
	})
	t.Run("ValidArguments", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		service := service.NewMockEvent(ctrl)
		logger := new(slog.Logger)

		if got := NewServer("", "", service, logger); got == nil {
			t.Errorf("NewServer() = %v, want non-nil", got)
		}
	})
}

func TestServer_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpServer := &http.Server{Handler: http.NotFoundHandler()}

	errCh := make(chan error, 1)
	s := Server{httpServer: httpServer, errCh: errCh}
	s.Start(ctx)
	s.httpServer.Close()

	want := http.ErrServerClosed
	if got := <-errCh; got != want {
		t.Errorf("NewServer() got = %v, want %v", got, want)
	}

	if got1 := httpServer.BaseContext(nil); got1 != ctx {
		t.Errorf("NewServer() got1 = %v, want %v", got1, ctx)
	}
}

func TestServer_Stop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := Server{httpServer: &http.Server{}}

	if err := s.Stop(ctx); err != nil {
		t.Errorf("Server.Stop() error = %v, wantErr %v", err, false)
	}
}

func TestServer_Err(t *testing.T) {
	want := make(chan error, 1)
	close(want)
	s := Server{httpServer: &http.Server{}, errCh: want}

	if got := s.Err(); got != want {
		t.Errorf("Server.Err() = %v, want %v", got, want)
	}
}
