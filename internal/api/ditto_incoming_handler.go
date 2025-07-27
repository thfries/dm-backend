package api

import (
	"context"
	"dm-backend/internal/config"
	"dm-backend/internal/workflow"
	"encoding/json"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/client"
)

func DittoIncomingHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg DittoProtocolMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid Ditto protocol message", http.StatusBadRequest)
			return
		}

		// Dispatch the message for further processing
		if err := DispatchDittoMessage(r.Context(), temporalClient, msg); err != nil {
			http.Error(w, "Failed to process Ditto message: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Ditto message processed"))
	}
}

// Example Ditto protocol message struct (adapt as needed)
type DittoProtocolMessage struct {
	Topic   string                 `json:"topic"`
	Payload map[string]interface{} `json:"payload"`
	Headers map[string]interface{} `json:"headers"`
}

// DispatchDittoMessage can be reused for both HTTP and broker-based incoming messages
func DispatchDittoMessage(ctx context.Context, temporalClient client.Client, msg DittoProtocolMessage) error {
	// Example dispatch logic:
	// - Start a workflow for certain topics
	// - Signal a running workflow/activity for others

	switch msg.Topic {
	case "things/create":
		// Start a workflow
		_, err := temporalClient.ExecuteWorkflow(
			ctx,
			client.StartWorkflowOptions{
				TaskQueue: config.TaskQueue,
			},
			workflow.CreateSiteWorkflow, // or another workflow
			msg.Payload,
		)
		return err
	case "things/update":
		// Signal a running workflow (example)
		workflowID := msg.Headers["workflowId"].(string)
		return temporalClient.SignalWorkflow(ctx, workflowID, "", "ThingUpdated", msg.Payload)
	default:
		// Unhandled topic
		return fmt.Errorf("unhandled Ditto topic: %s", msg.Topic)
	}
}
