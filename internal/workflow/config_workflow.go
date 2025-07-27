package workflow

import (
	"dm-backend/internal/models"
	"time"

	"go.temporal.io/sdk/workflow"
)

type ConfigWorkflowParams struct {
	RQLQuery             string
	DittoProtocolMessage models.DittoProtocolMessage
}

func MassDeviceConfigWorkflow(ctx workflow.Context, params ConfigWorkflowParams) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Fetch devices from Ditto
	var devices []models.Device
	if err := workflow.ExecuteActivity(ctx, "FetchDevicesFromDitto", params.RQLQuery).Get(ctx, &devices); err != nil {
		return err
	}

	// Send Ditto protocol message for each device
	futures := make([]workflow.Future, len(devices))
	for i, device := range devices {
		// Pass ThingId and the protocol message to the activity
		activityParams := struct {
			ThingId string
			Message models.DittoProtocolMessage
		}{
			ThingId: device.ThingId,
			Message: params.DittoProtocolMessage,
		}
		futures[i] = workflow.ExecuteActivity(ctx, "SendDittoProtocolMessage", activityParams)
	}
	for _, f := range futures {
		if err := f.Get(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}
