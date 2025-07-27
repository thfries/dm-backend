package main

import (
	"log"

	_ "dm-backend/docs"
	"dm-backend/internal/activities"
	"dm-backend/internal/api"
	"dm-backend/internal/config"
	"dm-backend/internal/workflow"

	"go.temporal.io/sdk/activity"
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

	dittoClient := &activities.DittoClient{
		Host:           cfg.DittoHost,
		Username:       cfg.DittoUsername,
		Password:       cfg.DittoPassword,
		DevopsUsername: cfg.DittoDevopsUsername,
		DevopsPassword: cfg.DittoDevopsPassword,
	}
	activitiesImpl := &activities.Activities{DittoClient: dittoClient}

	w := worker.New(c, config.TaskQueue, worker.Options{})
	w.RegisterWorkflow(workflow.MassDeviceConfigWorkflow)
	w.RegisterWorkflow(workflow.CreateSiteWorkflow)
	w.RegisterWorkflow(workflow.CreateSiteBatchWorkflow)
	w.RegisterActivity(activitiesImpl.FetchDevicesFromDitto)
	w.RegisterActivity(activitiesImpl.ConfigureDevice)
	w.RegisterActivity(activitiesImpl.CreateConnection)
	w.RegisterActivity(activitiesImpl.GetConnectionStatus)
	w.RegisterActivity(activitiesImpl.CreateThing)
	w.RegisterActivity(activitiesImpl.DeleteThing)
	w.RegisterActivityWithOptions(activitiesImpl.SendDittoProtocolMessage, activity.RegisterOptions{Name: "SendDittoProtocolMessage"})
	w.RegisterActivityWithOptions(activitiesImpl.SendDittoProtocolMessage, activity.RegisterOptions{Name: "UpdateGatewayPolicy"})

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("unable to start worker", err)
		}
	}()

	// Start API server
	api.RunServer(c)
}
