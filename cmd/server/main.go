package main

import (
    "log"
    "os"
    "dm-backend/internal/api"
    "dm-backend/internal/workflow"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
)

func main() {
    // Read Temporal endpoint from environment variable, default to localhost:7233
    temporalHostPort := os.Getenv("TEMPORAL_HOSTPORT")
    if temporalHostPort == "" {
        temporalHostPort = "localhost:7233"
        log.Printf("TEMPORAL_HOSTPORT not set, using default: %s", temporalHostPort)
    } else {
        log.Printf("Using TEMPORAL_HOSTPORT from environment: %s", temporalHostPort)
    }

    c, err := client.NewClient(client.Options{
        HostPort: temporalHostPort,
    })
    if err != nil {
        log.Fatalln("unable to create Temporal client", err)
    }
    defer c.Close()

    w := worker.New(c, "MASS_DEVICE_CONFIG_TASK_QUEUE", worker.Options{})
    w.RegisterWorkflow(workflow.MassDeviceConfigWorkflow)
    // Register activities here...

    go func() {
        if err := w.Run(worker.InterruptCh()); err != nil {
            log.Fatalln("unable to start worker", err)
        }
    }()

    // Start API server
    api.RunServer(c)
}