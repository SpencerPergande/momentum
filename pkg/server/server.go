// Package server runs a Momentum service's HTTP server: health checks,
// HTTP/1.1 plus unencrypted HTTP/2 (so Connect and gRPC both work behind the
// cluster's ingress), JSON logs and graceful shutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Options configures Run.
type Options struct {
	// Name is the service name, used in logs.
	Name string
	// Addr is the address to listen on, for example ":8080". Ignored when
	// Listener is set.
	Addr string
	// Listener, when set, is used instead of listening on Addr.
	Listener net.Listener
	// Handler serves every path except the health checks.
	Handler http.Handler
	// Ready reports whether the service can take traffic. When nil, the
	// service is always ready.
	Ready func(context.Context) error
	// ShutdownTimeout bounds graceful shutdown. Defaults to 15 seconds.
	ShutdownTimeout time.Duration
	// Logger defaults to slog.Default().
	Logger *slog.Logger
}

// NewLogger returns a JSON logger that tags every record with the service
// name. The level comes from MOMENTUM_LOG_LEVEL (debug, info, warn, error)
// and defaults to info.
func NewLogger(service string) *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("MOMENTUM_LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(h).With("service", service)
}

// Handler returns the full handler Run serves: /healthz, /readyz, and
// everything else routed to o.Handler.
func Handler(o Options) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if o.Ready != nil {
			if err := o.Ready(r.Context()); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = fmt.Fprintf(w, "not ready: %v\n", err)
				return
			}
		}
		_, _ = w.Write([]byte("ready\n"))
	})
	if o.Handler != nil {
		mux.Handle("/", o.Handler)
	}
	return mux
}

// Run serves until ctx is cancelled, then shuts down gracefully. It returns
// nil after a clean shutdown.
func Run(ctx context.Context, o Options) error {
	log := o.Logger
	if log == nil {
		log = slog.Default()
	}
	timeout := o.ShutdownTimeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}

	ln := o.Listener
	if ln == nil {
		var err error
		ln, err = net.Listen("tcp", o.Addr)
		if err != nil {
			return fmt.Errorf("server: listen on %s: %w", o.Addr, err)
		}
	}

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	srv := &http.Server{
		Handler:           Handler(o),
		Protocols:         protocols,
		ReadHeaderTimeout: 10 * time.Second,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}

	errc := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", ln.Addr().String())
		errc <- srv.Serve(ln)
	}()

	select {
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
	}

	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server: shutdown: %w", err)
	}
	if err := <-errc; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}
