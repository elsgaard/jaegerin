// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
// version, commit, date is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var (
	release = "unknown"
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

type Server struct {
	address string
	logger  *slog.Logger
	mux     *http.ServeMux
	server  *http.Server
}

type Options struct {
	Host string
	Log  *slog.Logger
	Port int
}

func New(opts Options) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := http.NewServeMux()
	return &Server{
		address: address,
		logger:  opts.Log,
		mux:     mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start the Server by setting up routes and listening for HTTP requests on the given address.
func (s *Server) Start() error {

	s.setupRoutes()

	s.logger.Info("Starting HTTP/OTLP server", slog.String("version", commit), slog.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Stop the Server gracefully within the timeout.
func (s *Server) Stop() error {
	s.logger.Info("Stopping HTTP/OLTP server", slog.String("version", commit), slog.String("address", s.address))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
