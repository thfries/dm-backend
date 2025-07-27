package api

import (
	"encoding/json"
	"net/http"

	"dm-backend/internal/config"
	"dm-backend/internal/workflow"

	"go.temporal.io/sdk/client"
)

func StartCreateSitesHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params workflow.CreateSiteBatchWorkflowParams
		if err := json.NewDecoder(r.Body).Decode(&params.Sites); err != nil {
			http.Error(w, "Invalid JSON input", http.StatusBadRequest)
			return
		}

		options := client.StartWorkflowOptions{
			TaskQueue: config.TaskQueue,
		}

		workflowRun, err := temporalClient.ExecuteWorkflow(
			r.Context(),
			options,
			workflow.CreateSiteBatchWorkflow,
			params,
		)
		if err != nil {
			http.Error(w, "Failed to start batch workflow: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"workflowID": workflowRun.GetID(),
			"runID":      workflowRun.GetRunID(),
		})
	}
}
