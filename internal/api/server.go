package api

import (
	"context"
	"dm-backend/internal/config"
	"dm-backend/internal/workflow"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.temporal.io/sdk/client"
)

func RunServer(temporalClient client.Client) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/config/start", StartMassDeviceConfigHandler(temporalClient))
	mux.HandleFunc("/api/config/status", GetWorkflowStatusHandler(temporalClient))
	mux.HandleFunc("/api/sites/create", StartCreateSitesHandler(temporalClient))

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

func StartCreateSitesHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sites []workflow.CreateSiteParams
		if err := json.NewDecoder(r.Body).Decode(&sites); err != nil {
			http.Error(w, "Invalid JSON input", http.StatusBadRequest)
			return
		}
		for _, site := range sites {
			we, err := temporalClient.ExecuteWorkflow(
				r.Context(),
				client.StartWorkflowOptions{
					TaskQueue: config.TaskQueue,
				},
				workflow.CreateSiteWorkflow,
				site,
			)
			if err != nil {
				http.Error(w, "Failed to start workflow: "+err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Started CreateSiteWorkflow for site %s, WorkflowID: %s", site.SiteName, we.GetID())
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Workflows started"))
	}
}
