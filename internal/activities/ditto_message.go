package activities

import (
	"context"
	"dm-backend/internal/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SendDittoProtocolMessageParams struct {
	ThingId string
	Message models.DittoProtocolMessage // Use the new model
}

func validateDittoTopicRegex(topic string) error {
	// Ditto protocol topics can have variants for things and policies.
	// Example valid topics:
	// <namespace>/<name>/things/twin/commands/create
	// <namespace>/<name>/things/live/commands/modify
	// <namespace>/<name>/things/twin/events/update
	// <namespace>/<name>/things/live/messages/modify
	// <namespace>/<name>/policies/commands/modify
	var dittoTopicRegex = regexp.MustCompile(`^[^/]+/[^/]+/(things/(twin|live)/(commands|events|messages)/(create|modify|delete|update)|policies/commands/(create|modify|delete|update))$`)

	if !dittoTopicRegex.MatchString(topic) {
		return fmt.Errorf("topic '%s' does not match Ditto protocol specification", topic)
	}
	return nil
}

// SendDittoProtocolMessageActivity connects to Ditto WS, sends the message, and closes the connection
func (c *DittoClient) SendDittoProtocolMessage(ctx context.Context, params SendDittoProtocolMessageParams) error {
	wsURL := fmt.Sprintf("ws://%s/ws/2", c.Host)

	// Split thingId into namespace and name
	parts := strings.SplitN(params.ThingId, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid thingId format, expected namespace:name")
	}
	namespace, name := parts[0], parts[1]

	params.Message.Topic = strings.ReplaceAll(params.Message.Topic, "<namespace>", namespace)
	params.Message.Topic = strings.ReplaceAll(params.Message.Topic, "<name>", name)

	// Validate topic
	if err := validateDittoTopicRegex(params.Message.Topic); err != nil {
		return err
	}

	// Add correlation-id header with a UUID
	if params.Message.Headers == nil {
		params.Message.Headers = make(map[string]interface{})
	}
	params.Message.Headers["correlation-id"] = uuid.NewString()

	if err := c.ConnectWebSocket(wsURL); err != nil {
		return fmt.Errorf("failed to connect to Ditto WebSocket: %w", err)
	}
	defer c.CloseWebSocket()

	c.wsMutex.Lock()
	defer c.wsMutex.Unlock()
	if c.wsConn == nil {
		return fmt.Errorf("websocket not connected")
	}

	jsonBytes, err := json.Marshal(params.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal Ditto protocol message: %w", err)
	}
	fmt.Printf("Sending to websocket: %s\n", string(jsonBytes))

	if err := c.wsConn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
		return fmt.Errorf("failed to send Ditto protocol message: %w", err)
	}
	return nil
}

func (a *Activities) SendDittoProtocolMessage(ctx context.Context, params SendDittoProtocolMessageParams) error {
	return a.DittoClient.SendDittoProtocolMessage(ctx, params)
}
