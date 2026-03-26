package server

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// New returns an *http.Server configured with safe timeouts. The caller is
// responsible for setting the Handler before calling ListenAndServe.
func New(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// GracefulShutdown returns a context that is cancelled when SIGINT or SIGTERM
// is received. Pass the returned cancel to defer to clean up resources.
func GracefulShutdown() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}
