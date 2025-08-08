// Package logbookserver is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"context"
	"flag"
	"jaegerin/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
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

func main() {
	os.Exit(start())
}

func start() int {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	hostPtr := flag.String("host", "0.0.0.0", "")
	portPtr := flag.Int("port", 4318, "")

	flag.Parse()

	s := server.New(server.Options{
		Host: *hostPtr,
		Log:  logger,
		Port: *portPtr,
	})

	var eg errgroup.Group
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eg.Go(func() error {
		<-ctx.Done()
		if err := s.Stop(); err != nil {
			logger.Error("Error stopping server", err)
			return err
		}
		return nil
	})

	if err := s.Start(); err != nil {
		logger.Error("Error starting server", err)
		return 1
	}

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}
