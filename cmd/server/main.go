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

	dittoClient := &activities.DittoClient{Host: cfg.DittoHost, Username: cfg.DittoUsername, Password: cfg.DittoPassword}
	activitiesImpl := &activities.Activities{DittoClient: dittoClient}

	w := worker.New(c, "MASS_DEVICE_CONFIG_TASK_QUEUE", worker.Options{})
	w.RegisterWorkflow(workflow.MassDeviceConfigWorkflow)
	w.RegisterActivity(activitiesImpl.FetchDevicesFromDitto)
	w.RegisterActivity(activitiesImpl.ConfigureDevice)

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("unable to start worker", err)
		}
	}()

	// Start API server
	api.RunServer(c)
}
