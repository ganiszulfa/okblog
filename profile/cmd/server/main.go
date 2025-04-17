package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ganis/okblog/profile/pkg/service"
	httptransport "github.com/ganis/okblog/profile/pkg/transport/http"
	"github.com/go-kit/log"
)

func main() {
	// Create a logger
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// Create service with logging middleware
	var svc service.Service
	svc = service.NewService()
	svc = service.LoggingMiddleware(logger)(svc)

	// Create HTTP server with the service
	server := httptransport.NewServer(svc, logger)

	// Create a channel to listen for errors coming from the listener.
	errs := make(chan error, 2)

	// Start the server
	go func() {
		logger.Log("transport", "HTTP", "addr", ":8080", "msg", "Starting server")
		errs <- http.ListenAndServe(":8080", server)
	}()

	// Listen for an interrupt signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		logger.Log("signal", "interrupt", "msg", "Shutting down")
		errs <- nil
	}()

	logger.Log("exit", <-errs)
}
