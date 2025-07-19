package api

import (
	"dm-backend/internal/workflow"
	"encoding/json"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

func GetWorkflowStatusHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflowID := r.URL.Query().Get("workflowID")
		runID := r.URL.Query().Get("runID")

		if workflowID == "" {
			http.Error(w, "workflowID required", http.StatusBadRequest)
			return
		}

		resp, err := temporalClient.DescribeWorkflowExecution(r.Context(), workflowID, runID)
		if err != nil {
			log.Printf("Failed to get workflow status: %v", err)
			http.Error(w, "failed to get workflow status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

type StartConfigRequest struct {
	RQLQuery     string            `json:"rql_query"`
	ConfigParams map[string]string `json:"config_params"`
}

// Returns a handler function that can be registered with http.HandleFunc
func StartMassDeviceConfigHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req StartConfigRequest
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		params := workflow.ConfigWorkflowParams{
			RQLQuery:     req.RQLQuery,
			ConfigParams: req.ConfigParams,
		}
		options := client.StartWorkflowOptions{
			TaskQueue: "MASS_DEVICE_CONFIG_TASK_QUEUE",
		}

		workflowRun, err := temporalClient.ExecuteWorkflow(
			r.Context(),
			options,
			workflow.MassDeviceConfigWorkflow,
			params,
		)
		if err != nil {
			log.Printf("Failed to start workflow: %v", err)
			http.Error(w, "Failed to start workflow", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"workflowID": workflowRun.GetID(),
			"runID":      workflowRun.GetRunID(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
