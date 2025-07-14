package main

import (
	"dm-backend/internal/activities"
	"dm-backend/internal/api"
	"dm-backend/internal/config"
	"dm-backend/internal/workflow"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("Using TEMPORAL_HOSTPORT from config: %s", cfg.TemporalHost)

	c, err := client.NewClient(client.Options{
		HostPort: cfg.TemporalHost,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	dittoClient := &activities.DittoClient{Host: cfg.DittoHost}
	activitiesImpl := &activities.Activities{Ditto: dittoClient}

	w := worker.New(c, "MASS_DEVICE_CONFIG_TASK_QUEUE", worker.Options{})
	w.RegisterWorkflow(workflow.MassDeviceConfigWorkflow)
	w.RegisterActivity(activitiesImpl.FetchDevicesFromDitto)

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("unable to start worker", err)
		}
	}()

	// Start API server
	api.RunServer(c)
}
