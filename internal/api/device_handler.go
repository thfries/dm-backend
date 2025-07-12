// Add below StartMassDeviceConfigHandler in the same file

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