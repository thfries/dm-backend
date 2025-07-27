package api

import (
	"dm-backend/internal/config"
	"dm-backend/internal/models"
	"dm-backend/internal/workflow"
	"encoding/json"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

// GetWorkflowStatusHandler handles the retrieval of a workflow's status by its ID and run ID.
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
	RQLQuery             string                      `json:"rql_query"`
	DittoProtocolMessage models.DittoProtocolMessage `json:"ditto_protocol_message"`
}

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
			RQLQuery:             req.RQLQuery,
			DittoProtocolMessage: req.DittoProtocolMessage,
		}
		options := client.StartWorkflowOptions{
			TaskQueue: config.TaskQueue,
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
