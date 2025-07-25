package workflow

import (
	"dm-backend/internal/activities"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CreateSiteParams defines the input for the createSite workflow
type CreateSiteParams struct {
	SiteName    string
	Host        string
	Port        string
	Username    string
	Password    string
	Description string
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
			"siteName":        params.SiteName, // unique attribute
			"siteDescription": params.Description,
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

	// 2. Create Gateway Connection
	createConnParams := activities.CreateConnectionParams{
		ConnectionName: params.SiteName + "-conn",
		TemplateName:   "mqtt5",
		Placeholders: map[string]string{
			"ThingID":        thingID,
			"ConnectionName": params.SiteName + "-conn",
			"MQTTHost":       params.Host,
			"MQTTPort":       params.Port,
			"Username":       params.Username,
			"Password":       params.Password,
		},
	}
	err = workflow.ExecuteActivity(ctx, "CreateConnection", createConnParams).Get(ctx, nil)
	if err != nil {
		// Compensation: delete the thing if connection creation fails after retries
		// Placeholder for DeleteThing activity
		_ = workflow.ExecuteActivity(ctx, "DeleteThing", thingID)
		return fmt.Errorf("createGatewayConnection failed after retries: %w", err)
	}
	return nil
}
