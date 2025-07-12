package api

import (
    "log"
    "net/http"
    "go.temporal.io/sdk/client"
)

func RunServer(temporalClient client.Client) {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/config/start", StartMassDeviceConfigHandler(temporalClient))
    mux.HandleFunc("/api/config/status", GetWorkflowStatusHandler(temporalClient))

    log.Println("Starting HTTP server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}