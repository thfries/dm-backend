package main

import (
    "log"
    "mass-device-config/internal/api"
    "mass-device-config/internal/workflow"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
)

func main() {
    // Connect to Temporal server
    c, err := client.NewClient(client.Options{})
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