package workflow

import (
	"dm-backend/internal/activities"
	"dm-backend/internal/models"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CreateSiteParams defines the input for the createSite workflow
type CreateSiteParams struct {
	Site models.Site
}

// CreateSiteWorkflow creates a gateway thing and a connection, with compensation on failure
func CreateSiteWorkflow(ctx workflow.Context, params CreateSiteParams) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Create Gateway Thing
	thingData := map[string]interface{}{
		"attributes": map[string]interface{}{
			"siteName":        params.Site.SiteName, // unique attribute
			"siteDescription": params.Site.Description,
		},
	}
	createThingParams := activities.CreateThingParams{
		Namespace:          "gateway", // or use from config if needed
		UniqueAttributeKey: "siteName",
		ThingData:          thingData,
	}
	var thingID string
	err := workflow.ExecuteActivity(ctx, "CreateThing", createThingParams).Get(ctx, &thingID)
	if err != nil {
		return err
	}

	var subject = fmt.Sprintf("integration:%s", thingID)

	// 1.5. Update Gateway Policy by sending Ditto Protocol Message
	updatePolicyParams := activities.SendDittoProtocolMessageParams{
		ThingId: thingID,
		Message: models.DittoProtocolMessage{
			Topic:   "<namespace>/<name>/policies/commands/modify",
			Headers: map[string]interface{}{"correlation-id": "update-policy"},
			Path:    "entries/DEVICE",
			Value: map[string]interface{}{
				"subjects": map[string]interface{}{
					subject: map[string]interface{}{
						"type": "connection",
					},
				},
				"resources": map[string]interface{}{
					"thing:/": map[string]interface{}{
						"grant":  []string{"READ", "WRITE"},
						"revoke": []string{},
					},
					"message:/": map[string]interface{}{
						"grant":  []string{"READ", "WRITE"},
						"revoke": []string{},
					},
				},
			},
		},
	}
	err = workflow.ExecuteActivity(ctx, "UpdateGatewayPolicy", updatePolicyParams).Get(ctx, nil)
	if err != nil {
		// Compensation: delete the thing if policy update fails after retries
		_ = workflow.ExecuteActivity(ctx, "DeleteThing", thingID)
		return fmt.Errorf("updateGatewayPolicy failed after retries: %w", err)
	}

	// 2. Create Gateway Connection
	createConnParams := activities.CreateConnectionParams{
		ConnectionName: params.Site.SiteName + "-conn",
		TemplateName:   "mqtt5",
		Placeholders: map[string]string{
			"ThingID":        thingID,
			"ConnectionName": params.Site.SiteName + "-conn",
			"MQTTHost":       params.Site.Host,
			"MQTTPort":       params.Site.Port,
			"Username":       params.Site.Username,
			"Password":       params.Site.Password,
		},
	}
	err = workflow.ExecuteActivity(ctx, "CreateConnection", createConnParams).Get(ctx, nil)
	if err != nil {
		// Compensation: delete the thing if connection creation fails after retries
		_ = workflow.ExecuteActivity(ctx, "DeleteThing", thingID)
		return fmt.Errorf("createGatewayConnection failed after retries: %w", err)
	}
	return nil
}
