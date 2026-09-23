package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

func run(ctx context.Context, listener net.Listener, handler http.Handler) error {
	httpServer := &http.Server{Handler: handler}
	serveErr := make(chan error, 1)
	go func() { serveErr <- httpServer.Serve(listener) }()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			_ = httpServer.Close()
			return fmt.Errorf("shutdown: %w", err)
		}
		if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve: %w", err)
		}
		return nil
	}
}
