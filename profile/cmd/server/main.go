package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ganis/okblog/profile/pkg/service"
	httptransport "github.com/ganis/okblog/profile/pkg/transport/http"
)

func main() {
	svc := service.NewService()
	server := httptransport.NewServer(svc)

	// Create a channel to listen for errors coming from the listener.
	errs := make(chan error, 2)

	// Start the server
	go func() {
		log.Printf("Starting server on :8080")
		errs <- http.ListenAndServe(":8080", server)
	}()

	// Listen for an interrupt signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		errs <- nil
	}()

	log.Printf("exit", <-errs)
}
