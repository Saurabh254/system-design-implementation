package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/saurabh254/system-design-implementation/ratelimiter/internal/utils"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	logger.Log.Info("starting server", "addr", s.httpServer.Addr)

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server crashed", "err", err)
		}
		errCh <- err
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		logger.Log.Info("shutdown triggered by context", "err", ctx.Err())

	case <-stop:
		logger.Log.Info("shutdown triggered by signal")

	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		logger.Log.Info("server stopped")
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

	logger.Log.Info("initiating graceful shutdown", "timeout", "5s")

	err := s.httpServer.Shutdown(shutdownCtx)
	cancel()

	if err != nil {
		logger.Log.Error("shutdown failed", "err", err)
		return err
	}

	logger.Log.Info("shutdown complete")
	return nil
}
