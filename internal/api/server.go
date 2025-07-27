package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.temporal.io/sdk/client"
)

// RunServer initializes and starts the HTTP server
func RunServer(temporalClient client.Client) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/config/start", StartMassDeviceConfigHandler(temporalClient))
	mux.HandleFunc("/api/config/status", GetWorkflowStatusHandler(temporalClient))
	mux.HandleFunc("/api/sites/create", StartCreateSitesHandler(temporalClient))
	mux.HandleFunc("/api/ditto/incoming", DittoIncomingHandler(temporalClient))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("docs"))))

	server := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting HTTP server on :18080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop // Wait for signal

	log.Println("Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server stopped gracefully")
}
